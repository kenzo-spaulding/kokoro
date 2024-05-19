package traces

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	api "go.opentelemetry.io/otel/sdk/trace"
)

// TODO endpoint for pushing traces and whether to use stdouttrace
type Traces struct {
	Style string `env:"TRACES_EXPORTER" envDefault:"CONSOLE"`
}

func Init(ctx context.Context, config Traces) error {
	var exporter api.SpanExporter
	var err error

	switch strings.ToUpper(config.Style) {
	case "CONSOLE":
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	default:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	if err != nil {
		return err
	}

	bsp := api.NewBatchSpanProcessor(exporter)
	provider := api.NewTracerProvider(
		api.WithSampler(api.AlwaysSample()),
		api.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(provider)

	go func() {
		select {
		case <-ctx.Done():
			err = provider.Shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to shutdown trace provider")
			}
		}
	}()

	return nil
}
