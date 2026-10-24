package paseto_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/alvarolucio2007/Scholarly/internal/adapters/paseto"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
	"github.com/go-openapi/testify/require"
)

func TestPasetoAuthenticator_CreateToken(t *testing.T) {
	keyB64 := "0wKIJsaVqOKCE9rWYnTZD8ASPn8Ks9YCXpQfwDjQ/b4="
	key, err := base64.StdEncoding.DecodeString(keyB64)
	require.NoError(t, err)
	claims, err := paseto.NewClaims("test@test.com", 5*time.Minute, 1)
	require.NoError(t, err)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		symmetricKey []byte
		// Named input parameters for target function.
		claims  ports.Claims
		wantErr bool
	}{
		{
			name:         "should create valid authenticator with valid key",
			symmetricKey: key,
			claims:       claims,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := paseto.NewPasetoAuthenticator(tt.symmetricKey)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := p.CreateToken(tt.claims)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateToken() succeeded unexpectedly")
			}
			verifiedClaims, err := p.VerifyToken(got)
			require.NoError(t, err)
			require.Equal(t, claims.ID, verifiedClaims.ID)
			require.Equal(t, claims.Email, verifiedClaims.Email)
			require.Equal(t, claims.UserID, verifiedClaims.UserID)
			require.WithinDuration(t, claims.IssuedAt, verifiedClaims.IssuedAt, time.Second)
			require.WithinDuration(t, claims.ExpiresAt, verifiedClaims.ExpiresAt, time.Second)
		})
	}
}
