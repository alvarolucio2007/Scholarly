package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func TestCourseRepo_Create(t *testing.T) {
	var id int64
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM courses WHERE id = $1", id)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		course  *domain.Course
		wantErr bool
	}{
		{
			name:    "should accept normal courses with valid data",
			course:  &domain.Course{TeacherID: teacher.UserID, Name: "valid", Code: "VALID", Semester: "VALID", MaxStudents: 10},
			wantErr: false,
		},
		{
			name:    "shouldn't accept empty courses",
			course:  &domain.Course{},
			wantErr: true,
		},
		{
			name:    "shouldn't acccept courses with invalid teachers",
			course:  &domain.Course{TeacherID: teacher.UserID - 1, Name: "INVALID", Code: "INVALID", Semester: "INVALID", MaxStudents: 10},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewCourseRepo(testDB, testLogger)
			gotErr := r.Create(t.Context(), tt.course)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
			id = tt.course.ID
		})
	}
}
