package validators

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}(\.[a-zA-Z]{2,})?$`)
	return re.MatchString(email)
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func IsValidName(name string) bool {
	return len(name) >= 2 && len(name) <= 50
}

func IsValidPassword(password string) bool {

	return len(password) >= 8
}

func IsValidPhone(phone int) bool {
	return phone != 0
}

func IsEmpty(address string) bool {
	return strings.TrimSpace(address) != ""

}
