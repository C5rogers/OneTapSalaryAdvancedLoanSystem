package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("Invalid credentials: email or password is incorrect.")
		check = false
	}
	return check, msg
}
