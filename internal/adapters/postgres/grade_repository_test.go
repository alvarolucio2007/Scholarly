package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func createGradeTest(t *testing.T, grade *domain.Grade) {
	r := NewGradeRepo(testDB, testLogger)
	err := r.Create(t.Context(), grade)
	require.NoError(t, err)
}

func createGradeConditionsTest(t *testing.T) (*domain.User, *domain.User, *domain.Teacher, *domain.Student, *domain.Course, *domain.Enrollment, *domain.Test) {
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	user2 := &domain.User{Name: "Create student test", Email: "createstudent@email.com", CPF: "00000000002", PasswordHash: []byte("123456789")}
	createUserTest(t, user2)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	student := &domain.Student{UserID: user2.ID, EnrollmentNumber: "Test"}
	createStudentTest(t, student)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)
	enrollment := &domain.Enrollment{StudentID: student.UserID, CourseID: course.ID, Status: "active"}
	createEnrollmentTest(t, enrollment)
	test := &domain.Test{CourseID: course.ID, Name: "Successful test creation", Weight: 6}
	createTestTest(t, test)
	return user, user2, teacher, student, course, enrollment, test
}

func TestGradeRepo_Create(t *testing.T) {
	var id int64
	user, user2, teacher, student, course, enrollment, test := createGradeConditionsTest(t)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM grades WHERE id = $1", id)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM tests WHERE id = $1", test.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM enrollments WHERE id = $1", enrollment.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", teacher.UserID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", student.UserID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user2.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		grade   *domain.Grade
		wantErr bool
	}{
		{
			name:    "should create grade with valid credentials",
			grade:   &domain.Grade{TestID: test.ID, EnrollmentID: enrollment.ID, Value: 10},
			wantErr: false,
		}, {
			name:    "shouldn't create grade with invalid testID",
			grade:   &domain.Grade{TestID: 0, EnrollmentID: enrollment.ID, Value: 10},
			wantErr: true,
		}, {
			name:    "shouldn't create grade with invalid enrollmentID",
			grade:   &domain.Grade{TestID: test.ID, EnrollmentID: 0, Value: 10},
			wantErr: true,
		}, {
			name:    "shouldn't create grade with invalid value",
			grade:   &domain.Grade{TestID: test.ID, EnrollmentID: 0, Value: 1000},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewGradeRepo(testDB, testLogger)
			gotErr := r.Create(t.Context(), tt.grade)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
		})
		id = tt.grade.ID
	}
}
