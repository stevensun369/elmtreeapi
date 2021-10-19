package student

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std
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
func studentMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.StudentClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if !tkn.Valid {
      return c.Status(500).SendString("token not valid")
    }

    studentIDBytes, _ := json.Marshal(claims.StudentID)
    studentIDJSON := string(studentIDBytes)
    c.Locals("studentID", studentIDJSON)

    gradeBytes, _ := json.Marshal(claims.Grade)
    gradeJSON := string(gradeBytes)
    c.Locals("grade", gradeJSON)

    subjectListBytes, _ := json.Marshal(claims.SubjectList)
    subjectListJSON := string(subjectListBytes)
    c.Locals("subjectList", subjectListJSON)
  }

  if (token == "") {
    return c.Status(500).SendString("no token")
  }

  return c.Next()
}

// @desc   Login student
// @route  POST /api/student/login
// @access Public
func postLogin(c *fiber.Ctx) error {

  // getting body and unmarshalling it into a body variable
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // getting the db collection
  collection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the student
  var student models.Student
  if err = collection.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&student); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // if the student doesn't have a password
  if student.Password == "" {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    var modifiedStudent models.Student
    collection.FindOneAndUpdate(context.Background(), bson.M{"cnp": body["cnp"]}, bson.D{
      {Key: "$set", Value: bson.D{{Key: "password",Value: string(hashedPassword)}}},
    }).Decode(&modifiedStudent)

    
    // jwt
    tokenString, err := utils.StudentGenerateToken(modifiedStudent.StudentID, modifiedStudent.Grade, modifiedStudent.SubjectList)
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    return c.JSON(bson.M{
      "studentID": modifiedStudent.StudentID,
      "firstName": modifiedStudent.FirstName,
      "dadInitials": modifiedStudent.DadInitials,
      "lastName": modifiedStudent.LastName,
      "cnp": modifiedStudent.CNP,
      "password": modifiedStudent.Password,
      "grade": modifiedStudent.Grade,
      "subjectList": modifiedStudent.SubjectList,   
      "token": tokenString,
    })
  } else {
    hashedPassword := student.Password

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

    tokenString, err := utils.StudentGenerateToken(student.StudentID, student.Grade, student.SubjectList)
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if compareErr == nil {
      return c.JSON(bson.M{
        "studentID": student.StudentID,
        "firstName": student.FirstName,
        "dadInitials": student.DadInitials,
        "lastName": student.LastName,
        "cnp": student.CNP,
        "password": student.Password,
        "grade": student.Grade,
        "subjectList": student.SubjectList,
        "token": tokenString,
      })
    } else {
      return c.Status(500).SendString(fmt.Sprintf("%v", compareErr))
    }
  } 
}