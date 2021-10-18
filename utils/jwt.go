package utils

import (
	"time"

	"backend-go/models"

	"github.com/dgrijalva/jwt-go"
)

var JWTKey = []byte("123456")

type TeacherClaims struct {
  TeacherID string `json:"teacherID"`
  HomeroomGrade models.Grade `json:"homeroomGrade"`
  SubjectList []models.Subject `json:"subjectList"`
  jwt.StandardClaims
}

func TeacherGenerateToken(id string, homeroomGrade models.Grade, subjectList []models.Subject) (tokenString string, err error) {
  // one year has 8760 hours
  expirationTime := time.Now().Add(8760 * time.Hour)

  // the "claims"
  claims := &TeacherClaims{
    TeacherID: id,
    HomeroomGrade: homeroomGrade,
    SubjectList: subjectList,
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

