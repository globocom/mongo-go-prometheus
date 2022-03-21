package mongoprom_test

import (
	"context"
	"strings"
	"testing"
	"time"

	mongoprom "github.com/globocom/mongo-go-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
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
			mongoprom.WithInstanceName("testdb"),
			mongoprom.WithDurationBuckets([]float64{.1, .2}),
		)

		evt := &event.CommandSucceededEvent{
			CommandFinishedEvent: event.CommandFinishedEvent{
				CommandName:   "insert",
				DurationNanos: 200 * time.Millisecond.Nanoseconds(),
			},
		}

		// act
		sut.Succeeded(context.Background(), evt)

		// assert
		expected := strings.NewReader(`
			# HELP namespace1_mongo_commands Histogram of MongoDB commands
			# TYPE namespace1_mongo_commands histogram
			namespace1_mongo_commands_bucket{command="insert",instance="testdb",le="0.1"} 0
			namespace1_mongo_commands_bucket{command="insert",instance="testdb",le="0.2"} 1
			namespace1_mongo_commands_bucket{command="insert",instance="testdb",le="+Inf"} 1
			namespace1_mongo_commands_sum{command="insert",instance="testdb"} 0.2
			namespace1_mongo_commands_count{command="insert",instance="testdb"} 1
		`)
		err := testutil.GatherAndCompare(prometheus.DefaultGatherer, expected, "namespace1_mongo_commands")
		if err != nil {
			assert.Fail(err.Error())
		}
	})

	t.Run("export metrics after a command fails", func(t *testing.T) {
		// arrange
		sut := mongoprom.NewCommandMonitor(
			mongoprom.WithNamespace("namespace2"),
			mongoprom.WithInstanceName("testdb"),
			mongoprom.WithDurationBuckets([]float64{.2, .3}),
		)

		evt := &event.CommandFailedEvent{
			CommandFinishedEvent: event.CommandFinishedEvent{
				CommandName:   "update",
				DurationNanos: 300 * time.Millisecond.Nanoseconds(),
			},
		}

		// act
		sut.Failed(context.Background(), evt)

		// assert
		expected := strings.NewReader(`
			# HELP namespace2_mongo_commands Histogram of MongoDB commands
			# TYPE namespace2_mongo_commands histogram
			namespace2_mongo_commands_bucket{command="update",instance="testdb",le="0.2"} 0
			namespace2_mongo_commands_bucket{command="update",instance="testdb",le="0.3"} 1
			namespace2_mongo_commands_bucket{command="update",instance="testdb",le="+Inf"} 1
			namespace2_mongo_commands_sum{command="update",instance="testdb"} 0.3
			namespace2_mongo_commands_count{command="update",instance="testdb"} 1
		`)
		err := testutil.GatherAndCompare(prometheus.DefaultGatherer, expected, "namespace2_mongo_commands")
		if err != nil {
			assert.Fail(err.Error())
		}

		expected = strings.NewReader(`
			# HELP namespace2_mongo_command_errors Number of MongoDB commands that have failed
			# TYPE namespace2_mongo_command_errors counter
			namespace2_mongo_command_errors{command="update",instance="testdb"} 1
		`)
		err = testutil.GatherAndCompare(prometheus.DefaultGatherer, expected, "namespace2_mongo_command_errors")
		if err != nil {
			assert.Fail(err.Error())
		}
	})
}
