package parent

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	// internal packages
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std
	"context"
	"encoding/json"
	"fmt"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// parent middleware
func parentMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.ParentClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if !tkn.Valid {
      return c.Status(500).SendString("token not valid")
    }

    parentIDBytes, _ := json.Marshal(claims.ParentID)
    parentIDJSON := string(parentIDBytes)
    c.Locals("parentID", parentIDJSON)

    studentIDListBytes, _ := json.Marshal(claims.StudentIDList)
    studentIDListJSON := string(studentIDListBytes)
    c.Locals("studentIDList", studentIDListJSON)
  }

  if (token == "") {
    return c.Status(500).SendString("no token")
  }

  return c.Next()
}

// @desc   Register parent
// @route  POST /api/parent/register
// @access Public
func parentRegister(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  parentsCollection, err := db.GetCollection("parents")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  // generating the parent id
  var parentID = utils.GenID()
  parentID = utils.GenID()
  fmt.Println(parentID)
  var parentGenID models.Parent
  parentsCollection.FindOne(context.Background(), bson.M{"parentID": parentID}).Decode(&parentGenID)
  for (parentGenID.ParentID != "") {
    parentID = utils.GenID()
    parentsCollection.FindOne(context.Background(), bson.M{"parentID": parentID}).Decode(&parentGenID)
    fmt.Println(parentID)
  } 

  // check if there is a parent account with cnp
  var checkParent models.Parent
  parentsCollection.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&checkParent)
  if (checkParent.ParentID != "") {
    return c.Status(401).SendString("Există deja un părinte cu CNP-ul introdus.")
  }

  // hashed password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  parent := models.Parent {
    ParentID: parentID,
    CNP: body["cnp"],
    FirstName: body["firstName"],
    LastName: body["lastName"],
    Password: string(hashedPassword),
    StudentIDList: []string {},
  }

  insertedResult, err := parentsCollection.InsertOne(context.Background(), parent)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // generate token
  tokenString, err := utils.ParentGenerateToken(parent.ParentID, []string {})
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "parentID": parent.ParentID,
    "cnp": parent.CNP,
    "firstName": parent.FirstName,
    "lastName": parent.LastName,
    "students": []string {},
    "token": tokenString,
  })
}

// @desc   login parent
// @route  POST /api/parent/login
// @access Public
func parentLogin(c *fiber.Ctx) error {
  // getting body and unmarshalling it into a body variable
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // getting the db collection
  parentsCollection, err := db.GetCollection("parents")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the parent
  var parent models.Parent
  if err = parentsCollection.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&parent); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  hashedPassword := parent.Password

  // get the students
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var students []models.Student
  options := options.Find()
  options.SetSort(bson.D{{Key: "grade.gradeNumber", Value: 1}, {Key: "grade.gradeLetter", Value: 1}})
  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "studentID": bson.M{"$in": parent.StudentIDList},
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }


  compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

  tokenString, err := utils.ParentGenerateToken(parent.ParentID, parent.StudentIDList)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if compareErr == nil {
    return c.JSON(bson.M{
      "parentID": parent.ParentID,
      "firstName": parent.FirstName,
      "lastName": parent.LastName,
      "cnp": parent.CNP,
      "studentIDList": parent.StudentIDList,
      "students": students,
      "password": parent.Password,
      "token": tokenString,
    })
  } else {
    return c.Status(401).SendString("Nu ati introdus parola corecta")
  }
} 