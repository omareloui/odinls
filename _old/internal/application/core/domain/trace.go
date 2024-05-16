package domain

type Trace struct {
	User            ID     // this could be empty as this could be the user's first product
	Slug            string // auto generated or custom
	UserEmail       string // the user email that he should sign in with (or fill the user id with if already signed)
	Merchant        ID
	Product         ID          // the product name or better yet, the product id but that's over engineering it
	ContactOverride ContactInfo // this could be in v2 or later
	IsLive          bool
}
