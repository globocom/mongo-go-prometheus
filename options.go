package mongoprom

type (
	// Options represents options to customize the exported metrics.
	Options struct {
		InstanceName    string
		Namespace       string
		DurationBuckets []float64
	}
	Option func(*Options)
)

type (
	// PoolOptions represents options to customize the exported metrics.
	PoolOptions struct {
		InstanceName string
		Namespace    string
	}
	PoolOption func(*PoolOptions)
)

// DefaultPoolOptions returns the default options.
func DefaultPoolOptions() *PoolOptions {
	return &PoolOptions{
		InstanceName: "unnamed",
		Namespace:    "default",
	}
}

func (options *PoolOptions) Merge(opts ...PoolOption) {
	for _, opt := range opts {
		opt(options)
	}
}

// DefaultOptions returns the default options.
func DefaultOptions() *Options {
	return &Options{
		InstanceName:    "unnamed",
		Namespace:       "",
		DurationBuckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	}
}

func (options *Options) Merge(opts ...Option) {
	for _, opt := range opts {
		opt(options)
	}
}

// PoolWithInstanceName sets the name of the MongoDB instance.
func PoolWithInstanceName(name string) PoolOption {
	return func(options *PoolOptions) {
		options.InstanceName = name
	}
}

// PoolWithNamespace sets the namespace of all metrics.
func PoolWithNamespace(namespace string) PoolOption {
	return func(options *PoolOptions) {
		options.Namespace = namespace
	}
}

// WithInstanceName sets the name of the MongoDB instance.
func WithInstanceName(name string) Option {
	return func(options *Options) {
		options.InstanceName = name
	}
}

// WithNamespace sets the namespace of all metrics.
func WithNamespace(namespace string) Option {
	return func(options *Options) {
		options.Namespace = namespace
	}
}

// WithDurationBuckets sets the duration buckets of commands.
func WithDurationBuckets(buckets []float64) Option {
	return func(options *Options) {
		options.DurationBuckets = buckets
	}
}
