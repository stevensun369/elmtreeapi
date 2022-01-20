package parent

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Get marks
// @route  GET /api/parent/marks/:studentID/:subjectID
// @access Private
func getMarks(c *fiber.Ctx) error {
  studentID := c.Params("studentID")
  subjectID := c.Params("subjectID")

  marks, err := db.GetMarks(bson.M{
    "studentID": studentID,
    "subject.subjectID": subjectID,
  }, db.DateSort)

  if err != nil {
    return utils.Error(c, err)
  }

  if len(marks) == 0 {
    marks = []models.Mark {}
  }

  return c.JSON(marks)
}

// @desc   Get truancies
// @route  GET /api/parent/truancies/:studentID/:subjectID
// @access Private
func getTruancies(c *fiber.Ctx) error {
  studentID := c.Params("studentID")
  subjectID := c.Params("subjectID")

  truancies, err := db.GetTruancies(bson.M{
    "studentID": studentID,
    "subject.subjectID": subjectID,
  }, db.DateSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(truancies) == 0 {
    truancies = []models.Truancy {}
  }

  return c.JSON(truancies)
}

// @desc   Get average marks of all students
// @route  GET /api/parent/averagemarks
// @access Private
func getAverageMarks(c *fiber.Ctx) error {
  var studentIDList []string
  utils.GetLocals(c.Locals("studentIDList"), &studentIDList)

  averageMarks, err := db.GetAverageMarks(bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, db.TermSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}

// @desc   Get term marks of all students
// @route  GET /api/parent/termmarks
// @access Private
func getTermMarks(c *fiber.Ctx) error {
  var studentIDList []string
  utils.GetLocals(c.Locals("studentIDList"), &studentIDList)

  termMarks, err := db.GetTermMarks(bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, db.TermSort)
  if err != nil {
    return utils.Error(c, err)
  }
  
  if len(termMarks) == 0 {
    termMarks = []models.TermMark {}
  }
  
  return c.JSON(termMarks)
}

// @desc    Get final marks by subjectID
// @route   GET /api/parent/finalmarks/:studentID/:subjectID
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
