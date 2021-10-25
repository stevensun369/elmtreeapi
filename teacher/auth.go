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
	"fmt"
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

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if !tkn.Valid {
      return c.Status(500).SendString("token not valid")
    }

    teacherIDBytes, _ := json.Marshal(claims.TeacherID)
    teacherIDJSON := string(teacherIDBytes)
    c.Locals("teacherID", teacherIDJSON)

    homeroomBytes, _ := json.Marshal(claims.HomeroomGrade)
    homeroomJSON := string(homeroomBytes)
    c.Locals("homeroomGrade", homeroomJSON)

    subjectListBytes, _ := json.Marshal(claims.SubjectList)
    subjectListJSON := string(subjectListBytes)
    c.Locals("subjectList", subjectListJSON)
  }

  if (token == "") {
    return c.Status(500).SendString("no token")
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

  // getting the db collection
  collection, err := db.GetCollection("teachers")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the teacher
  var teacher models.Teacher
  if err = collection.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&teacher); err != nil {
    return c.Status(401).JSON(bson.M{
      "message": "Nu există niciun profesor cu CNP-ul introdus.",
    })
  }

  // if the teacher doesn't have a password
  if teacher.Password == "" {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    var modifiedTeacher models.Teacher
    collection.FindOneAndUpdate(context.Background(), bson.M{"cnp": body["cnp"]}, bson.D{
      {Key: "$set", Value: bson.D{{Key: "password",Value: string(hashedPassword)}}},
    }).Decode(&modifiedTeacher)

    
    // jwt
    tokenString, err := utils.TeacherGenerateToken(modifiedTeacher.TeacherID, modifiedTeacher.HomeroomGrade, modifiedTeacher.SubjectList)
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

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
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

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
      return c.Status(401).JSON(bson.M{
        "message": "Nu ați introdus parola validă.",
      })
    }
  } 
}