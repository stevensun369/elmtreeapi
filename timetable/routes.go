package timetable

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/api/timetable")

  // get all periods by gradeID
  g.Get("/student/:gradeID", getByGradeID)

  g.Get("/teacher/:teacherID", getByTeacherID)

}