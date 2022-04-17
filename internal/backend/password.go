package backend

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type binaryHash [32]byte
type binarySalt [32]byte

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("must be at least 8 characters long")
	}
	return nil
}

func hashPassword(password string, salt binarySalt) binaryHash {
	passwordBytes := []byte(password)
	passwordBytes = append(passwordBytes, salt[:]...)

	const cost = 13
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, cost)
	if err != nil {
		panic(err)
	}

	return sha256.Sum256(hash)
}

func generatePasswordSalt() binarySalt {
	var salt binarySalt
	_, err := rand.Read(salt[:])
	if err != nil {
		panic(err)
	}
	return salt
}
