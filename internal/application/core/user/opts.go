package user

type (
	RetrieveOptsFunc func(*RetrieveOpts)
	RetrieveOpts     struct {
		PopulateRole     bool
		PopulateMerchant bool
	}
)

func WithPopulatedRole(opts *RetrieveOpts) {
	opts.PopulateRole = true
}

func WithPopulatedMerchant(opts *RetrieveOpts) {
	opts.PopulateMerchant = true
}

func ParseRetrieveOpts(funcs ...RetrieveOptsFunc) *RetrieveOpts {
	o := &RetrieveOpts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
