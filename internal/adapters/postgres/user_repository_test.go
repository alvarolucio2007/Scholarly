package postgres

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

func TestUserRepo_Create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		user    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserRepo(testDB, testLogger)
			gotErr := r.Create(t.Context(), tt.user)
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
