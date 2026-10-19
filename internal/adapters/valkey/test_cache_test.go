package valkey_test

import (
	"context"
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/valkey"
	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/go-openapi/testify/require"
	"github.com/redis/go-redis/v9"
)

func TestTestCache_Create(t *testing.T) {
	test := &domain.Test{ID: 1}
	t.Cleanup(func() {
		testRDB.FlushAll(t.Context())
	})
	tests := []struct {
		name    string // description of this test case
		rdb     *redis.Client
		item    *domain.Test
		wantErr bool
	}{
		{
			name:    "should cache valid test",
			rdb:     testRDB,
			item:    test,
			wantErr: false,
		}, {
			name:    "should not cache nil",
			rdb:     testRDB,
			item:    nil,
			wantErr: true,
		}, {
			name:    "should throw err if rdb is down",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			item:    test,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewTestCache(tt.rdb)
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
			cachedTest, err := c.Read(t.Context(), test.ID)
			require.NoError(t, err)
			require.EqualValues(t, test, cachedTest)
		})
	}
}

func TestCacheTest_CancelledCacheContext(t *testing.T) {
	c := valkey.NewTestCache(testRDB)
	test := &domain.Test{ID: 1}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	err := c.Create(ctx, test)
	require.ErrorIs(t, err, context.Canceled)
}

func TestTestCache_Read(t *testing.T) {
	test := &domain.Test{ID: 1}
	testCache := valkey.NewTestCache(testRDB)
	err := testCache.Create(t.Context(), test)
	require.NoError(t, err)
	t.Cleanup(func() {
		testRDB.FlushAll(t.Context())
	})
	tests := []struct {
		name    string // description of this test case
		rdb     *redis.Client
		id      int64
		want    *domain.Test
		wantErr bool
	}{
		{
			name:    "should read valid test",
			rdb:     testRDB,
			id:      test.ID,
			want:    test,
			wantErr: false,
		}, {
			name:    "shouldn't return anything if test doesn't exist",
			rdb:     testRDB,
			id:      0,
			want:    nil,
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			id:      test.ID,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewTestCache(tt.rdb)
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
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestTestCache_Delete(t *testing.T) {
	test := &domain.Test{ID: 1}
	testCache := valkey.NewTestCache(testRDB)
	err := testCache.Create(t.Context(), test)
	require.NoError(t, err)
	t.Cleanup(func() {
		testRDB.FlushAll(t.Context())
	})
	tests := []struct {
		name    string // description of this test case
		rdb     *redis.Client
		id      int64
		wantErr bool
	}{
		{
			name:    "should delete valid test",
			rdb:     testRDB,
			id:      test.ID,
			wantErr: false,
		}, {
			name:    "shouldn't function with invalid redis connection",
			rdb:     redis.NewClient(&redis.Options{Addr: "localhost:9999"}),
			id:      test.ID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valkey.NewTestCache(tt.rdb)
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
