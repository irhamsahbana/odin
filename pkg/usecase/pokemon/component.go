package pokemon

import (
	"context"
	"reflect"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/adapters"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/entity"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/usecase"
)

func init() {
	usecase.Register(usecase.Registration{
		Name: "pokemon",
		Inf:  reflect.TypeOf((*T)(nil)).Elem(),
		New: func() any {
			return &impl{}
		},
	})
}

// T is the interface implemented by all Pokemon Component implementations.
type T interface {
	GetAll(ctx context.Context) (*entity.Resource, error)
}

type impl struct {
	adapter *adapters.Adapter
}

// Init initializes the execution of a process involved in a Pokemon Component usecase.
func (i *impl) Init(adt *adapters.Adapter) error {
	i.adapter = adt
	return nil
}
