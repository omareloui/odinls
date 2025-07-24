package product

type (
	RetrieveOptsFunc func(*RetrieveOpts)
	RetrieveOpts     struct {
		PopulateUsedMaterial bool
	}
)

func WithPopulatedUserMaterial(opts *RetrieveOpts) {
	opts.PopulateUsedMaterial = true
}

func ParseRetrieveOpts(funcs ...RetrieveOptsFunc) *RetrieveOpts {
	o := &RetrieveOpts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
