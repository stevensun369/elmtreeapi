package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTKey = []byte("123456")

type Claims struct {
  ID string `json:"id"`
  jwt.StandardClaims
}

func GenerateToken(id string) (tokenString string, err error) {
  // one year has 8760 hours
  expirationTime := time.Now().Add(8760 * time.Hour)

  // the "claims"
  claims := &Claims{
    ID: id,
    StandardClaims: jwt.StandardClaims{
      ExpiresAt: expirationTime.Unix(),
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err = token.SignedString(JWTKey)


  return tokenString, err
}

// func ValidateToken(token string) {
//   jwt.Parse()
// }

