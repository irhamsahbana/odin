// Package adapters are the glue between components and external sources.
// # This manifest was generated by ymir. DO NOT EDIT.
package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPokemonRest(t *testing.T) {
	is := assert.New(t)

	adapter := &Adapter{}
	adapter.Sync(
		WithPokemonRest(&PokemonRest{URL: "https://pokeapi.co/api/v2/"}),
	)
	resp, err := adapter.PokemonRest.Get("pokemon/")
	is.Nil(err)
	is.NotNil(resp)

	// Asserts
	is.Nil(adapter.UnSync())
}
