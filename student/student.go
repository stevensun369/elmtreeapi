package student

import (
	"backend-go/db"
	"backend-go/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @desc   Get marks by subject id
// @route  GET /api/student/mark/:subjectID
// @access Private
func getMarks(c *fiber.Ctx) error {
  studentIDLocals := fmt.Sprintf("%v", c.Locals("studentID"))
  var studentID string
  json.Unmarshal([]byte(studentIDLocals), &studentID)

  subjectID := c.Params("subjectID")

  marksCollection, err := db.GetCollection("marks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var marks []models.Mark

  options := options.Find()
  options.SetSort(bson.D{{Key: "dateMonth", Value: 1}, {Key: "dateDay", Value: 1}})
  cursor, err := marksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &marks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(marks) == 0 {
    marks = []models.Mark {}
  }

  return c.JSON(marks)
}

// @desc   Get truancys by subject id
// @route  GET /api/student/truancy/:subjectID
// @access Private
func getTruancys(c *fiber.Ctx) error {
  studentIDLocals := fmt.Sprintf("%v", c.Locals("studentID"))
  var studentID string
  json.Unmarshal([]byte(studentIDLocals), &studentID)

  subjectID := c.Params("subjectID")

  truancysCollection, err := db.GetCollection("truancies")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var truancys []models.Truancy

  options := options.Find()
  options.SetSort(bson.D{{Key: "dateMonth", Value: 1}, {Key: "dateDay", Value: 1}})
  cursor, err := truancysCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &truancys); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(truancys) == 0 {
    truancys = []models.Truancy {}
  }

  return c.JSON(truancys)
}

// @desc   Get averageMarks
// @route  GET /api/student/average
// @access Private
func getAverageMarks(c *fiber.Ctx) error {
  studentIDLocals := fmt.Sprintf("%v", c.Locals("studentID"))
  var studentID string
  json.Unmarshal([]byte(studentIDLocals), &studentID)

  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var averageMarks []models.AverageMark

  options := options.Find()
  options.SetSort(bson.D{{Key: "term", Value: 1}})
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "studentID": studentID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}

// @desc   Get termMarks
// @route  GET /api/student/term
// @access Private
func getTermMarks(c *fiber.Ctx) error {
    studentIDLocals := fmt.Sprintf("%v", c.Locals("studentID"))
    var studentID string
    json.Unmarshal([]byte(studentIDLocals), &studentID)

    // getting the term marks
    termMarksCollection, err := db.GetCollection("termmarks")
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }
    var termMarks []models.TermMark
  
    cursor, err := termMarksCollection.Find(context.Background(), bson.M{
      "studentID": studentID,
    })
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }
  
    if err = cursor.All(context.Background(), &termMarks); err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if len(termMarks) == 0 {
      termMarks = []models.TermMark {}
    }
  
    return c.JSON(termMarks)
}