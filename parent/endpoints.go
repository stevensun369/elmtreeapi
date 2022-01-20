package parent

import "github.com/gofiber/fiber/v2"

func Endpoints(app *fiber.App) {
	group := app.Group("/api/parent")

  // register parent
  group.Post("/register", parentRegister)
  // login parent
  group.Post("/login", parentLogin)
  // update parent
  group.Get("/update", parentMiddleware, update)

  // get timetable
  group.Get("/timetable/:studentID", parentMiddleware, getPeriods)
  // get periods teacher
  group.Get("/timetable/:studentID/teachers", parentMiddleware, getPeriodsTeachers)
  // get school
  group.Get("/school/:studentID", parentMiddleware, getSchool)

  // add student
  group.Put("/students", parentMiddleware, parentAddStudent)

  // get marks and truancies
  group.Get("/marks/:studentID/:subjectID", parentMiddleware, getMarks)
  group.Get("/truancies/:studentID/:subjectID", parentMiddleware, getTruancies)

  // get averageMarks
  group.Get("/averagemarks", parentMiddleware, getAverageMarks)

  // get termMarks
  group.Get("/termmarks", parentMiddleware, getTermMarks)

  // get finalMarks
  group.Get("/finalmarks/:studentID/:subjectID", parentMiddleware, getFinalMarks)
}
