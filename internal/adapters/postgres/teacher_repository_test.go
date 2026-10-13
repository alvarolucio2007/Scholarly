package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

func createTeacherTest(t *testing.T, teacher *domain.Teacher) {
	r := NewTeacherRepo(testDB, testLogger)
	err := r.Create(t.Context(), teacher)
	require.NoError(t, err)
}

func TestTeacherRepo_Create(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		teacher *domain.Teacher
		wantErr bool
	}{
		{
			name:    "should accept valid userIDs",
			teacher: &domain.Teacher{UserID: user1.ID, Department: "Creation"},
			wantErr: false,
		},
		{
			name:    "should not accept invalid userIDs",
			teacher: &domain.Teacher{UserID: 0, Department: "Invalid"},
			wantErr: true,
		},
		{
			name:    "should not accept duplicated userIDs",
			teacher: &domain.Teacher{UserID: user1.ID, Department: "Duplicated creation"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTeacherRepo(testDB, testLogger)
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
}

func TestTeacherRepo_GetByID(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	teacher1 := &domain.Teacher{UserID: user1.ID, Department: "Testing"}
	createTeacherTest(t, teacher1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
	r := NewTeacherRepo(testDB, testLogger)
	tests := []struct {
		name      string // description of this test case
		teacherID int64
		want      *domain.Teacher
		wantErr   bool
	}{
		{
			name:      "should read valid userIDs",
			teacherID: user1.ID,
			want:      teacher1,
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
			require.EqualValues(t, got.UserID, tt.want.UserID)
			require.EqualValues(t, got.Department, tt.want.Department)
		})
	}
}

func TestTeacherRepo_List(t *testing.T) {
	incorrectStr := "INCORRECT"
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	teacher1 := &domain.Teacher{UserID: user1.ID, Department: "Testing"}
	createTeacherTest(t, teacher1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
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
			r := NewTeacherRepo(testDB, testLogger)
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
				require.EqualValues(t, got[i].UserID, tt.want[i].UserID)
				require.EqualValues(t, got[i].Department, tt.want[i].Department)
			}
		})
	}
}
