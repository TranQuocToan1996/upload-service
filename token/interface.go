package tokens

type TokenProvider interface {
	CreateToken(userClaims *UserClaims) (string, error)
	VerifyToken(token string) (*UserClaims, error)
}
