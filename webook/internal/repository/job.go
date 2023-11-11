package repository

import "context"

type JobRepository interface {
	Preempt(ctx context.Context) []
}
