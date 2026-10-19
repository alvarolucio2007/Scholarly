package postgres_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

type teacherScenario struct {
	User *domain.User
}

func newTeacherScenario(t *testing.T) teacherScenario {
	user := &domain.User{Name: "Create teacher", Email: "createTeacher@email.com", CPF: "00000000000", PasswordHash: []byte("123456789")}
	createUserTest(t, user)
	return teacherScenario{User: user}
}

func cleanupTeacher(t *testing.T, teacherID, userID int64) {
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM teachers WHERE user_id = $1", teacherID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users where id=$1", userID)
		require.NoError(t, err)
	})
}
func createTeacherTest(t *testing.T, teacher *domain.Teacher) {
	r := postgres.NewTeacherRepo(testDB, testLogger)
	err := r.Create(t.Context(), teacher)
	require.NoError(t, err)
}

func TestTeacherRepo_Create(t *testing.T) {
	r := postgres.NewTeacherRepo(testDB, testLogger)
	s := newTeacherScenario(t)
	teacher := &domain.Teacher{UserID: s.User.ID, Department: "Creation"}
	tests := []struct {
		name    string // description of this test case
		teacher *domain.Teacher
		wantErr bool
	}{
		{
			name:    "should accept valid userIDs",
			teacher: teacher,
			wantErr: false,
		},
		{
			name:    "should not accept invalid userIDs",
			teacher: &domain.Teacher{UserID: 0, Department: "Invalid"},
			wantErr: true,
		},
		{
			name:    "should not accept duplicated userIDs",
			teacher: &domain.Teacher{UserID: s.User.ID, Department: "Duplicated creation"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Create(t.Context(), tt.teacher)
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
	cleanupTeacher(t, teacher.UserID, s.User.ID)
}

func TestTeacherRepo_GetByID(t *testing.T) {
	s := newTeacherScenario(t)
	teacher := &domain.Teacher{UserID: s.User.ID, Department: "Testing"}
	createTeacherTest(t, teacher)
	r := postgres.NewTeacherRepo(testDB, testLogger)
	tests := []struct {
		name      string // description of this test case
		teacherID int64
		want      *domain.Teacher
		wantErr   bool
	}{
		{
			name:      "should read valid userIDs",
			teacherID: s.User.ID,
			want:      teacher,
			wantErr:   false,
		}, {
			name:      "should not return anything if invalid userIDs",
			teacherID: -1,
			want:      nil,
			wantErr:   true,
		}, {
			name:      "should not return anything if userID doesn't exist",
			teacherID: 100,
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := r.GetByID(t.Context(), tt.teacherID)
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
			require.EqualValues(t, got.UserID, tt.want.UserID)
			require.EqualValues(t, got.Department, tt.want.Department)
		})
	}
	cleanupTeacher(t, teacher.UserID, s.User.ID)
}

func TestTeacherRepo_List(t *testing.T) {
	incorrectStr := "INCORRECT"
	s := newTeacherScenario(t)
	teacher1 := &domain.Teacher{UserID: s.User.ID, Department: "Testing"}
	createTeacherTest(t, teacher1)
	r := postgres.NewTeacherRepo(testDB, testLogger)
	tests := []struct {
		name    string // description of this test case
		filter  ports.TeacherFilter
		want    []*domain.Teacher
		wantErr bool
	}{
		{
			name:    "should list all teachers if no deparmtmentID is given",
			filter:  ports.TeacherFilter{},
			want:    []*domain.Teacher{teacher1},
			wantErr: false,
		}, {
			name:    "shouldn't list any teacher if no teachers match",
			filter:  ports.TeacherFilter{Department: &incorrectStr},
			want:    []*domain.Teacher{},
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
	cleanupTeacher(t, teacher1.UserID, s.User.ID)
}

func TestTeacherRepo_Update(t *testing.T) {
	s := newTeacherScenario(t)
	teacher := &domain.Teacher{UserID: s.User.ID, Department: "Testing"}
	createTeacherTest(t, teacher)
	tests := []struct {
		name    string // description of this test case
		teacher *domain.Teacher
		wantErr bool
	}{
		{
			name:    "shouldn't edit teacher if all params are empty",
			teacher: &domain.Teacher{UserID: s.User.ID, Department: ""},
			wantErr: false,
		}, {
			name:    "should edit only department",
			teacher: &domain.Teacher{UserID: s.User.ID, Department: "Updated"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgres.NewTeacherRepo(testDB, testLogger)
			gotErr := r.Update(t.Context(), tt.teacher)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Update() succeeded unexpectedly")
			}
			updatedTeacher, err := r.GetByID(t.Context(), s.User.ID)
			require.NoError(t, err)
			if tt.teacher.Department != "" {
				require.NotEqual(t, teacher.Department, updatedTeacher.Department)
			}
		})
	}
	cleanupTeacher(t, teacher.UserID, s.User.ID)
}

func TestTeacherRepo_Delete(t *testing.T) {
	s := newTeacherScenario(t)
	teacher1 := &domain.Teacher{UserID: s.User.ID, Department: "Testing"}
	createTeacherTest(t, teacher1)
	r := postgres.NewTeacherRepo(testDB, testLogger)
	tests := []struct {
		name      string // description of this test case
		teacherID int64
		wantErr   bool
	}{
		{
			name:      "should delete with valid userID",
			teacherID: teacher1.UserID,
			wantErr:   false,
		}, {
			name:      "shouldn't delete with invalid ID",
			teacherID: teacher1.UserID + 100,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Delete(t.Context(), tt.teacherID)
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
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users WHERE id=$1", s.User.ID)
		require.NoError(t, err)
	})
}
