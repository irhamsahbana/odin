package tracer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	// SpanIDFieldName is the field name for the span ID.
	SpanIDFieldName = "span.id"

	// SpanContext is the field name for the span context.
	SpanContext = "span.context"

	// TraceIDFieldName is the field name for the trace ID.
	TraceIDFieldName = "trace.id"
)

// TraceContextHook returns a zerolog.Hook that will add any trace context
// contained in ctx to log events.
func TraceContextHook(ctx context.Context) zerolog.Hook {
	return traceContextHook{ctx}
}

var carrier http.Header

type traceContextHook struct {
	ctx context.Context
}

func (t traceContextHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	carrier = http.Header{}
	sc := trace.SpanFromContext(t.ctx).SpanContext()
	if !sc.TraceID().IsValid() || !sc.SpanID().IsValid() {
		return
	}
	b, err := sc.MarshalJSON()
	if err != nil {
		return
	}
	e.Bytes(SpanContext, b)
	e.Str(TraceIDFieldName, sc.TraceID().String())
	e.Str(SpanIDFieldName, sc.SpanID().String())
	otel.GetTextMapPropagator().Inject(t.ctx, propagation.HeaderCarrier(carrier))
}

// ZeroWriter represents severity level.
type ZeroWriter struct {
	// MinLevel holds the minimum level of logs to send to
	//
	// MinLevel must be greater than or equal to zerolog.ErrorLevel.
	// If it is less than this, zerolog.ErrorLevel will be used as
	// the minimum instead.
	MinLevel zerolog.Level
}

func (w *ZeroWriter) minLevel() zerolog.Level {
	minLevel := w.MinLevel
	if minLevel < zerolog.DebugLevel {
		minLevel = zerolog.DebugLevel
	}
	return minLevel
}

// Write is a no-op.
func (*ZeroWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// WriteLevel decodes the JSON-encoded log record in p, and reports it as an error using w.Tracer.
func (w *ZeroWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level < w.minLevel() || level >= zerolog.NoLevel {
		return len(p), nil
	}

	var logRecord logRecord
	events, err := logRecord.decode(bytes.NewReader(p))
	if err != nil {
		return 0, err
	}

	if logRecord.spanContext.IsValid() {
		ctx := otel.GetTextMapPropagator().Extract(
			context.Background(), propagation.HeaderCarrier(carrier))
		tr := otel.Tracer("logger")
		_, span := tr.Start(ctx, fmt.Sprintf("log.%s", level.String()))
		defer span.End()
		attrs := make([]attribute.KeyValue, 0)
		for key, value := range events {
			if key == SpanContext {
				continue
			}
			if key == zerolog.LevelFieldName {
				w.levelTo(value.(string), span)
				attrs = append(attrs, attribute.String(key, value.(string)))
				buf := &bytes.Buffer{}
				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(true)
				if err := enc.Encode(events[zerolog.MessageFieldName]); err != nil {
					span.RecordError(err)
					continue
				}
				attrs = append(attrs, attribute.String(zerolog.MessageFieldName, buf.String()))
			}

			switch v := value.(type) {
			case string:
				attrs = append(attrs, attribute.String(key, v))
			case json.Number:
				attrs = append(attrs, attribute.String(key, fmt.Sprintf("%v", v)))
			default:
				b, err := json.Marshal(v)
				if err != nil {
					span.RecordError(err)
				} else {
					attrs = append(attrs, attribute.String(key, fmt.Sprintf("%v", b)))
				}
			}
		}
		span.AddEvent(level.String(), trace.WithAttributes(attrs...))
	}
	return len(p), nil
}

func (*ZeroWriter) levelTo(level string, span trace.Span) {
	lvl, _ := zerolog.ParseLevel(level)
	switch lvl {
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		span.SetStatus(codes.Error, "logging")
	}
}

type spanContext struct {
	TraceID    string
	SpanID     string
	TraceFlags string
	TraceState string
	Remote     bool
}

type logRecord struct {
	message         string
	timestamp       time.Time
	traceID, spanID string
	spanContext     trace.SpanContext
}

func (l *logRecord) decode(r io.Reader) (events map[string]interface{}, err error) {
	m := make(map[string]interface{})
	d := json.NewDecoder(r)
	d.UseNumber()
	if err := d.Decode(&m); err != nil {
		return events, err
	}

	l.message, _ = m[zerolog.MessageFieldName].(string)
	if strVal, ok := m[zerolog.TimestampFieldName].(string); ok {
		if t, err := time.Parse(zerolog.TimeFieldFormat, strVal); err == nil {
			l.timestamp = t.UTC()
		}
	}
	if b, ok := m[SpanContext].(string); ok {
		var sc spanContext
		err := json.Unmarshal([]byte(b), &sc)
		if err != nil {
			return events, err
		}
		tid, err := trace.TraceIDFromHex(sc.TraceID)
		if err != nil {
			return events, err
		}
		sid, err := trace.SpanIDFromHex(sc.SpanID)
		if err != nil {
			return events, err
		}

		scc := trace.SpanContextConfig{
			TraceID:    tid,
			SpanID:     sid,
			TraceFlags: trace.FlagsSampled,
			TraceState: trace.TraceState{},
			Remote:     false,
		}
		if sc := trace.NewSpanContext(scc); sc.IsValid() {
			l.spanContext = sc
		}
	}

	if strVal, ok := m[SpanIDFieldName].(string); ok {
		l.spanID = strVal
	}

	if strVal, ok := m[TraceIDFieldName].(string); ok {
		l.traceID = strVal
	}
	return m, nil
}
