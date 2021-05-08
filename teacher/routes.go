package teacher

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  group := app.Group("/api/teacher")

  // login teacher
  group.Post("/login", postLogin)

  // get homeroom info
  group.Get("/homeroom/info", teacherMiddleware, getHomeroomInfo)

  // get homeroom students
  group.Get("/homeroom", teacherMiddleware, getHomeroomStudents)

  // get average marks of homeroom students
  group.Get("/homeroom/average", teacherMiddleware, getHomeroomAverageMarks)

  // create term mark for student at term
  group.Post("/homeroom/term", teacherMiddleware, createHomeroomTermMark)

  // get term marks for homeroom students
  group.Get("/homeroom/term", teacherMiddleware, getHomeroomTermMarks)

  // get the subjects of a homeroom student
  group.Get("/homeroom/:studentID", teacherMiddleware, getHomeroomStudentSubjects)

  // get teacher subjects
  group.Get("/subjects", teacherMiddleware, getSubjects)

  // get all the students that a non-homeroom teacher can interact with
  group.Get("/students", teacherMiddleware, getStudents)

  // get all the students of a subject
  group.Get("/:subjectID", teacherMiddleware, getSubjectStudents)
}