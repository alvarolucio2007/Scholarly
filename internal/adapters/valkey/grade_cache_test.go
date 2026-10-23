package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestGradeCache_Create(t *testing.T) {
	grade := domain.Grade{ID: 1}
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    domain.Grade
		wantErr bool
	}{
		{
			name:    "should cache valid teacher",
			rdb:     testRDB,
			item:    grade,
			wantErr: false,
		}, {
			name:    "shouldn't cache empty",
			rdb:     testRDB,
			item:    domain.Grade{},
			wantErr: true,
		}, {
			name:    "should throw err if rdb is invalid",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    grade,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewGradeCache(tt.rdb)
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
			cachedGrade, err := c.Read(t.Context(), grade.ID)
			require.NoError(t, err)
			require.Equal(t, grade, cachedGrade)
		})
	}
}

func TestGradeCache_Create_CanceledContext(t *testing.T) {
	c := valkey.NewGradeCache(testRDB)
	grade := domain.Grade{ID: 10000}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, grade)
	require.ErrorIs(t, err, context.Canceled)
	cleanupTest(t)
}

func TestGradeCache_Read(t *testing.T) {
	grade := domain.Grade{ID: 1}
	studentCache := valkey.NewGradeCache(testRDB)
	err := studentCache.Create(t.Context(), grade)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		id      int64
		want    domain.Grade
		wantErr bool
	}{
		{
			name:    "should read valid grade",
			rdb:     testRDB,
			id:      grade.ID,
			want:    grade,
			wantErr: false,
		}, {
			name:    "should not return anything if ID is invalid",
			rdb:     testRDB,
			id:      grade.ID + 1,
			want:    domain.Grade{},
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			id:      grade.ID,
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			want:    domain.Grade{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewGradeCache(tt.rdb)
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

func TestGradeCache_Delete(t *testing.T) {
	grade := domain.Grade{ID: 1}
	studentCache := valkey.NewGradeCache(testRDB)
	err := studentCache.Create(t.Context(), grade)
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
			name:    "should delete valid grade",
			rdb:     testRDB,
			id:      grade.ID,
			wantErr: false,
		}, {
			name:    "shouldn't function with an invalid redis connection",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			id:      grade.ID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewGradeCache(tt.rdb)
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
