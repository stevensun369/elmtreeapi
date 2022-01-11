package teacher

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std packages
	"context"
	"encoding/json"
	"strings"

	// jwt
	"github.com/dgrijalva/jwt-go"

	// bcrypt
	"golang.org/x/crypto/bcrypt"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
)

// teacher auth middleware
func teacherMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.TeacherClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    utils.CheckError(c, err)

    if !tkn.Valid {
      utils.MessageError(c, "token not valid")
    }

    utils.SetLocals(c, "teacherID", claims.TeacherID)
    utils.SetLocals(c, "homeroomGrade", claims.HomeroomGrade)
    utils.SetLocals(c, "subjectList", claims.SubjectList)
  }

  if (token == "") {
    utils.MessageError(c, "no token")
  }

  return c.Next()
}

// @desc   Login teacher
// @route  POST /api/teacher/login
// @access Public
func postLogin(c *fiber.Ctx) error {

  // getting body and unmarshalling it into a body variable
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  teacher, err := db.GetTeacherByCNP(body["cnp"])
  utils.CheckMessageError(c, err, "Nu există niciun profesor cu CNP-ul introdus.")

  // if the teacher doesn't have a password
  if teacher.Password == "" {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
    utils.CheckError(c, err)

    var modifiedTeacher models.Teacher
    db.Teachers.FindOneAndUpdate(context.Background(), bson.M{"cnp": body["cnp"]}, bson.D{
      {Key: "$set", Value: bson.D{{Key: "password",Value: string(hashedPassword)}}},
    }).Decode(&modifiedTeacher)
    
    // jwt
    tokenString, err := utils.TeacherGenerateToken(modifiedTeacher.TeacherID, modifiedTeacher.HomeroomGrade, modifiedTeacher.SubjectList)
    utils.CheckError(c, err)

    return c.JSON(bson.M{
      "teacherID": modifiedTeacher.TeacherID,
      "firstName": modifiedTeacher.FirstName,
      "lastName": modifiedTeacher.LastName,
      "cnp": modifiedTeacher.CNP,
      "password": modifiedTeacher.Password,
      "homeroomGrade": modifiedTeacher.HomeroomGrade,
      "subjectList": modifiedTeacher.SubjectList,   
      "token": tokenString,
    })
  } else {
    hashedPassword := teacher.Password

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

    tokenString, err := utils.TeacherGenerateToken(teacher.TeacherID, teacher.HomeroomGrade, teacher.SubjectList)
    utils.CheckError(c, err)

    if compareErr == nil {
      return c.JSON(bson.M{
        "teacherID": teacher.TeacherID,
        "firstName": teacher.FirstName,
        "lastName": teacher.LastName,
        "cnp": teacher.CNP,
        "password": teacher.Password,
        "homeroomGrade": teacher.HomeroomGrade,
        "subjectList": teacher.SubjectList,
        "token": tokenString,
      })
    } else {
      return utils.MessageError(c, "Nu ați introdus parola validă.")
    }
  } 
}