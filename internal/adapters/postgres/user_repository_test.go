package postgres_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/postgres"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

func createUserTest(t *testing.T, user *domain.User) {
	r := postgres.NewUserRepo(testDB, testLogger)
	err := r.Create(t.Context(), user)
	require.NoError(t, err)
}
func cleanupUser(t *testing.T, userID int64) {
	t.Cleanup(func() {
		_, err := testDB.Exec("DELETE FROM users where id=$1", userID)
		require.NoError(t, err)
	})
}

func TestUserRepo_Create(t *testing.T) {
	user := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	r := postgres.NewUserRepo(testDB, testLogger)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		user    *domain.User
		wantErr bool
	}{
		{
			name:    "should accept normal users with all default data",
			user:    user,
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
	cleanupUser(t, user.ID) //why is user.id here populated and not above? TODO: find out
}

func TestUserRepo_GetByID(t *testing.T) {
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	r := postgres.NewUserRepo(testDB, testLogger)
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
	cleanupUser(t, user1.ID)
}

func TestUserRepo_List(t *testing.T) {
	incorrectStr := "INCORRECT"
	r := postgres.NewUserRepo(testDB, testLogger)
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
	cleanupUser(t, user1.ID)
	cleanupUser(t, user2.ID)
}

func TestUserRepo_Update(t *testing.T) {
	r := postgres.NewUserRepo(testDB, testLogger)
	user1 := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user1)
	tests := []struct {
		name    string // description of this test case
		user    *domain.User
		wantErr bool
	}{
		{
			name:    "shouldn't edit user if all params are empty",
			user:    &domain.User{ID: user1.ID, Name: "", CPF: "", Email: "", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only name",
			user:    &domain.User{ID: user1.ID, Name: "changed", CPF: "", Email: "", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only CPF",
			user:    &domain.User{ID: user1.ID, Name: "", CPF: "00000000001", Email: "", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only email",
			user:    &domain.User{ID: user1.ID, Name: "", CPF: "", Email: "changed@gmail.com", PasswordHash: []byte("")},
			wantErr: false,
		}, {
			name:    "should edit only password",
			user:    &domain.User{ID: user1.ID, Name: "", CPF: "", Email: "", PasswordHash: []byte("changed")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Update(t.Context(), tt.user)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Update() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Update() succeeded unexpectedly")
			}
			updatedUser, err := r.GetByID(t.Context(), user1.ID)
			require.NoError(t, err)
			if tt.user.Name != "" {
				require.NotEqual(t, user1.Name, updatedUser.Name)
			}
			if tt.user.CPF != "" {
				require.NotEqual(t, user1.CPF, updatedUser.CPF)
			}
			if tt.user.Email != "" {
				require.NotEqual(t, user1.Email, updatedUser.Email)
			}
			if string(tt.user.PasswordHash) != "" {
				require.NotEqual(t, user1.PasswordHash, updatedUser.PasswordHash)
			}
		})
	}
	cleanupUser(t, user1.ID)
}

func TestUserRepo_Delete(t *testing.T) {
	r := postgres.NewUserRepo(testDB, testLogger)
	user := &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")}
	createUserTest(t, user)
	tests := []struct {
		name    string // description of this test case
		userID  int64
		wantErr bool
	}{
		{
			name:    "should delete with valid ID",
			userID:  user.ID,
			wantErr: false,
		}, {
			name:    "shouldn't delete with invalid ID",
			userID:  user.ID + 1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := r.Delete(t.Context(), tt.userID)
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
}
