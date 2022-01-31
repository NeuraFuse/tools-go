package jwt

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/neurafuse/tools-go/random"
	"time"
)

var signingKeyDefault = []byte("FAFWhdfli3209834z5hnAEFhklusefgli218AFESGliw3q9q3")
var SigningKeyActive = signingKeyDefault

func ResetSigningKey() {
	SigningKeyActive = signingKeyDefault
}

func GenerateToken(signingKey []byte) (string, error) {
	token := jwtgo.New(jwtgo.SigningMethodHS256)

	claims := token.Claims.(jwtgo.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func GenerateSigningKey() string {
	return random.GetString(64)
}
