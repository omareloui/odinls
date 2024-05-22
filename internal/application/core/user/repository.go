package user

type UserRepository interface {
	GetUsers() ([]User, error)
	FindUser(id string) (*User, error)
	FindUserByEmailOrUsername(emailOrUsername string) (*User, error)
	FindUserByEmailOrUsernameFromUser(*User) (*User, error)
	CreateUser(user *User) error
	UpdateUserByID(id string, user *User) error
}
