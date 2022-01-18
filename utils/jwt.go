package utils

import (
	"time"

	"backend-go/models"

	"github.com/dgrijalva/jwt-go"
)

var JWTKey = []byte("123456")

// claims
type TeacherClaims struct {
  TeacherID string `json:"teacherID"`
  HomeroomGrade models.Grade `json:"homeroomGrade"`
  SubjectList []models.Subject `json:"subjectList"`
  SchoolID string `json:"schoolID"`
  jwt.StandardClaims
}

type StudentClaims struct {
  StudentID string `json:"studentID"`
  Grade models.Grade `json:"grade"`
  SubjectList []models.ShortSubject `json:"subjectList"`
  jwt.StandardClaims
}

type ParentClaims struct {
  ParentID string `json:"parentID"`
  StudentIDList []string `json:"studentIDList"`
  jwt.StandardClaims
}

// generation functions
func TeacherGenerateToken(id string, homeroomGrade models.Grade, subjectList []models.Subject, schoolID string) (tokenString string, err error) {
  // one year has 8760 hours
  expirationTime := time.Now().Add(8760 * time.Hour)

  // the "claims"
  claims := &TeacherClaims{
    TeacherID: id,
    HomeroomGrade: homeroomGrade,
    SubjectList: subjectList,
    SchoolID: schoolID,
    StandardClaims: jwt.StandardClaims{
      ExpiresAt: expirationTime.Unix(),
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err = token.SignedString(JWTKey)

  return tokenString, err
}

func StudentGenerateToken(id string, grade models.Grade, subjectList []models.ShortSubject) (tokenString string, err error) {
  // one year has 8760 hours 
  expirationTime := time.Now().Add(8760 * time.Hour)

  // the "claims"
  claims := &StudentClaims{
    StudentID: id,
    Grade: grade,
    SubjectList: subjectList,
    StandardClaims: jwt.StandardClaims{
      ExpiresAt: expirationTime.Unix(),
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err = token.SignedString(JWTKey)

  return tokenString, err
}

func ParentGenerateToken(id string, students []string) (tokenString string, err error) {
  // one year has 8760 hours 
  expirationTime := time.Now().Add(8760 * time.Hour)

  // the "claims"
  claims := &ParentClaims{
    ParentID: id,
    StudentIDList: students,
    StandardClaims: jwt.StandardClaims{
      ExpiresAt: expirationTime.Unix(),
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err = token.SignedString(JWTKey)

  return tokenString, err
}

