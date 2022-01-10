package student

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  group := app.Group("/api/student")

  // login student
  group.Post("/login", postLogin)

  // update student data
  group.Get("/update", studentMiddleware, update)

  // get marks or truancies
  group.Get("/mark/:subjectID", studentMiddleware, getMarks)
  group.Get("/truancy/:subjectID", studentMiddleware, getTruancies)

  // get average marks
  group.Get("/average", studentMiddleware, getAverageMarks)

  // get term marks
  group.Get("/term", studentMiddleware, getTermMarks)

  // get timetable data
  group.Get("/timetable", studentMiddleware, getPeriods)
}