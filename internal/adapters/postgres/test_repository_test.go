package postgres

import (
	"testing"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func createTestTest(t *testing.T, test *domain.Test) {
	r := NewTestRepo(testDB, testLogger)
	err := r.Create(t.Context(), test)
	require.NoError(t, err)
}

func TestTestRepo_Create(t *testing.T) {
	var id int64
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM tests WHERE id = $1", id)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		test    *domain.Test
		wantErr bool
	}{
		{
			name:    "should create test with valid credentials",
			test:    &domain.Test{CourseID: course.ID, Name: "Successful test creation", Weight: 6, TestDate: time.Now()},
			wantErr: false,
		},
		{
			name:    "shouldn't create test with invalid course",
			test:    &domain.Test{CourseID: 0, Name: "Successful test creation", Weight: 6, TestDate: time.Now()},
			wantErr: true,
		},
		{
			name:    "shouldn't create test with invalid name",
			test:    &domain.Test{CourseID: course.ID, Name: "", Weight: 6, TestDate: time.Now()},
			wantErr: true,
		},
		{
			name:    "shouldn't create test with invalid weight",
			test:    &domain.Test{CourseID: course.ID, Name: "Unsuccessful test creation", Weight: 0, TestDate: time.Now()},
			wantErr: true,
		},
		{
			name:    "shouldn't create test with invalid date",
			test:    &domain.Test{CourseID: course.ID, Name: "Unsuccessful test creation", Weight: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTestRepo(testDB, testLogger)
			gotErr := r.Create(t.Context(), tt.test)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
			id = tt.test.ID
		})
	}
}

func TestTestRepo_GetByID(t *testing.T) {
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	teacher := &domain.Teacher{UserID: user.ID, Department: "CourseCreation"}
	createTeacherTest(t, teacher)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "GetByID", Code: "GETBYID", Semester: "GETBYIDSEMESTER", MaxStudents: 100}
	createCourseTest(t, course)
	test := &domain.Test{CourseID: course.ID, Name: "Successful test creation", Weight: 6}
	createTestTest(t, test)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM tests WHERE id = $1", test.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", course.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		testID  int64
		want    *domain.Test
		wantErr bool
	}{
		{
			name:    "should return correct test if has correct ID",
			testID:  test.ID,
			want:    test,
			wantErr: false,
		}, {
			name:    "should't return anything if incorrect ID",
			testID:  0,
			want:    nil,
			wantErr: true,
		}, {
			name:    "shouldn't return anything if ID is invalid",
			testID:  -1,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTestRepo(testDB, testLogger)
			got, gotErr := r.GetByID(t.Context(), tt.testID)
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
