// Copyright (c) 2024 Red Hat
// SPDX-License-Identifier: Apache-2.0

package logrustrace

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

var zeroTime time.Time

var _ trace.SpanExporter = &Exporter{}

// New creates an Exporter with the passed options.
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options...)

	return &Exporter{
		logger: cfg.Logger,
		level:  cfg.Level,
	}, nil
}

// Exporter is an implementation of trace.SpanSyncer that writes spans to logrus.
type Exporter struct {
	logger logrus.FieldLogger
	level  logrus.Level
	stop   atomic.Bool
}

// ExportSpans writes spans to logrus.
func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if len(spans) == 0 {
		return nil
	}

	for _, span := range spans {
		traceId := span.SpanContext().TraceID().String()
		spanId := span.SpanContext().SpanID().String()
		duration := span.EndTime().Sub(span.StartTime())
		name := span.Name()
		parentId := span.Parent().SpanID()
		statusCode := span.Status().Code
		statusMsg := span.Status().Description

		attrs := span.Attributes()
		fields := make(logrus.Fields, len(attrs))
		for _, kv := range attrs {
			key := string(kv.Key)
			if kv.Value.Type() == attribute.STRINGSLICE {
				fields[key] = kv.Value.AsStringSlice()
			}
			val := kv.Value.AsString()

			if val == "" {
				continue
			}

			if key == "msg" || key == "message" {
				key = "otel_msg"
			}

			fields[key] = val
		}

		lr := logrus.WithFields(fields).WithFields(logrus.Fields{
			"trace_id":  traceId,
			"span_id":   spanId,
			"duration":  duration,
			"name":      name,
			"parent_id": parentId,
		})

		if statusCode != 0 {
			lr = lr.WithField("status_code", statusCode)
		}
		if statusMsg != "" {
			lr = lr.WithField("status_msg", statusMsg)
		}
		lr.Logf(e.level, "Trace %s span %s name %s duration %.4fs", traceId, spanId, name, duration.Seconds())

		if err := ctx.Err(); err != nil || e.stop.Load() {
			return err
		}
	}
	return nil
}

// Shutdown is called to stop the exporter, it interrupts any possible span writing immediately.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stop.Store(true)

	return nil
}

// MarshalLog is the marshaling function used by the logging system to represent this Exporter.
func (e *Exporter) MarshalLog() interface{} {
	return struct {
		Type  string
		Level string
	}{
		Type:  "stdout",
		Level: e.level.String(),
	}
}
