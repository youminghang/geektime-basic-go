package connpool

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

var errUnknownPattern = errors.New("未知的双写 pattern")

type DoubleWritePool struct {
	pattern *atomicx.Value[string]
	src     gorm.ConnPool
	dst     gorm.ConnPool
}

func NewDoubleWritePool(srcDB *gorm.DB, dst *gorm.DB) *DoubleWritePool {
	return &DoubleWritePool{
		src:     srcDB.ConnPool,
		dst:     dst.ConnPool,
		pattern: atomicx.NewValueOf(patternSrcOnly)}
}

func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case patternSrcOnly:
		tx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			pattern: pattern,
			src:     tx,
		}, err
	case patternSrcFirst:
		return d.startTwoTx(d.src, d.dst, pattern, ctx, opts)
	case patternDstFirst:
		return d.startTwoTx(d.src, d.dst, pattern, ctx, opts)
	case patternDstOnly:
		tx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			pattern: pattern,
			src:     tx,
		}, err
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePool) startTwoTx(
	first gorm.ConnPool,
	second gorm.ConnPool,
	pattern string,
	ctx context.Context,
	opts *sql.TxOptions) (*DoubleWriteTx, error) {
	src, err := first.(gorm.TxBeginner).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	dst, err := second.(gorm.TxBeginner).BeginTx(ctx, opts)
	if err != nil {
		// 记录日志
		_ = src.Rollback()
	}
	return &DoubleWriteTx{src: src, dst: dst, pattern: pattern}, nil
}

func (d *DoubleWritePool) ChangePattern(pattern string) {
	d.pattern.Store(pattern)
}

func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// 禁用这个东西，因为我们没有办法创建出来 sql.Stmt 实例
	panic("implement me")
}

func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.dst.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				// 并且要通知修复数据
			}
		}
		return res, err
	case patternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.src.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				// 并且要通知修复数据
			}
		}
		return res, err
	case patternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern.Load() {
	case patternSrcFirst, patternSrcOnly:
		return d.src.QueryContext(ctx, query, args...)
	case patternDstFirst, patternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern.Load() {
	case patternSrcFirst, patternSrcOnly:
		return d.src.QueryRowContext(ctx, query, args...)
	case patternDstFirst, patternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		// 因为返回值里面咩有 error，只能 panic 掉
		panic(errUnknownPattern)
	}
}

type DoubleWriteTx struct {
	pattern string
	src     *sql.Tx
	dst     *sql.Tx
}

func (d *DoubleWriteTx) Commit() error {
	switch d.pattern {
	case patternSrcFirst:
		err := d.src.Commit()
		err1 := d.dst.Commit()
		if err1 != nil {
			// 记录日志
		}
		return err
	case patternSrcOnly:
		return d.src.Commit()
	case patternDstFirst:
		err := d.dst.Commit()
		err1 := d.src.Commit()
		if err1 != nil {
			// 记录日志
		}
		return err
	case patternDstOnly:
		return d.dst.Commit()
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWriteTx) Rollback() error {
	switch d.pattern {
	case patternSrcFirst:
		err := d.src.Rollback()
		err1 := d.dst.Rollback()
		if err1 != nil {
			// 记录日志
		}
		return err
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWriteTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("implement me")
}

func (d *DoubleWriteTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern {
	case patternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.dst.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				// 并且要通知修复数据
			}
		}
		return res, err
	case patternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil {
			_, err1 := d.src.ExecContext(ctx, query, args...)
			if err1 != nil {
				// 这边要记录日志
				// 并且要通知修复数据
			}
		}
		return res, err
	case patternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWriteTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern {
	case patternSrcFirst, patternSrcOnly:
		return d.src.QueryContext(ctx, query, args...)
	case patternDstFirst, patternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWriteTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern {
	case patternSrcFirst, patternSrcOnly:
		return d.src.QueryRowContext(ctx, query, args...)
	case patternDstFirst, patternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		// 因为返回值里面咩有 error，只能 panic 掉
		panic(errUnknownPattern)
	}
}

const (
	patternSrcOnly  = "src_only"
	patternSrcFirst = "src_first"
	patternDstFirst = "dst_first"
	patternDstOnly  = "dst_only"
)
