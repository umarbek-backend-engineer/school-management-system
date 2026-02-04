package utils

import (
	"errors"
)

func AuthorizeUsers(userRole string, allowedRoles ...string) (bool, error) {
	for _, value := range allowedRoles {
		if value == userRole {
			return true, nil
		}
	}
	return false, errors.New("User not Authorized")
}
