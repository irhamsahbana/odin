// Package usecase is implements component logic.
package usecase

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// registry is a state for registered components.
type registry struct {
	m          sync.Mutex
	components map[reflect.Type]*Registration
	byName     map[string]*Registration
}

// Registration is the configuration needed to register a useCase.
type Registration struct {
	Name string       // full package-prefixed usecase name
	Inf  reflect.Type // interface type for the usecase
	New  func() any   // returns a new instance of the implementation type
}

// register registers a UseCase.
func (r *registry) register(reg Registration) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.components[reg.Inf]; ok {
		return errors.New("usecase already registered")
	}
	if r.components == nil {
		r.components = map[reflect.Type]*Registration{}
	}
	if r.byName == nil {
		r.byName = map[string]*Registration{}
	}
	ptr := &reg
	r.components[reg.Inf] = ptr
	r.byName[reg.Name] = ptr
	return nil
}

// allUseCases returns all the registered components.
func (r *registry) allUseCases() []*Registration {
	r.m.Lock()
	defer r.m.Unlock()

	useCases := make([]*Registration, 0, len(r.components))
	for _, info := range r.components {
		useCases = append(useCases, info)
	}
	return useCases
}

func (r *registry) findByType(t reflect.Type) (*Registration, bool) {
	r.m.Lock()
	defer r.m.Unlock()
	reg, ok := r.components[t]
	return reg, ok
}

// globalRegistry is the global registry used by Register and Registered.
var globalRegistry registry

// Register registers a UseCase.
func Register(reg Registration) {
	if err := globalRegistry.register(reg); err != nil {
		panic(err)
	}
}

// Registered returns the useCases registered with Register.
func Registered() []*Registration {
	return globalRegistry.allUseCases()
}

// FindByType returns the Registration component by type definition.
func FindByType(t reflect.Type) (*Registration, error) {
	c, ok := globalRegistry.findByType(t)
	if !ok {
		return nil, fmt.Errorf("component of type %v was not registered; maybe you forgot to modified init func", t)
	}
	return c, nil
}
