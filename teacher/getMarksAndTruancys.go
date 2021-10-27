package teacher

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"backend-go/db"
	"backend-go/models"

	// std
	"context"
	"fmt"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @desc    Get marks
// @route   GET /api/teacher/mark/:subjectID/:studentID
// @access  Private
func getMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

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

// @desc    Get truancys
// @route   GET /api/teacher/truancy/:subjectID/:studentID
// @access  Private
func getTruancys(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

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

// @desc    Get Average Mark
// @route   POST /api/teacher/average/:subjectID/:studentID
// @access  Private
func getAverageMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var averageMarks []models.AverageMark

  options := options.Find()
  options.SetSort(bson.D{{Key: "term", Value: 1}})
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
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

// @desc    Get Average Marks by subjectID for multiple students
// @route   POST /api/teacher/average/:subjectID
// @access  Private
func getAverageMarksSubject(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")

  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var averageMarks []models.AverageMark

  options := options.Find()
  options.SetSort(bson.D{{Key: "term", Value: 1}})
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
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