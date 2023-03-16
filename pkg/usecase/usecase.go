package usecase

import (
	"fmt"
	"reflect"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/adapters"
)

// Get returns the usecase of type T, creating it if necessary.
func Get[T any](adapter *adapters.Adapter) (T, error) {
	var zero T
	inf := reflect.TypeOf(&zero).Elem()
	reg, err := FindByType(inf)
	if err != nil {
		return zero, err
	}
	obj := reg.New()
	// Call Init if available.
	if i, ok := obj.(interface {
		Init(*adapters.Adapter) error
	}); ok {
		if err := i.Init(adapter); err != nil {
			return zero, fmt.Errorf("component %q initialization failed: %w", reg.Name, err)
		}
	}
	return obj.(T), nil
}
