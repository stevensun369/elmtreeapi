package teacher

import (
	"github.com/gofiber/fiber/v2"
)

func Endpoints(app *fiber.App) {
  group := app.Group("/api/teacher")

  // login teacher
  group.Post("/login", postLogin)
  // get teacher subjects
  group.Get("/update", teacherMiddleware, update)
  
  // get timetable
  group.Get("/timetable", teacherMiddleware, getPeriods)
  // get school
  group.Get("/school", teacherMiddleware, getSchool)
  
  // get homeroom students
  group.Get("/homeroom", teacherMiddleware, getHomeroomStudents)
  // get homeroom periods
  group.Get("/homeroom/timetable", teacherMiddleware, getHomeroomPeriods)
  // get homeroom periods teachers
  group.Get("/homeroom/timetable/teachers", teacherMiddleware, getHomeroomPeriodsTeachers)
  // get average marks of homeroom students
  group.Get("/homeroom/averagemarks", teacherMiddleware, getHomeroomAverageMarks)
  // get term marks for homeroom students
  group.Get("/homeroom/termmarks", teacherMiddleware, getHomeroomTermMarks)
  // get the subjects of a homeroom student
  group.Get("/homeroom/subjects/:studentID", teacherMiddleware, getHomeroomSubjects)
  // create term mark for student at term
  group.Post("/homeroom/termmarks", teacherMiddleware, createHomeroomTermMark)

  // get all the students that a non-homeroom teacher can interact with
  group.Get("/students", teacherMiddleware, getStudents)
  // get marks or truancies
  group.Get("/marks/:subjectID/:studentID", teacherMiddleware, getMarks)
  group.Get("/truancies/:subjectID/:studentID", teacherMiddleware, getTruancies)
  // get final marks
  group.Get("/finalmarks/:subjectID/:studentID", teacherMiddleware, getFinalMarks)
  // get average marks for multiple students
  group.Get("/averagemarks/:subjectID", teacherMiddleware, getAverageMarks)

  // create mark or truancy
  group.Post("/marks", teacherMiddleware, createMark)
  group.Post("/truancies", teacherMiddleware, createTruancy)
  // motivate truancy
  group.Put("/truancies", teacherMiddleware, motivateTruancy)
  // create finalMark
  group.Post("/finalmarks", teacherMiddleware, createFinalMark)
  // create averageMark
  group.Post("/averagemarks", teacherMiddleware, createAverageMark)

}