package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestTeacherCache_Create(t *testing.T) {
	teacher := &domain.Teacher{UserID: 1}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    *domain.Teacher
		wantErr bool
	}{
		{
			name:    "should cache valid teacher",
			rdb:     testRDB,
			item:    teacher,
			wantErr: false,
		}, {
			name:    "shouldn't cache nil",
			rdb:     testRDB,
			item:    nil,
			wantErr: true,
		}, {
			name:    "should throw err if rdb is invalid",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    teacher,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewTeacherCache(tt.rdb)
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
		})
	}
}
func TestTeacherCache_Create_CanceledContext(t *testing.T) {
	c := valkey.NewTeacherCache(testRDB)
	teacher := &domain.Teacher{UserID: 1}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, teacher)
	require.ErrorIs(t, err, context.Canceled)
}
