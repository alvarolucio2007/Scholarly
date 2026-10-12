package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

func createUserTest(t *testing.T, user *domain.User) {
	r := NewUserRepo(testDB, testLogger)
	err := r.Create(t.Context(), user)
	require.NoError(t, err)
}

func TestUserRepo_Create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		user    *domain.User
		wantErr bool
	}{
		{
			name:    "should accept normal users with all default data",
			user:    &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")},
			wantErr: false,
		},
		{
			name:    "should not accept users without some data",
			user:    &domain.User{Name: "TestWithoutOtherParameters"},
			wantErr: true,
		},
		{
			name:    "should not accept users with same CPF",
			user:    &domain.User{Name: "TestSameCPF", CPF: "00000000000", Email: "test2@gmail.com", PasswordHash: []byte("1234567")},
			wantErr: true,
		},
		{
			name:    "should not accept users with same email",
			user:    &domain.User{Name: "TestSameEmail", CPF: "00000000001", Email: "test@gmail.com", PasswordHash: []byte("123456")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserRepo(testDB, testLogger)
			gotErr := r.Create(t.Context(), tt.user)
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
		_, err := testDB.Exec("DELETE FROM users WHERE id=$1", 1)
		require.NoError(t, err)
	})
}

func TestUserRepo_GetByID(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	user2 := &domain.User{Name: "Test", CPF: "00000000001", Email: "test2@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user2)
	tests := []struct {
		name    string // description of this test case
		userID  int64
		want    *domain.User
		wantErr bool
	}{
		{
			name:    "should accept normal data and return existing user",
			userID:  user1.ID,
			want:    user1,
			wantErr: false,
		},
		{
			name:    "should return error when user doesn't exist",
			userID:  0,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserRepo(testDB, testLogger)
			got, gotErr := r.GetByID(t.Context(), tt.userID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetByID() succeeded unexpectedly")
			}
			require.EqualValues(t, got.ID, tt.want.ID)
			require.EqualValues(t, got.Email, tt.want.Email)
			require.EqualValues(t, got.CPF, tt.want.CPF)
			require.EqualValues(t, got.CreatedAt, tt.want.CreatedAt)
			require.EqualValues(t, got.UpdatedAt, tt.want.UpdatedAt)
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user2.ID)
		require.NoError(t, err)
	})
}

func TestUserRepo_List(t *testing.T) {
	incorrectStr := "INCORRECT"
	user1 := &domain.User{Name: "TestOne", CPF: "00000000000", Email: "testOne@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	user2 := &domain.User{Name: "TestTwo", CPF: "00000000001", Email: "testTwo@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user2)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		filter  ports.UserFilter
		want    []*domain.User
		wantErr bool
	}{
		{
			name:    "with no filter, should return all users",
			filter:  ports.UserFilter{},
			want:    []*domain.User{user2, user1}, // reversed because of the ORDER BY created_at DESC in line 78 of user_repository
			wantErr: false,
		}, {
			name:    "should filter by name",
			filter:  ports.UserFilter{Name: &user1.Name},
			want:    []*domain.User{user1},
			wantErr: false,
		}, {
			name:    "should filter by CPF",
			filter:  ports.UserFilter{CPF: &user1.CPF},
			want:    []*domain.User{user1},
			wantErr: false,
		}, {
			name:    "should filter by email",
			filter:  ports.UserFilter{Email: &user1.Email},
			want:    []*domain.User{user1},
			wantErr: false,
		}, {
			name:    "should not return anyone if filters are invalid",
			filter:  ports.UserFilter{Name: &incorrectStr},
			want:    []*domain.User{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserRepo(testDB, testLogger)
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
				require.EqualValues(t, got[i].ID, tt.want[i].ID)
				require.EqualValues(t, got[i].Email, tt.want[i].Email)
				require.EqualValues(t, got[i].CPF, tt.want[i].CPF)
				require.EqualValues(t, got[i].CreatedAt, tt.want[i].CreatedAt)
				require.EqualValues(t, got[i].UpdatedAt, tt.want[i].UpdatedAt)
			}
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users WHERE id=$1", user1.ID)
		require.NoError(t, err)
		_, err = testDB.Exec("DELETE FROM users WHERE id=$1", user2.ID)
		require.NoError(t, err)
	})
}
