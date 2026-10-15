package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func createEnrollmentTest(t *testing.T, enrollment *domain.Enrollment) {
	r := NewEnrollmentRepo(testDB, testLogger)
	err := r.Create(t.Context(), enrollment)
	require.NoError(t, err)
}

func TestEnrollmentRepo_Create(t *testing.T) {
	var id int64
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
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM enrollments WHERE id=$1", id)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", user2.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user2.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name       string // description of this test case
		enrollment *domain.Enrollment
		wantErr    bool
	}{
		{
			name:       "should create enrollment with valid credentials",
			enrollment: &domain.Enrollment{StudentID: student.UserID, CourseID: course.ID, Status: "active"},
			wantErr:    false,
		}, {
			name:       "should not create enrollment with invalid student",
			enrollment: &domain.Enrollment{StudentID: -1, CourseID: course.ID, Status: "active"},
			wantErr:    true,
		}, {
			name:       "should not create enrollment with invalid course",
			enrollment: &domain.Enrollment{StudentID: student.UserID, CourseID: -1, Status: "active"},
			wantErr:    true,
		}, {
			name:       "should not create enrollment with invalid status",
			enrollment: &domain.Enrollment{StudentID: student.UserID, CourseID: -1, Status: "inactive"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewEnrollmentRepo(testDB, testLogger)
			gotErr := r.Create(t.Context(), tt.enrollment)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
			id = tt.enrollment.ID
		})
	}
}

func TestEnrollmentRepo_GetByID(t *testing.T) {
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
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM enrollments WHERE id=$1", enrollment.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", user2.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user2.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name         string // description of this test case
		enrollmentID int64
		want         *domain.Enrollment
		wantErr      bool
	}{
		{
			name:         "should read valid Enrollment IDs",
			enrollmentID: enrollment.ID,
			want:         enrollment,
			wantErr:      false,
		}, {
			name:         "should not return anything if enrollment ID is invalid",
			enrollmentID: -1,
			want:         nil,
			wantErr:      true,
		}, {
			name:         "should not return anything if enrollment ID doesn't exist",
			enrollmentID: enrollment.ID + 1,
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewEnrollmentRepo(testDB, testLogger)
			got, gotErr := r.GetByID(t.Context(), tt.enrollmentID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			require.EqualValues(t, &got, &tt.want)
		})
	}
}
