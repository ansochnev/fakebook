package backend

import (
	"fmt"
	"strconv"
	"strings"

	"fakebook/internal/db"
)

type AccountID uint64

const InvalidAccountID AccountID = 0

type Account struct {
	ID       AccountID `json:"ID"`
	Username string    `json:"Username"`
	Email    string    `json:"Email"`
}

type CreateAccountREQ struct {
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	Password2 string `json:"Password2"`
	Username  string `json:"Username,omitempty"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

func (b *Backend) CreateAccount(req *CreateAccountREQ) (*Account, error) {
	var err error

	err = b.validate(req)
	if err != nil {
		return nil, err
	}

	passwordSalt := generatePasswordSalt()
	passwordHash := hashPassword(req.Password, passwordSalt)

	tx, err := b.db.Begin()
	if err != nil {
		return nil, beginTransactionError(err)
	}
	defer tx.Rollback()

	accountID, err := insertAccount(tx, req)
	if err != nil {
		return nil, err
	}

	err = insertPassword(tx, accountID, passwordHash, passwordSalt)
	if err != nil {
		return nil, err
	}

	err = insertProfile(tx, accountID, req)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, commitTransactionError(err)
	}

	account := &Account{
		ID:       accountID,
		Email:    req.Email,
		Username: req.Username,
	}

	if account.Username == "" {
		account.Username = strconv.FormatUint(uint64(accountID), 10)
	}

	return account, nil
}

func (b *Backend) validate(r *CreateAccountREQ) error {
	if !isValidEmail(r.Email) {
		return invalidParamError("email", r.Email)
	}

	if r.Username != "" {
		err := validateUsername(r.Username)
		if err != nil {
			return invalidParamError("username", r.Username).Wrap(err)
		}
	}

	if r.FirstName == "" {
		return emptyParamError("firstName")
	}
	if r.LastName == "" {
		return emptyParamError("lastName")
	}

	err := validatePassword(r.Password)
	if err != nil {
		return passwordWeakError().Wrap(err)
	}
	if r.Password != r.Password2 {
		return passwordMismatchError()
	}

	return nil
}

func isValidEmail(email string) bool {
	return strings.ContainsRune(email, '@')
}

func validateUsername(username string) error {
	const charset = "0123456789.-_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for _, char := range username {
		if !strings.ContainsRune(charset, char) {
			return fmt.Errorf("character '%v' is not allowed", string(char))
		}
	}
	return nil
}

func insertAccount(tx db.Tx, req *CreateAccountREQ) (AccountID, error) {
	query := "INSERT INTO accounts (email, username) VALUES (?, ?)"

	result, err := tx.Exec(query, req.Email, req.Username)
	if db.IsDuplicateError(err) {
		columnName, _ := db.MustParseDuplicateError(err)
		if columnName == "email" {
			return InvalidAccountID, emailExistsError()
		}
		return InvalidAccountID, usernameExistsError()
	}
	if err != nil {
		return InvalidAccountID, internalError(err)
	}

	accountID, err := result.LastInsertId()
	if err != nil {
		return InvalidAccountID, internalError(err)
	}
	fmt.Println("accountID:", accountID)

	return AccountID(accountID), nil
}

func insertPassword(tx db.Tx, accountID AccountID, hash binaryHash, salt binarySalt) error {
	query := "INSERT INTO passwords (account_id, hash, salt) VALUES (?, ?, ?)"

	_, err := tx.Exec(query, accountID, hash[:], salt[:])
	if err != nil {
		return internalError(err)
	}

	return nil
}

func insertProfile(tx db.Tx, accountID AccountID, req *CreateAccountREQ) error {
	query := "INSERT INTO profiles (account_id, first_name, last_name) VALUES (?, ?, ?)"

	_, err := tx.Exec(query, accountID, req.FirstName, req.LastName)
	if err != nil {
		return internalError(err)
	}

	return nil
}
