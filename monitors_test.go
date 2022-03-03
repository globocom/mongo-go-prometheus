package mongoprom_test

import (
	"context"
	"strings"
	"testing"

	mongoprom "github.com/globocom/mongo-go-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/event"
)

func TestCommandMonitor(t *testing.T) {
	assert := assert.New(t)

	t.Run("create a new monitor", func(t *testing.T) {
		// act
		sut := mongoprom.NewCommandMonitor()

		// assert
		assert.NotNil(sut)
	})

	t.Run("do not panic if metrics are already registered", func(t *testing.T) {
		// act
		_ = mongoprom.NewCommandMonitor()

		// assert
		assert.NotPanics(func() {
			_ = mongoprom.NewCommandMonitor()
		})
	})

	t.Run("export metrics after a command succeeds", func(t *testing.T) {
		// arrange
		sut := mongoprom.NewCommandMonitor(
			mongoprom.WithNamespace("namespace1"),
			mongoprom.WithDurationBuckets([]float64{.1, .2}),
		)

		evt := &event.CommandSucceededEvent{}

		// act
		sut.Succeeded(context.Background(), evt)

		// assert
		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)

		assert.Contains(filter(metrics, "namespace1"), "namespace1_mongo_commands")
	})

	t.Run("export metrics after a command fails", func(t *testing.T) {
		// arrange
		sut := mongoprom.NewCommandMonitor(
			mongoprom.WithNamespace("namespace1"),
			mongoprom.WithDurationBuckets([]float64{.1, .2}),
		)

		evt := &event.CommandFailedEvent{}

		// act
		sut.Failed(context.Background(), evt)

		// assert
		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)

		assert.ElementsMatch(filter(metrics, "namespace1"), []string{
			"namespace1_mongo_commands",
			"namespace1_mongo_command_errors",
		})
	})

	t.Run("create a new pool monitor", func(t *testing.T) {
		npm := mongoprom.NewPoolMonitor()
		assert.NotNil(npm)
	})

	t.Run("do not panic if pool metrics are already registered", func(t *testing.T) {
		// act
		_ = mongoprom.NewPoolMonitor()

		// assert
		assert.NotPanics(func() {
			_ = mongoprom.NewPoolMonitor()
		})
	})

	t.Run("check if metric mongodb_connection_pool_usage  was registered successfully after event.GetSucceeded ", func(t *testing.T) {
		// arrange
		sut := mongoprom.NewPoolMonitor(
			mongoprom.PoolWithNamespace("namespace1"),
		)

		evt := &event.PoolEvent{Type: event.GetSucceeded, PoolOptions: &event.MonitorPoolOptions{
			MaxPoolSize: 0,
			MinPoolSize: 100,
		}}

		// act
		sut.Event(evt)

		// assert
		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)

		assert.Contains(filter(metrics, "namespace1"), "namespace1_mongodb_connection_pool_usage")
	})

	t.Run("check if metric mongodb_connection_pool_min  was registered successfully after event.GetSucceeded ", func(t *testing.T) {
		// arrange
		sut := mongoprom.NewPoolMonitor(
			mongoprom.PoolWithNamespace("namespace1"),
		)

		evt := &event.PoolEvent{Type: event.GetSucceeded, PoolOptions: &event.MonitorPoolOptions{
			MaxPoolSize: 0,
			MinPoolSize: 100,
		}}

		// act
		sut.Event(evt)

		// assert
		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)

		assert.Contains(filter(metrics, "namespace1"), "namespace1_mongodb_connection_pool_min")
	})

	t.Run("check if metric mongodb_connection_pool_max  was registered successfully after event.GetSucceeded ", func(t *testing.T) {
		// arrange
		sut := mongoprom.NewPoolMonitor(
			mongoprom.PoolWithNamespace("namespace1"),
		)

		evt := &event.PoolEvent{Type: event.GetSucceeded, PoolOptions: &event.MonitorPoolOptions{
			MaxPoolSize: 0,
			MinPoolSize: 100,
		}}

		// act
		sut.Event(evt)

		// assert
		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)

		assert.Contains(filter(metrics, "namespace1"), "namespace1_mongodb_connection_pool_max")
	})

	t.Run("check pool connection repository cleared ", func(t *testing.T) {
		pcr := &mongoprom.PoolConnectionsRepository{}
		pcr.Clear()
		assert.Equal(0, pcr.Value())
	})

	t.Run("check pool connection repository increasing. Using connection from pool ", func(t *testing.T) {
		pcr := &mongoprom.PoolConnectionsRepository{}
		pcr.Get()
		assert.Equal(1, pcr.Value())
	})

	t.Run("check pool connection repository returning. Using connection from pool ", func(t *testing.T) {
		pcr := &mongoprom.PoolConnectionsRepository{}
		pcr.Get()
		pcr.Get()
		pcr.Get()
		pcr.Return()
		assert.Equal(2, pcr.Value())
	})

	t.Run("check pool connection repository clearing. Using connection from pool ", func(t *testing.T) {
		pcr := &mongoprom.PoolConnectionsRepository{}
		pcr.Get()
		pcr.Get()
		pcr.Get()
		pcr.Return()
		pcr.Clear()
		assert.Equal(0, pcr.Value())
	})

}

func filter(metrics []*dto.MetricFamily, namespace string) []string {
	var result []string
	for _, metric := range metrics {
		if strings.HasPrefix(*metric.Name, namespace) {
			result = append(result, *metric.Name)
		}
	}
	return result
}
