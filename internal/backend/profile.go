package backend

import (
	"database/sql"
)

const (
	SexMale   = "male"
	SexFemale = "female"
)

type UserProfile struct {
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	DateOfBirth string `json:"DateOfBirth,omitempty"`
	Sex         string `json:"Sex,omitempty"`
	City        string `json:"City,omitempty"`
	Info        string `json:"Info,omitempty"`
}

func (b *Backend) GetProfileByUsername(username string) (*UserProfile, error) {
	if username == "" {
		return nil, nil
	}

	query :=
		`SELECT
			first_name,
			last_name,
			date_of_birth,
			sex,
			city,
			info
		FROM profiles
		JOIN accounts ON accounts.id = profiles.account_id
		WHERE username = ?`

	row := b.db.QueryRow(query, username)

	var (
		userProfile UserProfile
		dateOfBirth sql.NullString
		sex         sql.NullString
		city        sql.NullString
		info        sql.NullString
	)

	err := row.Scan(
		&userProfile.FirstName,
		&userProfile.LastName,
		&dateOfBirth,
		&sex,
		&city,
		&info)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	userProfile.DateOfBirth = dateOfBirth.String
	userProfile.Sex = sex.String
	userProfile.City = city.String
	userProfile.Info = info.String

	return &userProfile, nil
}
