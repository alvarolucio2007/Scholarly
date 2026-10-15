package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
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

func TestEnrollmentRepo_List(t *testing.T) {
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
			filter:  ports.EnrollmentFilter{StudentID: &teacher.UserID},
			want:    []*domain.Enrollment{},
			wantErr: false,
		}, {
			name:    "should list only select courses if right paramenters are given",
			filter:  ports.EnrollmentFilter{StudentID: &student.UserID},
			want:    []*domain.Enrollment{enrollment},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewEnrollmentRepo(testDB, testLogger)
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
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	user2 := &domain.User{Name: "Create student test", Email: "createstudent@email.com", CPF: "00000000002", PasswordHash: []byte("123456789")}
	createUserTest(t, user2)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	student := &domain.Student{UserID: user2.ID, EnrollmentNumber: "Test"}
	createStudentTest(t, student)
	student2 := &domain.Student{UserID: user.ID, EnrollmentNumber: "TestStudent"}
	createStudentTest(t, student2)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)
	course2 := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "COPY", Semester: "COPY", MaxStudents: 100}
	createCourseTest(t, course2)
	enrollment := &domain.Enrollment{StudentID: student.UserID, CourseID: course.ID, Status: "active"}
	createEnrollmentTest(t, enrollment)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM enrollments WHERE id=$1", enrollment.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", course2.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", user2.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM students WHERE user_id=$1", user.ID)
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
			name:       "shouldn't edit enrollment if all params are empty",
			enrollment: &domain.Enrollment{ID: enrollment.ID},
			wantErr:    false,
		}, {
			name:       "should only edit student ID",
			enrollment: &domain.Enrollment{ID: enrollment.ID, StudentID: student2.UserID, CourseID: 0},
			wantErr:    false,
		}, {
			name:       "should only edit course ID",
			enrollment: &domain.Enrollment{ID: enrollment.ID, StudentID: 0, CourseID: course2.ID},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewEnrollmentRepo(testDB, testLogger)
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
