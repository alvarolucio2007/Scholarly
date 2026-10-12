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
		{
			name:    "should accept normal users with all default data",
			user:    &domain.User{Name: "Test", CPF: "00000000000", Email: "test@gmail.com", PasswordHash: []byte("123456")},
			wantErr: false,
		},
		{
			name:    "should not accept users without some data",
			user:    &domain.User{Name: "TestWithoutOtherParameters"},
			wantErr: true,
		},
		{
			name:    "should not accept users with same CPF",
			user:    &domain.User{Name: "TestSameCPF", CPF: "00000000000", Email: "test2@gmail.com", PasswordHash: []byte("1234567")},
			wantErr: true,
		},
		{
			name:    "should not accept users with same email",
			user:    &domain.User{Name: "TestSameEmail", CPF: "00000000001", Email: "test@gmail.com", PasswordHash: []byte("123456")},
			wantErr: true,
		},
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
