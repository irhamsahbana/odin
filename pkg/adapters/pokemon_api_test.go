package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPokemon(t *testing.T) {
	is := assert.New(t)

	adapter := &Adapter{}
	adapter.Sync(
		WithPokemon(&PokemonAPI{URL: "https://pokeapi.co/api/v2/"}),
	)
	resp, err := adapter.PokemonRest.R().Get("pokemon/")
	is.Nil(err)
	is.NotNil(resp)

	// Asserts
	is.Nil(adapter.UnSync())
}
