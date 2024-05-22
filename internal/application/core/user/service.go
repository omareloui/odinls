package user

type UserService interface {
	GetUsers(opts ...RetrieveOptsFunc) ([]User, error)
	FindUser(id string, opts ...RetrieveOptsFunc) (*User, error)
	FindUserByEmailOrUsername(emailOrUsername string, opts ...RetrieveOptsFunc) (*User, error)
	FindUserByEmailOrUsernameFromUser(usr *User, opts ...RetrieveOptsFunc) (*User, error)
	CreateUser(user *User, opts ...RetrieveOptsFunc) error
	UpdateUserByID(id string, user *User, opts ...RetrieveOptsFunc) error
}
