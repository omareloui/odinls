package chisrv

type opts struct {
	hasToNotBeSigned bool
	protected        bool
}

type optsFunc func(*opts)

func withProtection(opts *opts) {
	opts.protected = true
}

func withHasToNotBeSigned(opts *opts) {
	opts.hasToNotBeSigned = true
}

func parseOpts(funcs ...optsFunc) *opts {
	o := &opts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
