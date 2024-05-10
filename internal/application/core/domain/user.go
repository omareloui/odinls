package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Name struct {
	First string `bson:"first"`
	Last  string `bson:"last"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      Name               `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Phone     string             `bson:"phone,omitempty"`
	Role      primitive.ObjectID `bson:"role,omitempty"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
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
