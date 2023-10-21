package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const accessAudience = "access"

type TokenProvider interface {
	CreateToken(userClaims *UserClaims) string
	VerifyToken(token string) (*UserClaims, error)
}

func NewJWTProvider() TokenProvider {
	return &JWTProvider{
		key:      []byte(viper.GetString("KEY")),
		duration: viper.GetDuration("TOKEN_DURATION"),
		signAlgo: jwt.SigningMethodHS256,
	}
}

type JWTProvider struct {
	key      []byte
	duration time.Duration
	signAlgo jwt.SigningMethod
}

func (j *JWTProvider) CreateToken(userClaims *UserClaims) string {
	signAt := time.Now()
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: signAt.Add(j.duration)},
			Issuer:    "uploadService",
			Subject:   fmt.Sprint(userClaims.UserID),
			Audience:  jwt.ClaimStrings{accessAudience},
			NotBefore: &jwt.NumericDate{Time: signAt},
			IssuedAt:  &jwt.NumericDate{Time: signAt},
			ID:        uuid.NewString(),
		},
		UserClaims: *userClaims,
	}
	token := jwt.NewWithClaims(j.signAlgo, claims)
	jwtToken, err := token.SignedString(j.key)
	if err != nil {
		panic(err)
	}

	return jwtToken
}

func (j *JWTProvider) VerifyToken(token string) (*UserClaims, error) {
	claims := JWTClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != j.signAlgo.Alg() {
				return nil, errors.New("wrong algo: Get " + token.Method.Alg() + " expect " + j.signAlgo.Alg())
			}
			return j.key, nil
		})
	if err != nil {
		return nil, err
	}
	return &claims.UserClaims, nil
}
