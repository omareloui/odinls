package domain

type AuthName struct {
	First string `validate:"required"`
	Last  string `validate:"required,gte=3"`
}

type Register struct {
	Name            AuthName `validate:"required"`
	Email           string   `validate:"required,email"`
	Password        string   `validate:"required,gte=8,eqfield=ConfirmPassword"`
	ConfirmPassword string   `validate:"required,eqfield=Password"`
}

type Login struct {
	Email    string
	Password string
}

func NewLogin(email, password string) *Login {
	return &Login{Email: email, Password: password}
}

func NewRegister(name AuthName, email, password, confirmPassword string) *Register {
	return &Register{Name: name, Email: email, Password: password, ConfirmPassword: confirmPassword}
}
