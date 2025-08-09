package order

type (
	RetrieveOptsFunc func(*RetrieveOpts)
	RetrieveOpts     struct {
		PopulateCraftsmen    bool
		PopulateClient       bool
		PopulateItemProducts bool
	}
)

func WithPopulatedCraftsman(opts *RetrieveOpts) {
	opts.PopulateCraftsmen = true
}

func WithPopulatedClient(opts *RetrieveOpts) {
	opts.PopulateClient = true
}

func WithPopulatedItemProducts(opts *RetrieveOpts) {
	opts.PopulateItemProducts = true
}

func ParseRetrieveOpts(funcs ...RetrieveOptsFunc) *RetrieveOpts {
	o := &RetrieveOpts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
