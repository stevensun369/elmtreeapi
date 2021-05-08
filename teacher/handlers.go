package teacher

import (
	// internal backend package
	"backend-go/db"
	"backend-go/models"
	"math"
	"strconv"
	"strings"

	// internal utils
	"backend-go/utils"

	"context"
	"encoding/json"
	"fmt"

	// fiber

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	// bcrypt
	"golang.org/x/crypto/bcrypt"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// teacher auth middleware
func teacherMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.Claims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if !tkn.Valid {
      return c.Status(500).SendString("token not valid")
    }

    c.Locals("userID", claims.ID)
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
  var teacher bson.M
  if err = collection.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&teacher); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // if there is no teacher, we 500
  if teacher == nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }


  // if the teacher doesn't have a password
  if teacher["password"] == "" {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    var modifiedTeacher bson.M
    collection.FindOneAndUpdate(context.Background(), bson.M{"cnp": body["cnp"]}, bson.D{
      {Key: "$set", Value: bson.D{{Key: "password",Value: string(hashedPassword)}}},
    }).Decode(&modifiedTeacher)
    
    // jwt
    tokenString, err := utils.GenerateToken(fmt.Sprintf("%v", modifiedTeacher["teacherID"]))
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    return c.JSON(bson.M{
      "teacherID": modifiedTeacher["teacherID"],
      "firstName": modifiedTeacher["firstName"],
      "lastName": modifiedTeacher["lastName"],
      "cnp": modifiedTeacher["cnp"],
      "password": modifiedTeacher["password"],
      "homeroomGrade": modifiedTeacher["homeroomGrade"],
      "subjectList": modifiedTeacher["subjectList"],   
      "token": tokenString,
    })
  } else {
    hashedPassword := fmt.Sprintf("%v", teacher["password"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

    tokenString, err := utils.GenerateToken(fmt.Sprintf("%v", teacher["teacherID"]))
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if compareErr == nil {
      return c.JSON(bson.M{
        "teacherID": teacher["teacherID"],
        "firstName": teacher["firstName"],
        "lastName": teacher["lastName"],
        "cnp": teacher["cnp"],
        "password": teacher["password"],
        "homeroomGrade": teacher["homeroomGrade"],
        "subjectList": teacher["subjectList"],
        "token": tokenString,
      })
    } else {
      return c.Status(500).SendString(fmt.Sprintf("%v", compareErr))
    }
  } 
}

// @desc   Get homeroom info
// @route  GET /api/teacher/homeroom/info
// @access Private
func getHomeroomInfo(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))

  // getting a teacher
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(teacher.HomeroomGrade)
}

// @desc   Get homeroom students
// @route  GET /api/teacher/homeroom
// @access Private
func getHomeroomStudents(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))

  // getting a teacher
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the homeroom students
  var students []models.Student
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "grade.gradeID": teacher.HomeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(students)
}

// @desc   Get subjects of a homeroom student
// @route  GET /api/teacher/homeroom/:studentID
// @access Private
func getHomeroomStudentSubjects(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  // getting the students
  var student models.Student
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  studentsCollection.FindOne(context.Background(), bson.M{
    "studentID": studentID,
    }).Decode(&student)

  return c.JSON(student.SubjectList)
}

// @desc   Get average marks for the students in the homeroom grade
// @route  GET /api/teacher/homeroom/average
// @access Private
func getHomeroomAverageMarks(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))

  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var averageMarks []models.AverageMark
  averageMarksCollection, err := db.GetCollection("averageMarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "grade.gradeID": teacher.HomeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(averageMarks)
}

// @desc   Create term mark for the students at term
// @route  POST /api/teacher/homeroom/term
// @access Private
func createHomeroomTermMark(c *fiber.Ctx) error {
  // getting studentID and term from body
  var body map[string]string
  json.Unmarshal(c.Body(), &body)
  studentID := body["studentID"]
  term := body["term"]
  termInt, _ := strconv.Atoi(term)

  // term Marks Collection
  termMarksCollection, err := db.GetCollection("termmarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // checking that there isn't already a termMark
  var termMarkTest models.TermMark
  termMarksCollection.FindOne(context.Background(), bson.M{
    "studentID": studentID,
    "term": termInt,
  }).Decode(&termMarkTest)

  if (termMarkTest != models.TermMark{Grade: models.Grade{}}) {
    return c.Status(500).SendString("Exista deja o medie pe semestrul " + term)
  }

  // get student and his subjectList
  student, err := models.GetStudentById(studentID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // get averageMarks
  var averageMarks []models.AverageMark
  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "studentID": studentID,
    "term": termInt,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(averageMarks) != len(student.SubjectList) {
    return c.Status(500).SendString("Nu toate mediile au fost incheiate pe semestrul" + term)
  }

  var value float64 = 0
  for _, averageMark := range averageMarks {
    value += float64(averageMark.Value)
  }
  // still have to find out how it is going to work: whether I have to round, or to floor. 
  // I'll go with floor, for now. I hope I remember
  value = math.Floor(float64(value) / float64(len(student.SubjectList)) * 100) / 100

  // generating id
  var termMarkID = utils.GenID()
  fmt.Println(termMarkID)
  var termMarkGenID models.TermMark
  termMarksCollection.FindOne(context.Background(), bson.M{"termMarkID": termMarkID}).Decode(&termMarkGenID)
  for (termMarkGenID != models.TermMark{Grade: models.Grade{}}) {
    termMarkID = utils.GenID()
    termMarksCollection.FindOne(context.Background(), bson.M{"termMarkID": termMarkID}).Decode(&termMarkGenID)
    fmt.Println(termMarkID)
  } 

  termMark := models.TermMark{
    TermMarkID: termMarkID,
    Value: value,
    StudentID: studentID,
    Grade: student.Grade,
    Term: termInt,
  }

  insertedResult, err := termMarksCollection.InsertOne(context.Background(), termMark)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "termMarkID": termMarkID,
    "value": value,
    "studentID": studentID,
    "grade": student.Grade,
    "term": termInt,
  })
}

// @desc   Get term marks fro homeroom students
// @route  GET /api/teacher/homeroom/term
// @access Private
func getHomeroomTermMarks(c *fiber.Ctx) error {

  // user id and teacher
  userID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getrting the term marks
  termMarksCollection, err := db.GetCollection("termmarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var termMarks []models.TermMark

  cursor, err := termMarksCollection.Find(context.Background(), bson.M{
    "grade.gradeID": teacher.HomeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &termMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(termMarks)
}

// @desc    Get teacher subjects
// @route   GET /api/teacher/subjects
// @access  Private
func getSubjects(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(teacher.SubjectList)
}

// @desc    Get all students that a non-homeroom teacher can interact with
// @route   GET /api/teacher/students
// @access  Private
func getStudents(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var subjectIDList []string
  for _, subject := range teacher.SubjectList {
    subjectIDList = append(subjectIDList, subject.SubjectID)
  }

  // getting the students themselves
  var students []models.Student
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "subjectList.subjectID":  bson.M{"$in": subjectIDList},
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(students)
}

// @desc    Get students that have a subject
// @route   GET /api/teacher/:subjectID
// @access  Private
func getSubjectStudents(c *fiber.Ctx) error {
  var students []models.Student 
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  subjectID := c.Params("subjectID")

  options := options.Find()
  options.SetSort(bson.D{{Key: "lastName", Value: 1}})
  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "subjectList.subjectID": subjectID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(students)
}

// TODO: the rest of it
