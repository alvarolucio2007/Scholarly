package domain_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

func TestEnrollmentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		es   domain.EnrollmentStatus
		want bool
	}{
		{
			"statusActive",
			domain.StatusActive,
			true,
		},
		{
			"statusCancelled",
			domain.StatusCancelled,
			true,
		},
		{
			"statusCompleted",
			domain.StatusCompleted,
			true,
		},
		{
			"should return false",
			"false",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.es.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
