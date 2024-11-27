// Copyright (c) 2024 Red Hat
// SPDX-License-Identifier: Apache-2.0

package logrustrace

import (
	"github.com/sirupsen/logrus"
)

var (
	defaultLogger = logrus.StandardLogger()
	defaultLevel  = logrus.DebugLevel
)

// config contains options for the STDOUT exporter.
type config struct {
	// Writer is the destination for the exported data.
	Logger logrus.FieldLogger

	// Level is the log level for the logger.
	Level logrus.Level
}

// newConfig creates a validated Config configured with options.
func newConfig(options ...Option) config {
	cfg := config{
		Logger: defaultLogger,
		Level:  defaultLevel,
	}
	for _, opt := range options {
		cfg = opt.apply(cfg)
	}
	return cfg
}

// Option sets the value of an option for a Config.
type Option interface {
	apply(config) config
}

// WithLogger sets the export stream destination.
func WithLogger(logger logrus.FieldLogger) Option {
	return writerOption{logger}
}

type writerOption struct {
	L logrus.FieldLogger
}

func (o writerOption) apply(cfg config) config {
	cfg.Logger = cfg.Logger
	return cfg
}

// WithLevel sets the log level for the logger.
func WithLevel(level logrus.Level) Option {
	return levelOption{level}
}

type levelOption struct {
	L logrus.Level
}

func (o levelOption) apply(cfg config) config {
	cfg.Level = o.L
	return cfg
}
