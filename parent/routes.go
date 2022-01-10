package parent

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	group := app.Group("/api/parent")

  // register parent
  group.Post("/register", parentRegister)

  // login parent
  group.Post("/login", parentLogin)

  // add student
  group.Put("/students", parentMiddleware, parentAddStudent)

  // get students
  group.Get("/students", parentMiddleware, getStudents)

  // get marks and truancies
  group.Get("/mark/:studentID/:subjectID", parentMiddleware, getMarks)
  group.Get("/truancy/:studentID/:subjectID", parentMiddleware, getTruancies)

  // get averageMarks
  group.Get("/average", parentMiddleware, getAverageMarks)

  // get termMarks
  group.Get("/term", parentMiddleware, getTermMarks)

}