package client

type (
	RetrieveOptsFunc func(*RetrieveOpts)
	RetrieveOpts     struct{}
)

func ParseRetrieveOpts(funcs ...RetrieveOptsFunc) *RetrieveOpts {
	o := &RetrieveOpts{}
	for _, fun := range funcs {
		fun(o)
	}
	return o
}
