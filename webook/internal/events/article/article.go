package article

import (
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/internal/events"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

const topicReadEvent = "article_read_event"

type ReadEvent struct {
	Aid int64
	Uid int64
}

type Producer interface {
	ProduceReadEvent(evt ReadEvent) error
}

var _ events.Consumer = &InteractiveReadEventConsumer{}

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveReadEventConsumer(
	client sarama.Client,
	l logger.LoggerV1,
	repo repository.InteractiveRepository) *InteractiveReadEventConsumer {
	ic := &InteractiveReadEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
	return ic
}

// Start 这边就是自己启动 goroutine 了
func (r *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicReadEvent},
			saramax.NewHandler[ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := r.repo.IncrReadCnt(ctx, "article", evt.Aid)
	return err
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}

func (s *SaramaSyncProducer) ProduceReadEvent(evt ReadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.
		SendMessage(&sarama.ProducerMessage{
			Topic: topicReadEvent,
			Key:   sarama.ByteEncoder(val),
		})
	return err
}
