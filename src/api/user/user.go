package user

import (
	"fmt"
	"postgres"

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

func Delete(id uint) {

	db := postgres.OpenDB()
	defer db.Close()

	theUser := User{
		ID: id,
	}
	res := db.Delete(&theUser)
	fmt.Println(res)

}

func Update(req *User, id uint) (*User, error) {

	db := postgres.OpenDB()
	defer db.Close()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	theUser := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}
	updatedUser := new(User)
	err = db.First(&updatedUser, id).Updates(theUser).Error
	if err != nil {
		if postgres.IsUniqueConstraintError(err, UniqueConstraintUsername) {
			return nil, &UsernameDuplicateError{Username: theUser.Username}
		}
		if postgres.IsUniqueConstraintError(err, UniqueConstraintEmail) {
			return nil, &EmailDuplicateError{Email: theUser.Email}
		}
		return nil, err
	}
	return updatedUser, nil
}
