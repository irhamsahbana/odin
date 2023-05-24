// Package pokemon implement all logic.
package pokemon

import (
	"context"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/entity"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/shared/tracer"
)

// GetAll returns resource pokemon api.
func (i *impl) GetAll(ctx context.Context) (*entity.Resource, error) {
	_, span, l := tracer.StartSpanLogTrace(ctx, "{{$Handler}}")
	defer span.End()

	pok := &entity.Resource{}
	client := i.adapter.PokemonResty
	resp, err := client.R().SetResult(pok).Get("pokemon/")
	if err != nil {
		l.Error().Err(err).Msg("GetAll")
		return nil, err
	}

	l.Info().Int("Status Code", resp.StatusCode()).Msg("fetching pokemon")
	return pok, nil
}
