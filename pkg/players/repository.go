package players

import (
	"context"
	"errors"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, player Player) (*Player, error)
	List(ctx context.Context) ([]Player, error)
	Get(ctx context.Context, firstName string, lastName string) (*Player, error)
	Update(ctx context.Context, id int64, player Player) (*Player, error)
	Delete(ctx context.Context, id int64, player Player) error
}
