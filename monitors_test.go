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

		assert.ElementsMatch([]string{
			"namespace1_mongo_commands",
		}, filter(metrics, "namespace1"))
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

		assert.ElementsMatch([]string{
			"namespace1_mongo_commands",
			"namespace1_mongo_command_errors",
		}, filter(metrics, "namespace1"))
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
