package user

type UserService interface {
	GetUsers() ([]User, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmailOrUsername(emailOrUsername string) (*User, error)
	GetUserByEmailOrUsernameFromUser(usr *User) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUserByID(id string, user *User) (*User, error)
	UnsetCraftsmanByID(id string) (*User, error)
}
