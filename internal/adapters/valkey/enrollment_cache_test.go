package valkey_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestEnrollmentCache_Create(t *testing.T) {
	enrollment := &domain.Enrollment{ID: 1}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		rdb *redis.Client
		// Named input parameters for target function.
		item    *domain.Enrollment
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
			item:    nil,
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
			require.EqualValues(t, enrollment, cachedEnrollment)
		})
	}
}
