package domain

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

var (
	ErrNameRequired = errors.New("name is required")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidCPF   = errors.New("invalid CPF")
)

type User struct {
	ID           int64
	Name         string
	CPF          string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

func (u *User) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return ErrNameRequired
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return ErrInvalidEmail
	}
	if !isValidCPF(u.CPF) {
		return ErrInvalidCPF
	}
	return nil
}

var reCPF = regexp.MustCompile(`\D`)

func allSame(s string) bool {
	for i := range len(s) {
		if s[i] != s[0] {
			return false
		}
	}
	return len(s) > 0
}

func isValidCPF(cpf string) bool {
	cpf = reCPF.ReplaceAllString(cpf, "")
	if len(cpf) != 11 || allSame(cpf) {
		return false
	}
	for i := 9; i < 11; i++ {
		soma := 0
		for j := 0; j < i; j++ {
			soma += int(cpf[j]-'0') * ((i + 1) - j)
		}
		dv := (soma * 10 % 11) % 10
		if dv != int(cpf[i]-'0') {
			return false
		}
	}
	return true
}
