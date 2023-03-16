package adapters

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// PokemonAPI is data of instances.
type PokemonAPI struct {
	URL    string
	driver *resty.Client
}

// Open is open the connection of Pokémon.
func (p *PokemonAPI) Open() (*resty.Client, error) {
	if p.driver == nil {
		return nil, fmt.Errorf("driver was failed to connected")
	}
	return p.driver, nil
}

// Connect is connected the connection of Pokémon.
func (p *PokemonAPI) Connect() (err error) {
	p.driver = resty.New().SetBaseURL(p.URL)
	return nil
}

// Disconnect is disconnect the connection of Pokémon.
func (*PokemonAPI) Disconnect() error {
	return nil
}

// WithPokemon option function to assign on adapters.
func WithPokemon(driver Driver[*resty.Client]) Option {
	return func(a *Adapter) {
		if err := driver.Connect(); err != nil {
			panic(err)
		}
		open, err := driver.Open()
		if err != nil {
			panic(err)
		}
		a.PokemonRest = open
	}
}
