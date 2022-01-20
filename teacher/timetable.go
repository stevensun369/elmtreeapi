package teacher

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc    getting timetable
// @route   GET /api/teacher/timetable
// @access  Private
func getPeriods(c *fiber.Ctx) error {
  var subjectList []models.Subject
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

// @desc    getting timetable
// @route   GET /api/teacher/school
// @access  Private
func getSchool(c *fiber.Ctx) error {
  var schoolID string
  utils.GetLocals(c.Locals("schoolID"), &schoolID)
 
  var school models.School
  err := db.Schools.FindOne(context.Background(), bson.M{
    "schoolID": schoolID,
  }).Decode(&school)
  if err != nil {
    return utils.Error(c, err)
  }
  
  return c.JSON(school)
}