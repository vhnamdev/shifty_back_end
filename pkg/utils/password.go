package utils

import (
	"shifty-backend/pkg/xerror"

	"golang.org/x/crypto/bcrypt"
)

// Password hashing function: retrieves the password from the user when creating an account and hashes that password.
func HashPassword(password string) (string, error) {
	hassPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", xerror.BadRequest("Failed To hashpassword")
	}
	return string(hassPassword), nil
}

// Compare the hashed password stored in the database with the password entered by the user.
func CompareHashAndPassword(password string, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
