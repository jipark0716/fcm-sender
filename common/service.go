package common

import (
	"context"
)

type Service interface {
	Run() error
	Shutdown(ctx context.Context) error
}
