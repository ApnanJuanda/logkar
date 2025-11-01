package encrypt

import (
	"bsnack/domain/api/account/model"
	"bsnack/lib/env"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GenerateFromPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encryptedPassword), nil
}

func CompareHashAndPassword(hashedPassword, password *string) error {
	return bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(*password))
}

func GenerateTokenLogin(account model.Account) (string, error) {
	var signature = []byte(env.String("TOKEN_SIGNATURE", "myPrivateSignature"))
	claims := jwt.MapClaims{
		"id":        account.Id,
		"email":     account.Email,
		"name":      account.Name,
		"is_seller": account.IsSeller,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iss":       "logkar",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(signature)
	return stringToken, err
}
