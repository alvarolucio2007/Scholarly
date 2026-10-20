package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestStudentCache_Create(t *testing.T) {
	student := &domain.Student{UserID: 1}
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    *domain.Student
		wantErr bool
	}{
		{
			name:    "should cache valid user",
			rdb:     testRDB,
			item:    student,
			wantErr: false,
		},
		{
			name:    "shouldn't cache nil",
			rdb:     testRDB,
			item:    nil,
			wantErr: true,
		}, {
			name:    "should throw err if rdb is invalid",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    student,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewStudentCache(tt.rdb)
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
			cachedStudent, err := c.Read(t.Context(), student.UserID)
			require.NoError(t, err)
			require.EqualValues(t, student, cachedStudent)
		})
	}
}
func TestStudentCache_Create_CanceledContext(t *testing.T) {
	c := valkey.NewStudentCache(testRDB)
	student := &domain.Student{UserID: 1}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, student)
	require.ErrorIs(t, err, context.Canceled)
	cleanupTest(t)
}

func TestStudentCache_Read(t *testing.T) {
	student := &domain.Student{UserID: 1}
	studentCache := valkey.NewStudentCache(testRDB)
	err := studentCache.Create(t.Context(), student)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		id      int64
		want    *domain.Student
		wantErr bool
	}{
		{
			name:    "should read valid student",
			rdb:     testRDB,
			id:      student.UserID,
			want:    student,
			wantErr: false,
		}, {
			name:    "should not return anything if ID is invalid",
			rdb:     testRDB,
			id:      student.UserID + 1,
			want:    nil,
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			id:      student.UserID,
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewStudentCache(tt.rdb)
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
