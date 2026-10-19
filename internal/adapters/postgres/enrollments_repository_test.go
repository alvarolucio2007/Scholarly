package postgres_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

type enrollmentScenario struct {
	User    *domain.User
	User2   *domain.User
	Teacher *domain.Teacher
	Student *domain.Student
	Course  *domain.Course
}

func newEnrollmentScenario(t *testing.T) enrollmentScenario {
	user := &domain.User{Name: "Create enrollment", Email: "createenrollment@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	user2 := &domain.User{Name: "Create student test", Email: "createstudent@email.com", CPF: "00000000002", PasswordHash: []byte("123456789")}
	createUserTest(t, user2)
	teacher := &domain.Teacher{UserID: user.ID, Department: "EnrollmentCreation"}
	createTeacherTest(t, teacher)
	student := &domain.Student{UserID: user2.ID, EnrollmentNumber: "Enrollment"}
	createStudentTest(t, student)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)

	return enrollmentScenario{User: user, User2: user2, Teacher: teacher, Student: student, Course: course}
}

func cleanupEnrollment(t *testing.T, enrollmentID, courseID, studentID, teacherID, user2ID, userID int64) {
	t.Cleanup(func() {

		_, err := testDB.Exec("DELETE FROM enrollments WHERE id = $1", enrollmentID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", courseID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", teacherID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", studentID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", userID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user2ID)
		require.NoError(t, err)
	})
}
func createEnrollmentTest(t *testing.T, enrollment *domain.Enrollment) {
	r := postgres.NewEnrollmentRepo(testDB, testLogger)
	err := r.Create(t.Context(), enrollment)
	require.NoError(t, err)
}

func TestEnrollmentRepo_Create(t *testing.T) {
	s := newEnrollmentScenario(t)
	enrollment := &domain.Enrollment{CourseID: s.Course.ID, StudentID: s.Student.UserID, Status: "active"}
	cleanupEnrollment(t, enrollment.ID, s.Course.ID, s.Student.UserID, s.Teacher.UserID, s.User2.ID, s.User.ID)
	r := postgres.NewEnrollmentRepo(testDB, testLogger)
	tests := []struct {
		name       string // description of this test case
		enrollment *domain.Enrollment
		wantErr    bool
	}{
		{
			name:       "should create enrollment with valid credentials",
			enrollment: &domain.Enrollment{StudentID: s.Student.UserID, CourseID: s.Course.ID, Status: "active"},
			wantErr:    false,
		}, {
			name:       "should not create enrollment with invalid student",
			enrollment: &domain.Enrollment{StudentID: -1, CourseID: s.Course.ID, Status: "active"},
			wantErr:    true,
		}, {
			name:       "should not create enrollment with invalid course",
			enrollment: &domain.Enrollment{StudentID: s.Student.UserID, CourseID: -1, Status: "active"},
			wantErr:    true,
		}, {
			name:       "should not create enrollment with invalid status",
			enrollment: &domain.Enrollment{StudentID: s.Student.UserID, CourseID: -1, Status: "inactive"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}

func TestEnrollmentRepo_GetByID(t *testing.T) {
	s := newEnrollmentScenario(t)
	enrollment := &domain.Enrollment{StudentID: s.Student.UserID, CourseID: s.Course.ID, Status: "active"}
	createEnrollmentTest(t, enrollment)
	cleanupEnrollment(t, enrollment.ID, s.Course.ID, s.Student.UserID, s.Teacher.UserID, s.User2.ID, s.User.ID)
	r := postgres.NewEnrollmentRepo(testDB, testLogger)
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

func TestEnrollmentRepo_List(t *testing.T) {
	s := newEnrollmentScenario(t)
	enrollment := &domain.Enrollment{StudentID: s.Student.UserID, CourseID: s.Course.ID, Status: "active"}
	createEnrollmentTest(t, enrollment)
	cleanupEnrollment(t, enrollment.ID, s.Course.ID, s.Student.UserID, s.Teacher.UserID, s.User2.ID, s.User.ID)
	r := postgres.NewEnrollmentRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		filter  ports.EnrollmentFilter
		want    []*domain.Enrollment
		wantErr bool
	}{
		{
			name:    "should list all enrollments if no paramenters are given",
			filter:  ports.EnrollmentFilter{},
			want:    []*domain.Enrollment{enrollment},
			wantErr: false,
		}, {
			name:    "shouldn't list any enrollment if no enrollments match",
			filter:  ports.EnrollmentFilter{StudentID: &s.Teacher.UserID},
			want:    []*domain.Enrollment{},
			wantErr: false,
		}, {
			name:    "should list only select courses if right paramenters are given",
			filter:  ports.EnrollmentFilter{StudentID: &s.Student.UserID},
			want:    []*domain.Enrollment{enrollment},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.List(t.Context(), tt.filter)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("List() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("List() succeeded unexpectedly")
			}
			require.Equal(t, len(got), len(tt.want))
			for i := range got {
				require.EqualValues(t, &got[i], &tt.want[i])
			}
		})
	}
}

func TestEnrollmentRepo_Update(t *testing.T) {
	s := newEnrollmentScenario(t)
	enrollment := &domain.Enrollment{StudentID: s.Student.UserID, CourseID: s.Course.ID, Status: "active"}
	createEnrollmentTest(t, enrollment)
	cleanupEnrollment(t, enrollment.ID, s.Course.ID, s.Student.UserID, s.Teacher.UserID, s.User2.ID, s.User.ID)

	tests := []struct {
		name       string // description of this test case
		enrollment *domain.Enrollment
		wantErr    bool
	}{
		{
			name:       "shouldn't edit enrollment if all params are empty",
			enrollment: &domain.Enrollment{ID: enrollment.ID},
			wantErr:    false,
		}, {
			name:       "should only try to edit student ID",
			enrollment: &domain.Enrollment{ID: enrollment.ID, StudentID: 999, CourseID: 0},
			wantErr:    true,
		}, {
			name:       "should only try to edit course ID",
			enrollment: &domain.Enrollment{ID: enrollment.ID, StudentID: 0, CourseID: 999},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgres.NewEnrollmentRepo(testDB, testLogger)
			gotErr := r.Update(t.Context(), tt.enrollment)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			updatedEnrollment, err := r.GetByID(t.Context(), enrollment.ID)
			require.NoError(t, err)
			if tt.enrollment.StudentID != 0 {
				require.NotEqual(t, enrollment.StudentID, updatedEnrollment.StudentID)
			}
			if tt.enrollment.CourseID != 0 {
				require.NotEqual(t, enrollment.CourseID, updatedEnrollment.CourseID)
			}
		})
	}
}

func TestEnrollmentRepo_Delete(t *testing.T) {
	s := newEnrollmentScenario(t)
	enrollment := &domain.Enrollment{StudentID: s.Student.UserID, CourseID: s.Course.ID, Status: "active"}
	createEnrollmentTest(t, enrollment)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM enrollments WHERE id=$1", enrollment.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", s.Course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", s.Teacher.UserID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", s.Student.UserID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", s.Student.UserID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", s.User.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", s.User2.ID)
		require.NoError(t, err)
	}) //TODO: implement cleanup func
	r := postgres.NewEnrollmentRepo(testDB, testLogger)
	tests := []struct {
		name         string // description of this test case
		enrollmentID int64
		wantErr      bool
	}{
		{
			name:         "should delete valid enrollment",
			enrollmentID: enrollment.ID,
			wantErr:      false,
		}, {
			name:         "shouldn't delete invalid enrollmentID",
			enrollmentID: -1,
			wantErr:      true,
		}, {
			name:         "shouldn't delete enrollmentID which doesn't exist",
			enrollmentID: 999,
			wantErr:      true,
		}, {
			name:         "should return error if enrollment was already deleted",
			enrollmentID: enrollment.ID,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Delete(t.Context(), tt.enrollmentID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Delete() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Delete() succeeded unexpectedly")
			}
		})
	}
}
