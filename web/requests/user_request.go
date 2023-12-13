package requests

type RegisterBody struct {
	Username        string `validate:"required"`
	Email           string `validate:"required,email"`
	Phone           string `validate:"required,e164"`
	Password        string `validate:"required,password"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}

type LoginBody struct {
	Username string `validate:"omitempty"`
	Email    string `validate:"omitempty,email"`
	Phone    string `validate:"omitempty,e164"`
	Password string `validate:"required,password"`
}

type UpdateProfileBody struct {
	Username string `validate:"omitempty"`
	Email    string `validate:"omitempty,email"`
	Phone    string `validate:"omitempty,e164"`
}

type UpdatePasswordBody struct {
	CurrentPassword string `validate:"required,password"`
	Password        string `validate:"required,password"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}

type ForgetPasswordBody struct {
	Email string `validate:"required,email"`
}

type ResetPasswordBody struct {
	Password        string `validate:"required,password"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}
