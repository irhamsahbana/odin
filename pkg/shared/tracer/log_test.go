package tracer

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func ExampleTraceContextHook() {
	handleRequest := func(w http.ResponseWriter, req *http.Request) {
		logger := zerolog.Ctx(req.Context()).Hook(TraceContextHook(req.Context()))
		logger.Error().Msg("message")
	}
	http.HandleFunc("/", handleRequest)
}

func TestTraceContextHookNothing(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Hook(TraceContextHook(context.Background()))
	logger.Info().Msg("message")

	require.Equal(t, "{\"level\":\"info\",\"message\":\"message\"}\n", buf.String())
}
