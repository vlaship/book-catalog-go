package service

// Services is a struct for services
type Services struct {
	BookService     *BookService
	AuthorService   *AuthorService
	AuthService     *AuthService
	SendMailService *SendMailService
	UserService     *UserService
	OTPService      *OTPService
	PasswordService *PasswordService
	TosService      *TosService
}
