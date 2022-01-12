package timetable

import (
	"backend-go/db"
	"backend-go/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func getByGradeID(c *fiber.Ctx) error {
  gradeID := c.Params("gradeID")

  periods, err := db.GetPeriods(bson.M{
    "grade.gradeID": gradeID,
  }, db.PeriodSort)
  utils.CheckError(c, err)

  return c.JSON(periods)
}

func getByTeacherID(c *fiber.Ctx) error {
  teacherID := c.Params("teacherID")

  periods, err := db.GetPeriods(bson.M{
    "teacherID": teacherID,
  }, db.PeriodSort)
  utils.CheckError(c, err)

  return c.JSON(periods)
}