package user

type (
	RetrieveOptsFunc func(*RetrieveOpts)
	RetrieveOpts     struct {
		PopulateRole bool
	}
)

func WithPopulatedRole(opts *RetrieveOpts) {
	opts.PopulateRole = true
}

func ParseRetrieveOpts(funcs ...RetrieveOptsFunc) *RetrieveOpts {
	o := &RetrieveOpts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
