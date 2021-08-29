package mongoprom_test

import (
	"testing"

	mongoprom "github.com/globocom/mongo-go-prometheus"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	assert := assert.New(t)

	t.Run("return default options", func(t *testing.T) {
		assert.Equal(&mongoprom.Options{
			InstanceName:    "unnamed",
			Namespace:       "",
			DurationBuckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		}, mongoprom.DefaultOptions())
	})

	t.Run("merge default options with custom ones", func(t *testing.T) {
		// arrange
		custom1 := func(options *mongoprom.Options) { options.Namespace = "custom" }
		custom2 := func(options *mongoprom.Options) { options.DurationBuckets = []float64{.1} }

		options := mongoprom.DefaultOptions()

		// act
		options.Merge(custom1, custom2)

		// assert
		assert.Equal(&mongoprom.Options{
			InstanceName:    "unnamed",
			Namespace:       "custom",
			DurationBuckets: []float64{.1},
		}, options)
	})

	t.Run("customize the name of the MongoDB instance", func(t *testing.T) {
		// arrange
		options := mongoprom.DefaultOptions()

		// act
		mongoprom.WithInstanceName("database")(options)

		// assert
		assert.Equal("database", options.InstanceName)
	})

	t.Run("customize metrics namespace", func(t *testing.T) {
		// arrange
		options := mongoprom.DefaultOptions()

		// act
		mongoprom.WithNamespace("custom")(options)

		// assert
		assert.Equal("custom", options.Namespace)
	})

	t.Run("customize metrics duration buckets", func(t *testing.T) {
		// arrange
		options := mongoprom.DefaultOptions()

		// act
		mongoprom.WithDurationBuckets([]float64{.01})(options)

		// assert
		assert.Equal([]float64{.01}, options.DurationBuckets)
	})
}
