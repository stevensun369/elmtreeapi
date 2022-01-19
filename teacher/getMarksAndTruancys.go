package teacher

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
)

// @desc    Get marks
// @route   GET /api/teacher/mark/:subjectID/:studentID
// @access  Private
func getMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  marks, err := db.GetMarks(bson.M{
    "subject.subjectID": subjectID,
   "studentID": studentID,
  }, db.DateSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(marks) == 0 {
    marks = []models.Mark {}
  }

  return c.JSON(marks)
}

// @desc    Get truancies
// @route   GET /api/teacher/truancy/:subjectID/:studentID
// @access  Private
func getTruancies(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  truancies, err := db.GetTruancies(bson.M{
    "subject.subjectID": subjectID,
   "studentID": studentID,
  }, db.DateSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(truancies) == 0 {
    truancies = []models.Truancy {}
  }

  return c.JSON(truancies)
}

// @desc    Get Average Mark
// @route   GET /api/teacher/average/:subjectID/:studentID
// @access  Private
func getAverageMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  averageMarks, err := db.GetAverageMarks(bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, db.TermSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}

// @desc    Get Average Marks by subjectID for multiple students
// @route   GET /api/teacher/average/:subjectID
// @access  Private
func getAverageMarksSubject(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")

  averageMarks, err := db.GetAverageMarks(bson.M{
    "subject.subjectID": subjectID,
  }, db.TermSort)
  if err != nil {
    return utils.Error(c, err)
  } 

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}

// @desc    Get final marks by subjectID for one studet
// @route   GET /api/teacher/final/:subjectID/:studentID
// @access  Private
func getFinalMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  finalMarks, err := db.GetFinalMarks(bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, db.EmptySort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(finalMarks) == 0 {
    finalMarks = []models.FinalMark {}
  }

  return c.JSON(finalMarks)
}