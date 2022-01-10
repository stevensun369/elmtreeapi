package student

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Get marks by subject id
// @route  GET /api/student/mark/:subjectID
// @access Private
func getMarks(c *fiber.Ctx) error {
  // student id from locals
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  // subject id from params
  subjectID := c.Params("subjectID")

  marks, err := db.GetMarks(bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, db.DateSort)
  utils.CheckError(c, err)
  
  if len(marks) == 0 {
    marks = []models.Mark {}
  }

  return c.JSON(marks)
}

// @desc   Get truancies by subject id
// @route  GET /api/student/truancy/:subjectID
// @access Private
func getTruancies(c *fiber.Ctx) error {
  // student id from locals
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  // subject id from params
  subjectID := c.Params("subjectID")

  // getting truancies
  truancies, err := db.GetTruancies(bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, db.DateSort)
  utils.CheckError(c, err)

  if len(truancies) == 0 {
    truancies = []models.Truancy {}
  }

  return c.JSON(truancies)
}

// @desc   Get averageMarks
// @route  GET /api/student/average
// @access Private
func getAverageMarks(c *fiber.Ctx) error {
  // student id from locals
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  averageMarks, err := db.GetAverageMarks(bson.M{
    "studentID": studentID,
  }, db.TermSort)
  utils.CheckError(c, err)

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}

// @desc   Get termMarks
// @route  GET /api/student/term
// @access Private
func getTermMarks(c *fiber.Ctx) error {
  // student id from locals
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  termMarks, err := db.GetTermMarks(bson.M{
    "studentID": studentID,
  }, db.EmptySort) 
  utils.CheckError(c, err)

  if len(termMarks) == 0 {
    termMarks = []models.TermMark {}
  }

  return c.JSON(termMarks)
}

// @desc   Get timetable
// @route  GET /api/student/timetable
// @access Private
func getPeriods(c *fiber.Ctx) error {
  var grade models.Grade
  utils.GetLocals(c.Locals("grade"), &grade)

  gradeID := grade.GradeID

  periods, err := db.GetPeriods(bson.M{
    "grade.gradeID": gradeID,
  }, db.PeriodSort)
  utils.CheckError(c, err)

  return c.JSON(periods)
}