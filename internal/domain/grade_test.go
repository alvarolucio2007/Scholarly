package domain_test

import (
	"testing"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
)

func TestCalculateWeighedAverage(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		grades []domain.Grade
		tests  map[int64]domain.Test
		want   float64
	}{
		{
			"same weights and same test",
			[]domain.Grade{{TestID: 1, Value: 7}, {TestID: 1, Value: 9}},
			map[int64]domain.Test{1: {ID: 1, Weight: 1}, 2: {ID: 2, Weight: 1}},
			8,
		},
		{
			"different weights and same test",
			[]domain.Grade{{TestID: 1, Value: 7}, {TestID: 1, Value: 9}},
			map[int64]domain.Test{1: {ID: 1, Weight: 1}, 2: {ID: 2, Weight: 2}},
			25 / 3,
		},
		{
			"same weights and different tests",
			[]domain.Grade{{TestID: 1, Value: 7}, {TestID: 2, Value: 9}},
			map[int64]domain.Test{1: {ID: 1, Weight: 1}, 2: {ID: 2, Weight: 1}},
			8,
		},
		{
			"different weights and different tests",
			[]domain.Grade{{TestID: 1, Value: 7}, {TestID: 2, Value: 9}},
			map[int64]domain.Test{1: {ID: 1, Weight: 1}, 2: {ID: 2, Weight: 2}},
			25 / 3.0,
		},
		{
			"non-existent tests",
			[]domain.Grade{{TestID: 100, Value: 7}, {TestID: 200, Value: 9}},
			map[int64]domain.Test{1: {ID: 1, Weight: 1}, 2: {ID: 2, Weight: 1}},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.CalculateWeighedAverage(tt.grades, tt.tests)
			if got != tt.want {
				t.Errorf("CalculateWeighedAverage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsApproved(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		average float64
		want    bool
	}{
		{
			"should be approved",
			10.0,
			true,
		}, {
			"should not be approved",
			6.9,
			false,
		}, {
			"exactly 7",
			7,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.IsApproved(tt.average)
			if got != tt.want {
				t.Errorf("IsApproved() = %v, want %v", got, tt.want)
			}
		})
	}
}
