package order

type (
	RetrieveOptsFunc func(*RetrieveOpts)
	RetrieveOpts     struct {
		PopulateMerchant     bool
		PopulateCraftsmen    bool
		PopulateClient       bool
		PopulateItemProducts bool
		PopulateItemVariants bool
	}
)

func WithPopulatedMerchant(opts *RetrieveOpts) {
	opts.PopulateMerchant = true
}

func WithPopulatedCraftsman(opts *RetrieveOpts) {
	opts.PopulateCraftsmen = true
}

func WithPopulatedClient(opts *RetrieveOpts) {
	opts.PopulateClient = true
}

func WithPopulatedItemProducts(opts *RetrieveOpts) {
	opts.PopulateItemProducts = true
}

func WithPopulatedItemVariants(opts *RetrieveOpts) {
	WithPopulatedItemProducts(opts)
	opts.PopulateItemVariants = true
}

func ParseRetrieveOpts(funcs ...RetrieveOptsFunc) *RetrieveOpts {
	o := &RetrieveOpts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
