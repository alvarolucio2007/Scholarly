package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func TestTeacherRepo_Create(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
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
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM teachers WHERE user_id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
		require.NoError(t, err)
	})
}
