package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestEnrollmentCache_Create(t *testing.T) {
	enrollment := domain.Enrollment{ID: 1}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    domain.Enrollment
		wantErr bool
	}{
		{
			name:    "should cache valid enrollment",
			rdb:     testRDB,
			item:    enrollment,
			wantErr: false,
		},
		{
			name:    "shouldn't cache nil",
			rdb:     testRDB,
			item:    domain.Enrollment{},
			wantErr: true,
		}, {
			name:    "should throw err if rdb is invalid",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    enrollment,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewEnrollmentCache(tt.rdb)
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
			cachedEnrollment, err := c.Read(t.Context(), enrollment.ID)
			require.NoError(t, err)
			require.Equal(t, enrollment, cachedEnrollment)
		})
	}
}
func TestEnrollmentCache_Create_CanceledContext(t *testing.T) {
	c := valkey.NewEnrollmentCache(testRDB)
	enrollment := domain.Enrollment{ID: 1}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, enrollment)
	require.ErrorIs(t, err, context.Canceled)
	cleanupTest(t)
}

func TestEnrollmentCache_Read(t *testing.T) {
	enrollment := domain.Enrollment{ID: 1}
	studentCache := valkey.NewEnrollmentCache(testRDB)
	err := studentCache.Create(t.Context(), enrollment)
	require.NoError(t, err)
	cleanupTest(t)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		id      int64
		want    domain.Enrollment
		wantErr bool
	}{
		{
			name:    "should read valid grade",
			rdb:     testRDB,
			id:      enrollment.ID,
			want:    enrollment,
			wantErr: false,
		}, {
			name:    "should not return anything if ID is invalid",
			rdb:     testRDB,
			id:      enrollment.ID + 1,
			want:    domain.Enrollment{},
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			id:      enrollment.ID,
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			want:    domain.Enrollment{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewEnrollmentCache(tt.rdb)
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

func TestEnrollmentCache_Delete(t *testing.T) {
	enrollment := domain.Enrollment{ID: 1}
	studentCache := valkey.NewEnrollmentCache(testRDB)
	err := studentCache.Create(t.Context(), enrollment)
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
			id:      enrollment.ID,
			wantErr: false,
		}, {
			name:    "shouldn't function with an invalid redis connection",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			id:      enrollment.ID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewEnrollmentCache(tt.rdb)
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
