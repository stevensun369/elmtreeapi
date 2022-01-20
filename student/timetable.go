package student

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

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
