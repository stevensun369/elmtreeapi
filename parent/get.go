package parent

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Get a list of the parent's students
// @route  GET /api/parent/students
// @access Private
func getStudents(c *fiber.Ctx) error {
  var parentID string
  utils.GetLocals(c.Locals("parentID"), &parentID)

  var studentIDList []string
  utils.GetLocals(c.Locals("studentIDList"), &studentIDList)

  students, err := db.GetStudents(bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, db.GradeSort)

  if err != nil {
    return utils.Error(c, err)
  }

  tokenString, err := utils.ParentGenerateToken(parentID, studentIDList)
  if err != nil {
    return utils.Error(c, err)
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
// @route  GET /api/parent/truancy/:studentID/:subjectID
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
// @route  GET /api/parent/average
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
// @route  GET /api/parent/term
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
// @route   GET /api/parent/final/:studentID/:subjectID
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

func getPeriods(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    utils.Error(c, err)
  }

  subjectList := student.SubjectList
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
// @route  GET /api/parent/timetable/:studentID/teachers
// @access Private
func getPeriodsTeachers(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    return utils.Error(c, err)
  }

  subjectList := student.SubjectList
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

// @desc   Get school
// @route  GET /api/parent/school/:studentID/
// @access Private
func getSchool(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    return utils.Error(c, err)
  }
  
  var school models.School
  err = db.Schools.FindOne(context.Background(), bson.M{
    "schoolID": student.Grade.SchoolID,
  }).Decode(&school)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(school)
}