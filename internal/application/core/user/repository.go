package user

type UserRepository interface {
	GetUsers() ([]User, error)
	GetUser(id string) (*User, error)
	GetUserByEmailOrUsername(emailOrUsername string) (*User, error)
	GetUserByEmailOrUsernameFromUser(usr *User) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUserByID(id string, user *User) (*User, error)
	UnsetCraftsmanByID(id string) (*User, error)
}
