package student

import (
	"github.com/gofiber/fiber/v2"
)

func Endpoints(app *fiber.App) {
  group := app.Group("/api/student")

  // login student
  group.Post("/login", postLogin)
  // update student data
  group.Get("/update", studentMiddleware, update)

  // get timetable data
  group.Get("/timetable", studentMiddleware, getPeriods)
  // get timetable teachers
  group.Get("/timetable/teachers", studentMiddleware, getPeriodsTeachers)
  // get school data
  group.Get("/school", studentMiddleware, getSchool)

  // get marks or truancies
  group.Get("/marks/:subjectID", studentMiddleware, getMarks)
  group.Get("/truancies/:subjectID", studentMiddleware, getTruancies)

  // get average marks
  group.Get("/averagemarks", studentMiddleware, getAverageMarks)

  // get term marks
  group.Get("/termmarks", studentMiddleware, getTermMarks)

  // get finalMarks
  group.Get("/finalmarks/:subjectID", studentMiddleware, getFinalMarks)
}
