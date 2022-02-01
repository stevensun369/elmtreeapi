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
    _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    if err != nil {
      utils.Error(c, err)
    }

    // if !tkn.Valid {
    //   return c.Status(500).SendString("token not valid")
    // }

    utils.SetLocals(c, "studentID", claims.StudentID)
    utils.SetLocals(c, "grade", claims.Grade)
    utils.SetLocals(c, "subjectList", claims.SubjectList)
  }

  if (token == "") {
    return utils.MessageError(c, "no token")
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

  // getting student by studentID
  student, err := db.GetStudentByID(body["studentID"])
  if err != nil {
    return utils.MessageError(c, "Nu există niciun elev cu ID-ul introdus.")
  }

  // if the student doesn't have a password
  if student.Password == "" {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)

    if err != nil {
      return utils.Error(c, err)
    }

    var modifiedStudent models.Student
    db.Students.FindOneAndUpdate(context.Background(), bson.M{"studentID": body["studentID"]}, bson.D{
      {Key: "$set", Value: bson.D{{Key: "password",Value: string(hashedPassword)}}},
    }).Decode(&modifiedStudent)
    
    // jwt
    tokenString, err := utils.StudentGenerateToken(modifiedStudent.StudentID, modifiedStudent.Grade, modifiedStudent.SubjectList)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
      "studentID": modifiedStudent.StudentID,
      "firstName": modifiedStudent.FirstName,
      "dadInitials": modifiedStudent.DadInitials,
      "lastName": modifiedStudent.LastName,
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
      return utils.Error(c, err)
    }

    if compareErr == nil {
      return c.JSON(bson.M{
        "studentID": student.StudentID,
        "firstName": student.FirstName,
        "dadInitials": student.DadInitials,
        "lastName": student.LastName,
        "password": student.Password,
        "grade": student.Grade,
        "subjectList": student.SubjectList,
        "token": tokenString,
      })
    } else {
      return utils.MessageError(c, "Nu ați introdus parola validă.",)
    }
  } 
}

// @desc   update student
// @route  GET /api/student/update
// @access Private
func update(c *fiber.Ctx) error {
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    return utils.Error(c, err)
  }
  
  tokenString, err := utils.StudentGenerateToken(student.StudentID, student.Grade, student.SubjectList)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(bson.M{
    "subjectList": student.SubjectList,
    "grade": student.Grade,
    "token": tokenString,
  })
}
