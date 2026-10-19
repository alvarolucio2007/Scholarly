package postgres_test

import (
	"testing"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

type testScenario struct {
	User    *domain.User
	Teacher *domain.Teacher
	Course  *domain.Course
}

func newTestScenario(t *testing.T) testScenario {
	user := &domain.User{Name: "Create test", Email: "createtest@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	teacher := &domain.Teacher{UserID: user.ID, Department: "TestCreation"}
	createTeacherTest(t, teacher)
	course := &domain.Course{TeacherID: teacher.UserID, Name: "TestCreation", Code: "TESTCREATION", Semester: "TESTCREATION", MaxStudents: 100}
	createCourseTest(t, course)
	return testScenario{User: user, Teacher: teacher, Course: course}
}
func cleanupTest(t *testing.T, testID, courseID, teacherID, userID int64) {
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM tests WHERE id = $1", testID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM courses WHERE id = $1", courseID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM teachers WHERE user_id=$1", teacherID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", userID)
		require.NoError(t, err)
	})
}

func createTestTest(t *testing.T, test *domain.Test) {
	r := postgres.NewTestRepo(testDB, testLogger)
	err := r.Create(t.Context(), test)
	require.NoError(t, err)
}

func TestTestRepo_Create(t *testing.T) {
	s := newTestScenario(t)
	test := &domain.Test{}
	r := postgres.NewTestRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		test    *domain.Test
		wantErr bool
	}{
		{
			name:    "should create test with valid credentials",
			test:    &domain.Test{CourseID: s.Course.ID, Name: "Successful test creation", Weight: 6, TestDate: time.Now()},
			wantErr: false,
		},
		{
			name:    "shouldn't create test with invalid course",
			test:    &domain.Test{CourseID: 0, Name: "Successful test creation", Weight: 6, TestDate: time.Now()},
			wantErr: true,
		},
		{
			name:    "shouldn't create test with invalid name",
			test:    &domain.Test{CourseID: s.Course.ID, Name: "", Weight: 6, TestDate: time.Now()},
			wantErr: true,
		},
		{
			name:    "shouldn't create test with invalid weight",
			test:    &domain.Test{CourseID: s.Course.ID, Name: "Unsuccessful test creation", Weight: 0, TestDate: time.Now()},
			wantErr: true,
		},
		{
			name:    "shouldn't create test with invalid date",
			test:    &domain.Test{CourseID: s.Course.ID, Name: "Unsuccessful test creation", Weight: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
	cleanupTest(t, test.ID, s.Course.ID, s.Teacher.UserID, s.Teacher.UserID)
}

func TestTestRepo_GetByID(t *testing.T) {
	s := newTestScenario(t)
	test := &domain.Test{CourseID: s.Course.ID, Name: "Successful test creation", Weight: 6}
	createTestTest(t, test)
	r := postgres.NewTestRepo(testDB, testLogger)
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
	cleanupTest(t, test.ID, s.Course.ID, s.Teacher.UserID, s.User.ID)
}

func TestTestRepo_List(t *testing.T) {
	incorrectName := "Not a valid test name"
	s := newTestScenario(t)
	test := &domain.Test{CourseID: s.Course.ID, Name: "Successful test creation", Weight: 6}
	createTestTest(t, test)
	test2 := &domain.Test{CourseID: s.Course.ID, Name: "Test 2", Weight: 4}
	createTestTest(t, test2)
	r := postgres.NewTestRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		filter  ports.TestFilter
		want    []*domain.Test
		wantErr bool
	}{
		{
			name:    "should list all tests if no paramameters are given",
			filter:  ports.TestFilter{},
			want:    []*domain.Test{test2, test},
			wantErr: false,
		}, {
			name:    "shouldn't list any test if no tests match",
			filter:  ports.TestFilter{Name: &incorrectName},
			want:    []*domain.Test{},
			wantErr: false,
		}, {
			name:    "should list only select tests if right parameters are given",
			filter:  ports.TestFilter{Name: &test.Name},
			want:    []*domain.Test{test},
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
				require.EqualValues(t, got[i], tt.want[i])
			}
		})
	}
	cleanupTest(t, test2.ID, s.Course.ID, s.Teacher.UserID, s.User.ID)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM tests WHERE id = $1", test.ID)
		require.NoError(t, err)
	})
}

func TestTestRepo_Update(t *testing.T) {
	s := newTestScenario(t)
	course2 := &domain.Course{TeacherID: s.Teacher.UserID, Name: "Course2", Code: "Test", MaxStudents: 100}
	test := &domain.Test{CourseID: s.Course.ID, Name: "Successful test creation", Weight: 6}
	createTestTest(t, test)
	test2 := &domain.Test{CourseID: s.Course.ID, Name: "Test 2", Weight: 4}
	createTestTest(t, test2)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM tests WHERE id = $1", test.ID)
		require.NoError(t, err)
	})
	r := postgres.NewTestRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		test    *domain.Test
		wantErr bool
	}{
		{
			name:    "shouldn't edit the test if all params are empty",
			test:    &domain.Test{ID: test.ID},
			wantErr: false,
		}, {
			name:    "should only edit courseID",
			test:    &domain.Test{ID: test.ID, CourseID: course2.ID, Name: ""},
			wantErr: false,
		}, {
			name:    "should only edit name",
			test:    &domain.Test{ID: test.ID, CourseID: 0, Name: "New name", Weight: 0},
			wantErr: false,
		}, {
			name:    "should only edit weight",
			test:    &domain.Test{ID: test.ID, CourseID: 0, Name: "", Weight: 9},
			wantErr: false,
		}, {
			name:    "should return error if not found",
			test:    &domain.Test{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Update(t.Context(), tt.test)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Update() succeeded unexpectedly")
			}
			updatedTest, err := r.GetByID(t.Context(), test.ID)
			require.NoError(t, err)
			if tt.test.CourseID != 0 {
				require.NotEqual(t, test.CourseID, updatedTest.CourseID)
			}
			if tt.test.Weight != 0 {
				require.NotEqual(t, test.Weight, updatedTest.Weight)
			}
			if tt.test.Name != "" {
				require.NotEqual(t, test.Name, updatedTest.Name)
			}
		})
	}
	cleanupTest(t, test2.ID, s.Course.ID, s.Teacher.UserID, s.User.ID)
}

func TestTestRepo_Delete(t *testing.T) {
	s := newTestScenario(t)
	test := &domain.Test{CourseID: s.Course.ID, Name: "Successful test creation", Weight: 6}
	createTestTest(t, test)
	tests := []struct {
		name    string // description of this test case
		testID  int64
		wantErr bool
	}{
		{
			name:    "should delete test if ID is valid",
			testID:  test.ID,
			wantErr: false,
		}, {
			name:    "shouldn't delete test if testID is invalid",
			testID:  -1,
			wantErr: true,
		}, {
			name:    "shoudln't delete test if testID doesn't exist",
			testID:  999,
			wantErr: true,
		}, {
			name:    "should return error if test was already deleted",
			testID:  test.ID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgres.NewTestRepo(testDB, testLogger)
			gotErr := r.Delete(t.Context(), tt.testID)
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
	cleanupTest(t, test.ID, s.Course.ID, s.Teacher.UserID, s.User.ID)
}
