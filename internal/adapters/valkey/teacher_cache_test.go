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
	teacher := domain.Teacher{UserID: 1}
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    domain.Teacher
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
			item:    domain.Teacher{},
			wantErr: true,
		}, {
			name:    "should throw err if rdb is invalid",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    domain.Teacher{},
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
			cachedTeacher, err := c.Read(t.Context(), teacher.UserID)
			require.NoError(t, err)
			require.Equal(t, teacher, cachedTeacher)
		})
	}
}
func TestTeacherCache_Create_CanceledContext(t *testing.T) {
	c := valkey.NewTeacherCache(testRDB)
	teacher := domain.Teacher{UserID: 1}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, teacher)
	require.ErrorIs(t, err, context.Canceled)
	cleanupTest(t)
}

func TestTeacherCache_Read(t *testing.T) {
	teacher := domain.Teacher{UserID: 1}
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
		want    domain.Teacher
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
			want:    domain.Teacher{},
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			id:      teacher.UserID,
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			want:    domain.Teacher{},
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
			require.Equal(t, got, tt.want)
		})
	}
}

func TestTeacherCache_Delete(t *testing.T) {
	teacher := domain.Teacher{UserID: 1}
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
		wantErr bool
	}{
		{
			name:    "should delete valid teacher",
			rdb:     testRDB,
			id:      teacher.UserID,
			wantErr: false,
		}, {
			name:    "shouldn't function with an invalid redis connection",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			id:      teacher.UserID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewTeacherCache(tt.rdb)
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
