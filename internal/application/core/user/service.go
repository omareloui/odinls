package user

type UserService interface {
	GetUsers(opts ...RetrieveOptsFunc) ([]User, error)
	GetUserByID(id string, opts ...RetrieveOptsFunc) (*User, error)
	GetUserByEmailOrUsername(emailOrUsername string, opts ...RetrieveOptsFunc) (*User, error)
	GetUserByEmailOrUsernameFromUser(usr *User, opts ...RetrieveOptsFunc) (*User, error)
	CreateUser(user *User, opts ...RetrieveOptsFunc) (*User, error)
	UpdateUserByID(id string, user *User, opts ...RetrieveOptsFunc) (*User, error)
	UnsetCraftsmanByID(id string) (*User, error)
}
