package teacher

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  group := app.Group("/api/teacher")

  // login teacher
  group.Post("/login", postLogin)

  // get homeroom info

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
  group.Get("/update", teacherMiddleware, update)

  // get all the students that a non-homeroom teacher can interact with
  group.Get("/students", teacherMiddleware, getStudents)

  // get all the students of a subject
  group.Get("/:subjectID", teacherMiddleware, getSubjectStudents)

  // create mark or truancy
  group.Post("/mark", teacherMiddleware, createMark)
  group.Post("/truancy", teacherMiddleware, createTruancy)

  // get marks or truancies
  group.Get("/mark/:subjectID/:studentID", teacherMiddleware, getMarks)
  group.Get("/truancy/:subjectID/:studentID", teacherMiddleware, getTruancys)

  // motivate truancy
  group.Put("/truancy/motivate", teacherMiddleware, motivateTruancy)

  // create averageMark
  group.Post("/average", teacherMiddleware, createAverageMark)

  // get average marks
  group.Get("/average/:subjectID/:studentID", teacherMiddleware, getAverageMarks)

  // get average marks for multiple students
  group.Get("/average/:subjectID", teacherMiddleware, getAverageMarksSubject)

}