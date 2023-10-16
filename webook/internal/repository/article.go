package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./article.go -package=repomocks -destination=mocks/article.mock.go ArticleRepository
type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	List(ctx context.Context, author int64,
		offset int, limit int) ([]domain.Article, error)

	// Sync 本身要求先保存到制作库，再同步到线上库
	Sync(ctx context.Context, art domain.Article) (int64, error)
	// SyncStatus 仅仅同步状态
	SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error
	GetById(ctx context.Context, id int64) (domain.Article, error)

	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	// 操作单一的库
	dao article.ArticleDAO
	// 在单体应用下，直接复用 DAO 是可以接受的
	userRepo UserRepository
	cache    cache.ArticleCache

	// SyncV1 用
	authorDAO article.ArticleAuthorDAO
	readerDAO article.ArticleReaderDAO

	// SyncV2 用
	db *gorm.DB
	l  logger.LoggerV1
}

func NewArticleRepository(dao article.ArticleDAO,
	c cache.ArticleCache,
	userRepo UserRepository,
	l logger.LoggerV1) ArticleRepository {
	return &CachedArticleRepository{
		dao:      dao,
		l:        l,
		cache:    c,
		userRepo: userRepo,
	}
}

func NewArticleRepositoryV1(authorDAO article.ArticleAuthorDAO,
	readerDAO article.ArticleReaderDAO) ArticleRepository {
	return &CachedArticleRepository{
		authorDAO: authorDAO,
		readerDAO: readerDAO,
	}
}

func (repo *CachedArticleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := repo.cache.GetPub(ctx, id)
	if err == nil {
		return res, err
	}
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := repo.userRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res = domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		},
	}
	// 也可以同步
	go func() {
		if err = repo.cache.SetPub(ctx, res); err != nil {
			repo.l.Error("缓存已发表文章失败",
				logger.Error(err), logger.Int64("aid", res.Id))
		}
	}()
	return res, nil
}

func (repo *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	cachedArt, err := repo.cache.Get(ctx, id)
	if err == nil {
		return cachedArt, nil
	}
	art, err := repo.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return repo.toDomain(art), nil
}

func (repo *CachedArticleRepository) List(ctx context.Context, author int64,
	offset int, limit int) ([]domain.Article, error) {
	// 只有第一页才走缓存，并且假定一页只有 100 条
	// 也就是说，如果前端允许创作者调整页的大小
	// 那么只有 100 这个页大小这个默认情况下，会走索引
	if offset == 0 && limit == 100 {
		data, err := repo.cache.GetFirstPage(ctx, author)
		if err == nil {
			go func() {
				repo.preCache(ctx, data)
			}()
			return data, nil
		}
		// 这里记录日志
		if err != cache.ErrKeyNotExist {
			repo.l.Error("查询缓存文章失败",
				logger.Int64("author", author), logger.Error(err))
		}
	}
	// 慢路径
	arts, err := repo.dao.GetByAuthor(ctx, author, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[article.Article, domain.Article](arts,
		func(idx int, src article.Article) domain.Article {
			return repo.toDomain(src)
		})
	// 一般都是让调用者来控制是否异步。
	go func() {
		repo.preCache(ctx, res)
	}()
	// 你这个也可以做成异步的
	err = repo.cache.SetFirstPage(ctx, author, res)
	if err != nil {
		repo.l.Error("刷新第一页文章的缓存失败",
			logger.Int64("author", author), logger.Error(err))
	}
	return res, nil
}

func (repo *CachedArticleRepository) preCache(ctx context.Context,
	arts []domain.Article) {
	// 1MB
	const contentSizeThreshold = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) <= contentSizeThreshold {
		// 你也可以记录日志
		if err := repo.cache.Set(ctx, arts[0]); err != nil {
			repo.l.Error("提前准备缓存失败", logger.Error(err))
		}
	}
}

func (repo *CachedArticleRepository) SyncStatus(ctx context.Context,
	uid, id int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, uid, id, status.ToUint8())
}

func (repo *CachedArticleRepository) Sync(ctx context.Context,
	art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	if err != nil {
		return 0, err
	}
	go func() {
		author := art.Author.Id
		err = repo.cache.DelFirstPage(ctx, author)
		if err != nil {
			repo.l.Error("删除第一页缓存失败",
				logger.Int64("author", author), logger.Error(err))
		}
		user, err := repo.userRepo.FindById(ctx, author)
		if err != nil {
			repo.l.Error("提前设置缓存准备用户信息失败",
				logger.Int64("uid", author), logger.Error(err))
		}
		art.Author = domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		}
		err = repo.cache.SetPub(ctx, art)
		if err != nil {
			repo.l.Error("提前设置缓存失败",
				logger.Int64("author", author), logger.Error(err))
		}
	}()
	return id, nil
}

func (repo *CachedArticleRepository) SyncV2(ctx context.Context,
	art domain.Article) (int64, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 直接 defer Rollback
	// 如果我们后续 Commit 了，这里会得到一个错误，但是没关系
	defer tx.Rollback()
	authorDAO := article.NewGORMArticleDAO(tx)
	readerDAO := article.NewGORMArticleReaderDAO(tx)

	// 下面代码和 SyncV1 一模一样
	artn := repo.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = authorDAO.Insert(ctx, artn)
		if err != nil {
			return 0, err
		}
	} else {
		err = authorDAO.UpdateById(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = readerDAO.UpsertV2(ctx, article.PublishedArticle(artn))
	if err != nil {
		// 依赖于 defer 来 rollback
		return 0, err
	}
	tx.Commit()
	return artn.Id, nil
}

func (repo *CachedArticleRepository) SyncV1(ctx context.Context,
	art domain.Article) (int64, error) {
	artn := repo.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = repo.authorDAO.Create(ctx, artn)
		if err != nil {
			return 0, err
		}
	} else {
		err = repo.authorDAO.UpdateById(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = repo.readerDAO.Upsert(ctx, artn)
	return id, err
}

func (repo *CachedArticleRepository) Create(ctx context.Context,
	art domain.Article) (int64, error) {
	id, err := repo.dao.Insert(ctx, repo.toEntity(art))
	if err != nil {
		return 0, err
	}
	author := art.Author.Id
	err = repo.cache.DelFirstPage(ctx, author)
	if err != nil {
		repo.l.Error("删除缓存失败",
			logger.Int64("author", author), logger.Error(err))
	}
	return id, nil
}

func (repo *CachedArticleRepository) Update(ctx context.Context,
	art domain.Article) error {
	err := repo.dao.UpdateById(ctx, repo.toEntity(art))
	if err != nil {
		return err
	}
	author := art.Author.Id
	err = repo.cache.DelFirstPage(ctx, author)
	if err != nil {
		repo.l.Error("删除缓存失败",
			logger.Int64("author", author), logger.Error(err))
	}
	return nil
}

func (repo *CachedArticleRepository) toDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
	}
}

func (repo *CachedArticleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		// 这一步，就是将领域状态转化为存储状态。
		// 这里我们就是直接转换，
		// 有些情况下，这里可能是借助一个 map 来转
		Status: uint8(art.Status),
	}
}
