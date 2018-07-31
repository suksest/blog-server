package post

import (
	"postgres"
)

type TitleNotExistsError struct{}
type IDNotExistsError struct{}

func (*TitleNotExistsError) Error() string {
	return "Title not exists"
}

func (*IDNotExistsError) Error() string {
	return "ID not exists"
}

func FindAll() ([]Post, error) {
	db := postgres.OpenDB()
	defer db.Close()

	posts := []Post{}
	db.Find(&posts)

	return posts, nil
}

func FindByTitle(title string) (*Post, error) {
	var post Post

	db := postgres.OpenDB()
	defer db.Close()

	res := db.Find(&post, &Post{Title: title})
	if res.RecordNotFound() {
		return nil, &TitleNotExistsError{}
	}
	return &post, nil
}

func FindByID(id uint) (*Post, error) {
	var post Post

	db := postgres.OpenDB()
	defer db.Close()

	res := db.Find(&post, &Post{ID: id})
	if res.RecordNotFound() {
		return nil, &IDNotExistsError{}
	}
	return &post, nil
}
