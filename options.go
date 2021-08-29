package mongoprom

type (
	Options struct{
		InstanceName string
		Namespace string
		DurationBuckets []float64
	}
	Option  func(*Options)
)

func DefaultOptions() *Options {
	return &Options{
		InstanceName: "unnamed",
		Namespace: "",
		DurationBuckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	}
}

func (options *Options) Merge(opts ...Option) {
	for _, opt := range opts {
		opt(options)
	}
}

func WithInstanceName(name string) Option {
	return func(options *Options) {
		options.InstanceName = name
	}
}
func WithNamespace(namespace string) Option {
	return func(options *Options) {
		options.Namespace = namespace
	}
}

func WithDurationBuckets(buckets []float64) Option {
	return func(options *Options) {
		options.DurationBuckets = buckets
	}
}
