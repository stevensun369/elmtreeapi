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

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
)

// parent middleware
func parentMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.ParentClaims{}
    _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    if err != nil {
      return utils.Error(c, err)
    }

    // if !tkn.Valid {
    //   return utils.MessageError(c, "token not valid")
    // }

    utils.SetLocals(c, "parentID", claims.ParentID)
    utils.SetLocals(c, "studentIDList", claims.StudentIDList)
  }

  if (token == "") {
    return utils.MessageError(c, "no token provided")
  }

  return c.Next()
}

// @desc   Register parent
// @route  POST /api/parent/register
// @access Public
func parentRegister(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // generating the parent id
  var parentID = utils.GenID()
  parentID = utils.GenID()
  var parentGenID models.Parent
  db.Parents.FindOne(context.Background(), bson.M{"parentID": parentID}).Decode(&parentGenID)
  for (parentGenID.ParentID != "") {
    parentID = utils.GenID()
    db.Parents.FindOne(context.Background(), bson.M{"parentID": parentID}).Decode(&parentGenID)
  } 

  // check if there is a parent account with cnp
  var checkParent models.Parent
  db.Parents.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&checkParent)
  if (checkParent.ParentID != "") {
    return utils.MessageError(c, "Există deja un părinte cu CNP-ul introdus.")
  }

  // hashed password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
  if err != nil {
    return utils.Error(c, err)
  }

  parent := models.Parent {
    ParentID: parentID,
    CNP: body["cnp"],
    FirstName: body["firstName"],
    LastName: body["lastName"],
    Password: string(hashedPassword),
    StudentIDList: []string {},
  }

  insertedResult, err := db.Parents.InsertOne(context.Background(), parent)
  if err != nil {
    return utils.Error(c, err)
  }

  // generate token
  tokenString, err := utils.ParentGenerateToken(parent.ParentID, []string {})
  if err != nil {
    return utils.Error(c, err)
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

  // getting the parent
  var parent models.Parent
  if err := db.Parents.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&parent); err != nil {
    return utils.MessageError(c, "Nu există niciun părinte cu CNP-ul introdus.")  
  }
  hashedPassword := parent.Password

  // get students
  students, err := db.GetStudents(bson.M{
    "studentID": bson.M{"$in": parent.StudentIDList},
  }, db.GradeSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(students) == 0 {
    students = []models.Student {}
  }

  compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

  tokenString, err := utils.ParentGenerateToken(parent.ParentID, parent.StudentIDList)
  if err != nil {
    return utils.Error(c, err)  
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
    return utils.MessageError(c, "Nu ați introdus parola validă.") 
  }
} 