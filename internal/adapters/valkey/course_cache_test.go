package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestCourseCache_Create(t *testing.T) {
	course := &domain.Course{}
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    *domain.Course
		wantErr bool
	}{
		{
			name:    "should cache valid teacher",
			rdb:     testRDB,
			item:    course,
			wantErr: false,
		}, {
			name:    "shouldn't cache nil",
			rdb:     testRDB,
			item:    nil,
			wantErr: true,
		}, {
			name:    "should throw err if rdb is invalid",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    course,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewCourseCache(tt.rdb)
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
			cachedCourse, err := c.Read(t.Context(), course.ID)
			require.NoError(t, err)
			require.EqualValues(t, course, cachedCourse)
		})
	}
}

func TestCourseCache_Create_CanceledContext(t *testing.T) {
	c := valkey.NewCourseCache(testRDB)
	course := &domain.Course{ID: 10000}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, course)
	require.ErrorIs(t, err, context.Canceled)
	cleanupTest(t)
}
