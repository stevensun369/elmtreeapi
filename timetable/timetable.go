package timetable

import (
	"backend-go/db"
	"backend-go/models"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

func getByGradeID(c *fiber.Ctx) error {
  gradeID := c.Params("gradeID")
  periodsCollection, err := db.GetCollection("periods")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var periods []models.Period
  cursor, err := periodsCollection.Find(context.Background(), bson.M{
    "subject.grade.gradeID": gradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &periods); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(periods)
}

func getByTeacherID(c *fiber.Ctx) error {
  teacherID := c.Params("teacherID")
  periodsCollection, err := db.GetCollection("periods")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var periods []models.Period
  cursor, err := periodsCollection.Find(context.Background(), bson.M{
    "teacherID": teacherID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &periods); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(periods)
}