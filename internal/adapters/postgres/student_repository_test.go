package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func createStudentTest(t *testing.T, student *domain.Student) {
	r := NewStudentRepo(testDB)
	err := r.Create(t.Context(), student)
	require.NoError(t, err)
}

func TestStudentRepo_Create(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM students WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
	})
	tests := []struct {
		name    string // description of this test case
		student *domain.Student
		wantErr bool
	}{
		{
			name:    "should accept valid userIDs",
			student: &domain.Student{UserID: user1.ID, EnrollmentNumber: "TestingCreation"},
			wantErr: false,
		}, {
			name:    "should not accept invalid userIDs",
			student: &domain.Student{UserID: 0, EnrollmentNumber: "Invalid"},
			wantErr: true,
		}, {
			name:    "should not accept duplicated userIDs",
			student: &domain.Student{UserID: user1.ID, EnrollmentNumber: "DuplicatedCreation"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStudentRepo(testDB)
			gotErr := r.Create(t.Context(), tt.student)
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

func TestStudentRepo_GetByID(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		studentID int64
		want      *domain.Student
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewStudentRepo(testDB)
			got, gotErr := r.GetByID(t.Context(), tt.studentID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
