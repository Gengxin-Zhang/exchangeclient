package wsclient

type Options struct {
	Signature func([]byte) []byte
	Inflact   func([]byte) ([]byte, error)
	Host      string
}

type Opt func(*Options)

func NewOptions(opts ...Opt) *Options {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithInflate(inflact func([]byte) ([]byte, error)) Opt {
	return func(o *Options) {
		o.Inflact = inflact
	}
}

func WithSignature(signature func([]byte) []byte) Opt {
	return func(o *Options) {
		o.Signature = signature
	}
}

func WithHost(host string) Opt {
	return func(o *Options) {
		o.Host = host
	}
}
