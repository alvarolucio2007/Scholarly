package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestUserCache_Create(t *testing.T) {
	user := &domain.User{ID: 10000, Name: "TestCreateRedis", CPF: "00000000000", Email: "test@createuserredis.com", PasswordHash: []byte("Test")}
	cleanupTest(t)
	tests := []struct {
		name    string // description of this test case
		item    *domain.User
		rdb     *redis.Client
		wantErr bool
	}{
		{
			name:    "should cache valid user",
			item:    user,
			rdb:     testRDB,
			wantErr: false,
		}, {
			name:    "should not cache nil",
			item:    nil,
			rdb:     testRDB,
			wantErr: true,
		}, {
			name:    "should throw err if db is down",
			item:    user,
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewUserCache(tt.rdb)
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
	user := &domain.User{ID: 10000, Name: "Test"}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, user)
	require.ErrorIs(t, err, context.Canceled)
	cleanupTest(t)
}

func TestUserCache_Read(t *testing.T) {
	user := &domain.User{ID: 10000, Name: "Test"}
	userCache := valkey.NewUserCache(testRDB)
	err := userCache.Create(t.Context(), user)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name    string // description of this test case
		id      int64
		rdb     *redis.Client
		want    *domain.User
		wantErr bool
	}{
		{
			name:    "should read valid user",
			id:      user.ID,
			rdb:     testRDB,
			want:    user,
			wantErr: false,
		}, {
			name:    "shouldn't return anything if user doesnt exist",
			id:      user.ID + 1,
			rdb:     testRDB,
			want:    nil,
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			id:      user.ID,
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewUserCache(tt.rdb)
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

func TestUserCache_Delete(t *testing.T) {
	user := &domain.User{ID: 10000, Name: "Test"}
	userCache := valkey.NewUserCache(testRDB)
	err := userCache.Create(t.Context(), user)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name    string // description of this test case
		rdb     *redis.Client
		id      int64
		wantErr bool
	}{
		{
			name:    "should delete valid user",
			rdb:     testRDB,
			id:      user.ID,
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			id:      user.ID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewUserCache(tt.rdb)
			gotErr := c.Delete(t.Context(), tt.id)
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
