package util

import (
	"errors"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"github.com/golang-jwt/jwt"
	uploader "github.com/mahdimehrabi/uploader/minio"
	"strings"
)

func DecodeToken(tokenString string, secret string) (bool, jwt.MapClaims, error) {

	Claims := jwt.MapClaims{}

	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrUnexpectedSigningMethod := errors.New("unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, Claims, key)
	var valid bool
	if token == nil {
		valid = false
	} else {
		valid = token.Valid
	}

	return valid, Claims, err
}

func GeneratePublicURL(file string, env *godotenv.Env) string {
	strs := strings.Split(file, "/")
	if len(strs) == 2 {
		scheme := "http://"
		if env.Environment == "production" {
			scheme = "https://"
		}
		return uploader.GeneratePublicURL(scheme, env.MinioHost, strs[0], strs[1])
	}
	return ""
}
