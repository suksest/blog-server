package user

import (
	"api/postgres"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type PasswordMismatchError struct{}

func (e *PasswordMismatchError) Error() string {
	return "password didn't match"
}

func Signup(req *User) (uint, error) {

	db := postgres.OpenDB()
	defer db.Close()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	newUser := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := Create(db, newUser)
	if err != nil {
		return 0, err
	}
	return id, err
}

func Create(db *gorm.DB, user *User) (uint, error) {
	err := db.Create(user).Error
	if err != nil {
		if postgres.IsUniqueConstraintError(err, UniqueConstraintUsername) {
			return 0, &UsernameDuplicateError{Username: user.Username}
		}
		if postgres.IsUniqueConstraintError(err, UniqueConstraintEmail) {
			return 0, &EmailDuplicateError{Email: user.Email}
		}
		return 0, err
	}
	return user.ID, nil
}

func Login(req *User) (*User, error) {

	user, err := FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.PasswordHash))
	if err != nil {
		return nil, &PasswordMismatchError{}
	}
	return user, nil
}
