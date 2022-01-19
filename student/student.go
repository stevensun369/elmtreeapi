package student

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"

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
  if err != nil {
    return utils.Error(c, err)
  } 
  
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
  if err != nil {
    return utils.Error(c, err)
  }

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
  if err != nil {
    return utils.Error(c, err)
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
  // student id from locals
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  termMarks, err := db.GetTermMarks(bson.M{
    "studentID": studentID,
  }, db.EmptySort) 
  if err != nil {
    return utils.Error(c, err)
  }

  if len(termMarks) == 0 {
    termMarks = []models.TermMark {}
  }

  return c.JSON(termMarks)
}

// @desc    Get final marks by subjectID
// @route   GET /api/student/final/:subjectID
// @access  Private
func getFinalMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  
  // getting studentID from locals
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)
  
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

// @desc   Get timetable
// @route  GET /api/student/timetable
// @access Private
func getPeriods(c *fiber.Ctx) error {
  var subjectList []models.ShortSubject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)

  var subjectIDList []string
  for _, subject := range subjectList {
    subjectIDList = append(subjectIDList, subject.SubjectID)
  }

  periods, err := db.GetPeriods(bson.M{
    "subject.subjectID": bson.M{"$in": subjectIDList},
  }, db.PeriodSort)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(periods)
}

// @desc   Get timetable
// @route  GET /api/student/timetable/teachers
// @access Private
func getPeriodsTeachers(c *fiber.Ctx) error {
  var subjectList []models.ShortSubject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)
  

  var subjectIDList []string
  for _, subject := range subjectList {
    subjectIDList = append(subjectIDList, subject.SubjectID)
  }

  teachers, err := db.GetTeachers(bson.M{
    "subjectList.subjectID": bson.M{"$in" : subjectIDList},
  }, db.EmptySort)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(teachers)
}

// @desc   Get timetable
// @route  GET /api/student/timetable
// @access Private
func getSchool(c *fiber.Ctx) error {
  var grade models.Grade
  utils.GetLocals(c.Locals("grade"), &grade)
  
  var school models.School
  err := db.Schools.FindOne(context.Background(), bson.M{
    "schoolID": grade.SchoolID,
  }).Decode(&school)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(school)
}

