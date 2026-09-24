package domain

import "time"

type Grade struct {
	ID           int64
	TestID       int64
	EnrollmentID int64
	Value        float64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

func CalculateWeighedAverage(grades []Grade, tests map[int64]Test) float64 {
	var sum, weightSum float64
	for _, g := range grades {
		t, ok := tests[g.TestID]
		if !ok {
			continue
		}
		sum += g.Value * t.Weight
		weightSum += t.Weight
	}
	if weightSum == 0 {
		return 0
	}
	return sum / weightSum
}

func IsApproved(average float64) bool {
	return average >= 7.0
}
