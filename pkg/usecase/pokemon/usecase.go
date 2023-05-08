// Package pokemon implement all logic.
package pokemon

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/entity"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/shared/tracer"
)

// GetAll returns resource pokemon api.
func (i *impl) GetAll(ctx context.Context) (*entity.Resource, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

	pok := &entity.Resource{}
	client := i.adapter.PokemonRest
	resp, err := client.R().SetResult(pok).Get("pokemon/")
	if err != nil {
		l.Error().Err(err).Msg("GetAll")
		return nil, err
	}

	l.Info().Int("Status Code", resp.StatusCode()).Msg("fetching pokemon")
	return pok, nil
}
