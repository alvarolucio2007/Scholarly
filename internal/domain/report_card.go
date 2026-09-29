package domain

type ReportCard struct {
	StudentID   int64
	StudentName string
	CourseID    int64
	CourseName  string
	CourseCode  string
	TeacherName string
	Average     float64
	Situation   EnrollmentStatus
}
