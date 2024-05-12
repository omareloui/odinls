package domain

type AuthName struct {
	First string `validate:"required,min=3"`
	Last  string `validate:"required,min=3"`
}

type Register struct {
	Name            AuthName `validate:"required"`
	Email           string   `validate:"required,email,unique_email"`
	Password        string   `validate:"required,min=8,max=64,eqfield=ConfirmPassword"`
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
