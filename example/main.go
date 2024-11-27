// Copyright (c) 2024 Red Hat
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/osbuild/logrustrace"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("echo-server")

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	r := echo.New()
	r.Use(otelecho.Middleware("my-server"))

	r.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		name := getUser(c.Request().Context(), id)
		return c.JSON(http.StatusOK, struct {
			ID   string
			Name string
		}{
			ID:   id,
			Name: name,
		})
	})
	logrus.SetLevel(logrus.DebugLevel)
	_ = r.Start(":8080")
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := logrustrace.New()
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func getUser(ctx context.Context, id string) string {
	ctx, span := tracer.Start(ctx, "getUser", oteltrace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	return dbUser(ctx, id)
}

var ErrNotFound = echo.NewHTTPError(http.StatusNotFound, "user not found")

func dbUser(ctx context.Context, id string) string {
	_, span := tracer.Start(ctx, "dbUser", oteltrace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	if id == "1" {
		return "first"
	}

	span.RecordError(ErrNotFound)
	span.SetStatus(codes.Error, "user not found")
	return "not found"
}
