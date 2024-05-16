package user

type UserService interface {
	FindUser(id string) (*User, error)
	FindUserByEmailOrUsername(emailOrUsername string) (*User, error)
	FindUserByEmailOrUsernameFromUser(*User) (*User, error)
	CreateUser(user *User) error
}
