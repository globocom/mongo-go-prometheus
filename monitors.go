package mongoprom

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/event"
)

var (
	labelNames = []string{"instance", "command"}
)

func NewCommandMonitor(opts ...Option) *event.CommandMonitor {
	options := DefaultOptions()
	options.Merge(opts...)

	commands := register(prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: options.Namespace,
		Name:      "mongo_commands",
		Help:      "Histogram of MongoDB commands",
		Buckets:   options.DurationBuckets,
	}, labelNames)).(*prometheus.HistogramVec)

	errors := register(prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: options.Namespace,
		Name:      "mongo_command_errors",
		Help:      "Number of MongoDB commands that have failed",
	}, labelNames)).(*prometheus.CounterVec)

	observeDuration := func(evt event.CommandFinishedEvent) {
		duration := time.Duration(evt.DurationNanos) / time.Second
		commands.WithLabelValues(options.InstanceName, evt.CommandName).Observe(float64(duration))
	}

	return &event.CommandMonitor{
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			observeDuration(evt.CommandFinishedEvent)
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			observeDuration(evt.CommandFinishedEvent)
			errors.WithLabelValues(options.InstanceName, evt.CommandName).Inc()
		},
	}
}

func register(collector prometheus.Collector) prometheus.Collector {
	err := prometheus.DefaultRegisterer.Register(collector)
	if err == nil {
		return collector
	}

	if arErr, ok := err.(prometheus.AlreadyRegisteredError); ok {
		return arErr.ExistingCollector
	}

	panic(err)
}
