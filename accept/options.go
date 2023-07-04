package accept

type (
	Option func(o *options)

	options struct {
		address  string
		certFile string
		keyFile  string
		chanSize int
	}
)

func getDefaultOptions() options {
	return options{
		certFile: "",
		keyFile:  "",
		chanSize: 256,
	}
}

func WithTLS(certFile, keyFile string) Option {
	return func(o *options) {
		o.certFile, o.keyFile = certFile, keyFile
	}
}

func WithChanSize(size int) Option {
	return func(o *options) {
		o.chanSize = size
	}
}
