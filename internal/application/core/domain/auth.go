package domain

type Register struct {
	Name            Name
	Email           string
	Password        string
	ConfirmPassword string
}

type Login struct {
	Email    string
	Password string
}

func NewLogin(email, password string) *Login {
	return &Login{Email: email, Password: password}
}

func NewRegister(name Name, email, password, confirmPassword string) *Register {
	return &Register{Name: name, Email: email, Password: password, ConfirmPassword: confirmPassword}
}
