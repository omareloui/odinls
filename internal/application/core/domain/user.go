package domain

type Name struct {
	First string
	Last  string
}

type User struct {
	ID       ID
	Name     Name
	Email    string
	Password string
	Phone    string
	Role     ID
}

type Craftsman struct {
	ID         ID
	User       ID
	HourlyRate float64
	Merchant   ID
}

type Client struct {
	ID                 ID
	User               ID
	Merchants          []ID
	Notes              string
	WholesaleAsDefault bool
	Locations          []string
	ContactInfo        ContactInfo
}

type ContactInfo struct {
	PhoneNumbers map[string]string
	Emails       map[string]string
	Links        []string          // auto deduct if it's a social media link in FE and show the icon accordingly
	Locations    map[string]string // places you can fine the user at e.g. (home, work). more of a v2 thing
}
