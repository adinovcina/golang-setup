package encryption

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestIsValidPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{"valid_password", "$2a$10$bnadVN3crP0gmgt0uajYDuNC6evYIP/36lfnWPSlgZY40qt0wWT2i", "test123", false},
		{"invalid_password", "$2a$10$bnadVN3crP0gmgt0uajYDuNC6evYIP/36lfnWPSlgZY40qt0wWT2i", "test", true},
		{"empty_hashed_password", "", "test", true},
		{"empty_password", "$2a$10$bnadVN3crP0gmgt0uajYDuNC6evYIP/36lfnWPSlgZY40qt0wWT2i", "", true},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := IsValid(tt.hashedPassword, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid(hashedPassword, password string) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestEncrypt(t *testing.T) { //nolint:funlen // need to test multiple different cases
	t.Parallel()

	password := "MyPassword123"

	// Test that a hashed password of the correct length is generated
	t.Run("Encrypt length", func(t *testing.T) {
		t.Parallel()

		hashedPassword, err := Encrypt(password)
		if err != nil {
			t.Error("Error encrypting password: ", err)
		}

		if len(hashedPassword) < 60 {
			t.Error("Hashed password should be at least 60 characters long but was ", len(hashedPassword))
		}
	})

	// Test that the generated hashed password is unique
	t.Run("Encrypt uniqueness", func(t *testing.T) {
		t.Parallel()

		hashedPassword1, err := Encrypt(password)
		if err != nil {
			t.Error("Error encrypting password: ", err)
		}

		hashedPassword2, err := Encrypt(password)
		if err != nil {
			t.Error("Error encrypting password: ", err)
		}

		if hashedPassword1 == hashedPassword2 {
			t.Error("Generated hashed passwords are not unique")
		}
	})

	// Test that the function returns an error if an empty password is provided
	t.Run("Encrypt empty password", func(t *testing.T) {
		t.Parallel()

		_, err := Encrypt("")
		if err != nil {
			t.Error("Expected an error but got nil")
		}
	})

	// Test that the generated hashed password can be used to verify the original password
	t.Run("Encrypt verification", func(t *testing.T) {
		t.Parallel()

		password := "MyPassword123"

		hashedPassword, err := Encrypt(password)
		if err != nil {
			t.Error("Error encrypting password: ", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			t.Error("Error comparing hashed password to original password: ", err)
		}
	})
}
