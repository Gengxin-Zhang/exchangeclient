package public

type Opt func(*Options)

type Register interface {
	OnDepth(interface{})
	OnCandle(interface{})
	OnTicker(interface{})
	OnTrade(interface{})
}

type Options struct {
	UseChannel bool
	Registers  []Register
}

func (o *Options) Use(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithChannel(useChannel bool) Opt {
	return func(o *Options) {
		o.UseChannel = useChannel
	}
}

func WithRegister(registers ...Register) Opt {
	return func(o *Options) {
		o.Registers = append(o.Registers, registers...)
	}
}
