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
	cleanupTest(t)
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
	cleanupTest(t)
}

func TestTeacherCache_Read(t *testing.T) {
	teacher := &domain.Teacher{UserID: 1}
	teacherCache := valkey.NewTeacherCache(testRDB)
	err := teacherCache.Create(t.Context(), teacher)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		id      int64
		want    *domain.Teacher
		wantErr bool
	}{
		{
			name:    "should read valid teacher",
			rdb:     testRDB,
			id:      teacher.UserID,
			want:    teacher,
			wantErr: false,
		}, {
			name:    "shouldn't return anything if teacher doesnt exist",
			id:      teacher.UserID + 1,
			rdb:     testRDB,
			want:    nil,
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			id:      teacher.UserID,
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewTeacherCache(tt.rdb)
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
