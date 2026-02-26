package validates

import (
	"regexp"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Minimum 8 characters, maximum 50 characters (configurable).

// At least one uppercase letter (A-Z).

// At least one lowercase letter (a-z).

// At least one number (0-9).

// At least one special character (e.g., !@#$%^&*).

// No spaces.

// Non-nil and non-empty.

func ValidatePassword(val any, fieldName string) (interface{}, error) {
	re, _ := regexp.Compile(`^[^\s]{8,50}$`)

	password := val.(string)

	// Periksa panjang dan tanpa spasi
	if !re.MatchString(password) {
		return nil, gqlerror.Errorf("Field '%s' harus berupa password yang valid (8-50 karakter, minimal satu huruf besar, huruf kecil, angka, dan karakter khusus)", fieldName)
	}

	// Periksa jenis karakter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	isPassword := hasUpper && hasLower && hasDigit && hasSpecial

	if !isPassword {
		return nil, gqlerror.Errorf("Field '%s' harus berupa password yang valid (8-50 karakter, minimal satu huruf besar, huruf kecil, angka, dan karakter khusus)", fieldName)
	}

	return val, nil

}
