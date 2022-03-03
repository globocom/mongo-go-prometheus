package mongoprom

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/event"
)

var (
	monitorLabels     = []string{"instance", "command"}
	poolMonitorLabels = []string{"instance"}
)

// NewCommandMonitor creates a event.CommandMonitor that exports metrics of Mongo commands.
// It also registers Prometheus collectors.
//
// The following metrics are exported:
//
// - Histogram of command duration.
// - Counter of command errors.
func NewCommandMonitor(opts ...Option) *event.CommandMonitor {
	options := DefaultOptions()
	options.Merge(opts...)

	commands := register(prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: options.Namespace,
		Name:      "mongo_commands",
		Help:      "Histogram of MongoDB commands",
		Buckets:   options.DurationBuckets,
	}, monitorLabels)).(*prometheus.HistogramVec)

	errors := register(prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: options.Namespace,
		Name:      "mongo_command_errors",
		Help:      "Number of MongoDB commands that have failed",
	}, monitorLabels)).(*prometheus.CounterVec)

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

func NewPoolMonitor(opts ...PoolOption) *event.PoolMonitor {
	options := DefaultPoolOptions()
	options.Merge(opts...)

	pcr := &PoolConnectionsRepository{
		usedConnections: 0,
		mux:             sync.Mutex{},
	}

	poolUsage := register(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: options.Namespace,
			Name:      "mongodb_connection_pool_usage",
			Help:      "MongoDB Connection Pool Usage",
		}, poolMonitorLabels)).(*prometheus.GaugeVec)

	poolMin := register(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: options.Namespace,
			Name:      "mongodb_connection_pool_min",
			Help:      "MongoDB Connection Pool Minimum",
		}, poolMonitorLabels)).(*prometheus.GaugeVec)

	poolMax := register(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: options.Namespace,
			Name:      "mongodb_connection_pool_max",
			Help:      "MongoDB Connection Pool Maximum",
		}, poolMonitorLabels)).(*prometheus.GaugeVec)

	observePoolEvents := func(poolEvt *event.PoolEvent) {
		if poolEvt != nil {
			if event.ConnectionReturned == poolEvt.Type {
				pcr.Return()
			}
			if event.GetSucceeded == poolEvt.Type {
				pcr.Get()
			}
			if event.PoolCleared == poolEvt.Type || event.PoolCreated == poolEvt.Type {
				pcr.Clear()
			}
			poolUsage.WithLabelValues(options.InstanceName).Set(float64(pcr.usedConnections))
			if poolEvt.PoolOptions != nil {
				poolMin.WithLabelValues(options.InstanceName).Set(float64(poolEvt.PoolOptions.MinPoolSize))
				poolMax.WithLabelValues(options.InstanceName).Set(float64(poolEvt.PoolOptions.MaxPoolSize))
			}
		}
	}

	return &event.PoolMonitor{
		Event: func(poolEvent *event.PoolEvent) {
			observePoolEvents(poolEvent)
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

type PoolConnectionsRepository struct {
	usedConnections int
	mux             sync.Mutex
}

func (c *PoolConnectionsRepository) Get() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.usedConnections++
}

func (c *PoolConnectionsRepository) Clear() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.usedConnections = 0
}

func (c *PoolConnectionsRepository) Return() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.usedConnections--
}

func (c *PoolConnectionsRepository) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.usedConnections
}
