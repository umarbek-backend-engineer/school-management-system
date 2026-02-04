package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("Please enter password")
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("Failed to generate salt")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBased64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	encodedhash := fmt.Sprintf("%s.%s", saltBased64, hashBase64)
	return encodedhash, nil
}

func VerifyPassword(password string, w http.ResponseWriter, req_password string) error {
	parts := strings.Split(password, ".")
	if len(parts) != 2 {
		http.Error(w, "Ooops, something went wrong", http.StatusInternalServerError)
		return fmt.Errorf("Invalid hash format")
	}

	saltBased64 := parts[0]
	hashBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBased64)
	if err != nil {
		http.Error(w, "Failed to decode the salt", http.StatusInternalServerError)
		return errors.New("Failed to decode the salt")
	}

	hash, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		log.Println("Failed to decode the hash:", err)
		http.Error(w, "Failed to decode the hash", http.StatusInternalServerError)
	}

	hashedPassword := argon2.IDKey([]byte(req_password), salt, 1, 64*1024, 4, 32)
	if len(hashedPassword) != len(hash) {
		log.Println("Incorrect Password")
		http.Error(w, "Incorrect Password", http.StatusBadRequest)
		return err
	}

	if subtle.ConstantTimeCompare(hash, hashedPassword) == 1 {
		// Password correct
	} else {
		log.Println("Incorrect Password")
		http.Error(w, "Incorrect Password", http.StatusBadRequest)
		return err
	}
	return nil
}
