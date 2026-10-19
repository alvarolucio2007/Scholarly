package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
)

func TestUserCache_Create(t *testing.T) {
	user := &domain.User{ID: 1, Name: "TestCreateRedis", CPF: "00000000000", Email: "test@createuserredis.com", PasswordHash: []byte("Test")}
	c := valkey.NewUserCache(testRDB)
	tests := []struct {
		name    string // description of this test case
		item    *domain.User
		wantErr bool
	}{
		{
			name:    "should cache valid user",
			item:    user,
			wantErr: false,
		}, {
			name:    "should not cache nil",
			item:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
func TestUserCache_Create_CanceledContext(t *testing.T) {
	c := valkey.NewUserCache(testRDB)
	user := &domain.User{ID: 1, Name: "Test"}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, user)
	require.ErrorIs(t, err, context.Canceled)
}

func TestUserCache_Read(t *testing.T) {
	user := &domain.User{ID: 1, Name: "Test"}
	c := valkey.NewUserCache(testRDB)
	err := c.Create(t.Context(), user)
	require.NoError(t, err)
	tests := []struct {
		name    string // description of this test case
		id      int64
		want    *domain.User
		wantErr bool
	}{
		{
			name:    "should read valid user",
			id:      user.ID,
			want:    user,
			wantErr: false,
		}, {
			name:    "shouldn't return anything if user doesnt exist",
			id:      user.ID + 1,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := c.Read(t.Context(), tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Read() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Read() succeeded unexpectedly")
			}
			require.EqualValues(t, got, tt.want)
		})
	}
}
