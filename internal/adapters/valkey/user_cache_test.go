package valkey_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func TestUserCache_Create(t *testing.T) {
	user := &domain.User{ID: 1, Name: "TestCreateRedis", CPF: "00000000000", Email: "test@createuserredis.com", PasswordHash: []byte("Test")}
	tests := []struct {
		name    string // description of this test case
		item    *domain.User
		wantErr bool
	}{
		{
			name:    "should cache valid user",
			item:    user,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewUserCache(testRDB)
			gotErr := c.Create(t.Context(), tt.item)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Create() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Create() succeeded unexpectedly")
			}
			cachedUser, err := c.Read(t.Context(), user.ID)
			require.NoError(t, err)
			require.EqualValues(t, user, cachedUser)
		})
	}
}
