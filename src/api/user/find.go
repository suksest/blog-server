package user

import (
	"api/postgres"
)

type EmailNotExistsError struct{}
type IDNotExistsError struct{}

func (*EmailNotExistsError) Error() string {
	return "email not exists"
}

func (*IDNotExistsError) Error() string {
	return "ID not exists"
}

func FindAll() ([]User, error) {
	db := postgres.OpenDB()
	defer db.Close()

	users := []User{}
	db.Find(&users)

	return users, nil
}

func FindByEmail(email string) (*User, error) {
	var user User

	db := postgres.OpenDB()
	defer db.Close()

	res := db.Find(&user, &User{Email: email})
	if res.RecordNotFound() {
		return nil, &EmailNotExistsError{}
	}
	return &user, nil
}

func FindByID(id uint) (*User, error) {
	var user User

	db := postgres.OpenDB()
	defer db.Close()

	res := db.Find(&user, &User{ID: id})
	if res.RecordNotFound() {
		return nil, &IDNotExistsError{}
	}
	return &user, nil
}
