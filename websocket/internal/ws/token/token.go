package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"upbit/internal/domain"
)

func CreateToken(token domain.Token) string {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_key": token.AccessKey,
		"nonce":      uuid.New().String(),
	})

	jwtToken, err := tkn.SignedString([]byte(token.SecretKey))
	if err != nil {
		log.Fatalf("Error caused by Token Creation process", zap.Error(err))
	}
	return jwtToken
}
