package passwords

type Helper interface {
	HashPassword(password string) string
	CreateJWTToken(userId int, email string) string
	ComparePasswords(currentPwHash string, testingPassword string) bool
	TokenGenerator() string
	ValidateToken(tokenHeader string) (int, string, error)
	ValidatePassword(password string) error
}
