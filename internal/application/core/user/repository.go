package user

type UserRepository interface {
	FindUser(id string) (*User, error)
	FindUserByEmailOrUsername(emailOrUsername string) (*User, error)
	CreateUser(user *User) error
}