package domain

import "testing"

func Test_allSame(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			"all same",
			"11111",
			true,
		},
		{
			"all different",
			"12345",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allSame(tt.s)
			if got != tt.want {
				t.Errorf("allSame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidCPF(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cpf  string
		want bool
	}{
		{
			"valid CPF with punctuation",
			"014.808.178-99",
			true,
		},
		{
			"valid CPF, no punctuation",
			"34693992046",
			true,
		},
		{
			"invalid CPF with punctuation",
			"999.999.999-99",
			false,
		},
		{
			"invalid CPF with no punctuation",
			"99999999999",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidCPF(tt.cpf)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("isValidCPF() = %v, want %v", got, tt.want)
			}
		})
	}
}
