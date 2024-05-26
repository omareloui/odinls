package product

type (
	RetrieveOptsFunc func(*RetrieveOpts)
	RetrieveOpts     struct {
		PopulateMerchant  bool
		PopulateCraftsman bool
	}
)

func WithPopulatedMerchant(opts *RetrieveOpts) {
	opts.PopulateMerchant = true
}

func WithPopulatedCraftsman(opts *RetrieveOpts) {
	opts.PopulateCraftsman = true
}

func ParseRetrieveOpts(funcs ...RetrieveOptsFunc) *RetrieveOpts {
	o := &RetrieveOpts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
