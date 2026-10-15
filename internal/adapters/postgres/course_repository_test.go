package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

func createCourseTest(t *testing.T, course *domain.Course) {
	r := NewCourseRepo(testDB, testLogger)
	err := r.Create(t.Context(), course)
	require.NoError(t, err)
}

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

func TestCourseRepo_GetByID(t *testing.T) {
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name     string // description of this test case
		courseID int64
		want     *domain.Course
		wantErr  bool
	}{
		{
			name:     "should read valid testIDs",
			courseID: course.ID,
			want:     course,
			wantErr:  false,
		}, {
			name:     "should not return anything if courseID is invalid",
			courseID: -1,
			want:     nil,
			wantErr:  true,
		}, {
			name:     "should not return anything if courseID doesn't exist",
			courseID: 1000,
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewCourseRepo(testDB, testLogger)
			got, gotErr := r.GetByID(t.Context(), tt.courseID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			require.EqualValues(t, &course, &got)
		})
	}
}

func TestCourseRepo_List(t *testing.T) {
	incorrectStr := "INCORRECT"
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		filter  ports.CourseFilter
		want    []*domain.Course
		wantErr bool
	}{
		{
			name:    "should list all courses if no parameters are given",
			filter:  ports.CourseFilter{},
			want:    []*domain.Course{course},
			wantErr: false,
		}, {
			name:    "shouldn't list any course if no courses match",
			filter:  ports.CourseFilter{Name: &incorrectStr},
			want:    []*domain.Course{},
			wantErr: false,
		}, {
			name:    "should list only select courses if right paramenters are given",
			filter:  ports.CourseFilter{Name: &course.Name},
			want:    []*domain.Course{course},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewCourseRepo(testDB, testLogger)
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
