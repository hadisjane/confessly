package errs

import "errors"

var (
	ErrConfessionNotFound             = errors.New("confession not found")
	ErrConfessionInvalid              = errors.New("confession is invalid")
	ErrConfessionTextEmpty           = errors.New("confession text is empty")
	ErrInvalidId                = errors.New("invalid id")
	ErrNotFound                 = errors.New("not found")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrIncorrectUsernameOrPassword = errors.New("incorrect username or password")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrForbidden                = errors.New("you don't have permission to access this resource")
	ErrForbiddenDelete          = errors.New("you don't have permission to delete this confession")
	ErrReportExists             = errors.New("you have already reported this confession")
	ErrUserBanned               = errors.New("your account has been banned")
	ErrUserAlreadyBanned        = errors.New("user is already banned")
	ErrUserNotBanned            = errors.New("user is not banned")
	ErrYouCannotBanYourself     = errors.New("you cannot ban yourself")
	ErrYouCannotBanOtherAdmin   = errors.New("you cannot ban other administrators")
	ErrInternalServer           = errors.New("internal server error")
)