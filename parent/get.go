package parent

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @desc   Get a list of the parent's students
// @route  GET /api/parent/students
// @access Private
func getStudents(c *fiber.Ctx) error {
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  parentsCollection, err := db.GetCollection("parents")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var parentID string
  json.Unmarshal([]byte(fmt.Sprintf("%v", c.Locals("parentID"))), &parentID)

  var parent models.Parent
  if err = parentsCollection.FindOne(context.Background(), bson.M{"parentID": parentID}).Decode(&parent); err != nil {
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

  tokenString, err := utils.ParentGenerateToken(parent.ParentID, parent.StudentIDList)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(students) == 0 {
    students = []models.Student {}
  }

  return c.JSON(
    bson.M{
      "students": students,
      "token": tokenString,
    },
  )
}

// @desc   Get marks
// @route  GET /api/parent/marks/:studentID/:subjectID
// @access Private
func getMarks(c *fiber.Ctx) error {
  studentID := c.Params("studentID")
  subjectID := c.Params("subjectID")

  marksCollection, err := db.GetCollection("marks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var marks []models.Mark
  options := options.Find()
  options.SetSort(bson.D{{Key: "dateMonth", Value: 1}, {Key: "dateDay", Value: 1}})
  cursor, err := marksCollection.Find(context.Background(), bson.M{
    "studentID": studentID,
    "subject.subjectID": subjectID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &marks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(marks)
}

// @desc   Get truancys
// @route  GET /api/parent/truancy/:studentID/:subjectID
// @access Private
func getTruancys(c *fiber.Ctx) error {
  studentID := c.Params("studentID")
  subjectID := c.Params("subjectID")

  truancysColection, err := db.GetCollection("truancies")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var truancys []models.Truancy
  options := options.Find()
  options.SetSort(bson.D{{Key: "dateMonth", Value: 1}, {Key: "dateDay", Value: 1}})
  cursor, err := truancysColection.Find(context.Background(), bson.M{
    "studentID": studentID,
    "subject.subjectID": subjectID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &truancys); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(truancys)
}

// @desc   Get average marks of all students
// @route  GET /api/parent/average
// @access Private
func getAverageMarks(c *fiber.Ctx) error {
  var studentIDList []string
  json.Unmarshal([]byte(fmt.Sprintf("%v", c.Locals("studentIDList"))), &studentIDList)

  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var averageMarks []models.AverageMark
  options := options.Find()
  options.SetSort(bson.D{{Key: "term", Value: 1}})
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(averageMarks)
}

// @desc   Get term marks of all students
// @route  GET /api/parent/term
// @access Private
func getTermMarks(c *fiber.Ctx) error {
  var studentIDList []string
  json.Unmarshal([]byte(fmt.Sprintf("%v", c.Locals("studentIDList"))), &studentIDList)

  termMarksCollection, err := db.GetCollection("termmarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var termMarks []models.TermMark
  options := options.Find()
  options.SetSort(bson.D{{Key: "term", Value: 1}})
  cursor, err := termMarksCollection.Find(context.Background(), bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &termMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(termMarks)
}
