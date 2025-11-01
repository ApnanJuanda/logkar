package middleware

import (
	"bsnack/domain/api/account/model"
	"bsnack/lib/env"
	"bsnack/lib/response"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func WithAuh() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}

		auths := strings.Split(authHeader, " ")
		if len(auths) != 2 {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}
		data, err := DecryptJWT(auths[1])
		fmt.Println(data)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}
		authAccount := model.Account{
			Id:       data["id"].(string),
			Email:    data["email"].(string),
			Name:     data["name"].(string),
			IsSeller: data["is_seller"].(bool),
		}
		ctxUser := context.WithValue(ctx.Request.Context(), "auth_account", authAccount)
		ctx.Request = ctx.Request.WithContext(ctxUser)
		ctx.Next()
	}
}

func DecryptJWT(token string) (map[string]interface{}, error) {
	var signature = []byte(env.String("TOKEN_SIGNATURE", "myPrivateSignature"))
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return signature, nil
	})
	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}

	if !parsedToken.Valid {
		return data, errors.New("invalid token")
	}
	fmt.Println("Decript JWT", parsedToken.Claims.(jwt.MapClaims))
	return parsedToken.Claims.(jwt.MapClaims), nil
}
