package post

import (
	"postgres"

	"github.com/jinzhu/gorm"
)

func Create(db *gorm.DB, post *Post) (uint, error) {
	res := db.Create(post)
	if res.Error != nil {
		if postgres.IsUniqueConstraintError(res.Error, UniqueConstraintTitle) {
			return 0, &TitleDuplicateError{}
		}
		return 0, res.Error
	}
	return post.ID, nil
}

// func NewPost(post *Post) (uint, error) {

// 	db := postgres.OpenDB()
// 	defer db.Close()

// 	tx := db.Begin()
// 	if tx.Error != nil {
// 		return 0, tx.Error
// 	}
// 	newPost := &Post{
// 		AuthorID:    post.AuthorID,
// 		Title:       post.Title,
// 		Body:        post.Body,
// 		PublishedAt: time.Now().UTC(),
// 	}
// 	_, err := Create(tx, newPost)
// 	if err != nil {
// 		return 0, err
// 	}
// 	for _, tagName := range post.Tags {
// 		t, err := tag.CreateIfNotExists(tx, tagName)
// 		if err != nil {
// 			tx.Rollback()
// 			return 0, err
// 		}
// 		err = AddTag(tx, newPost, t)
// 		if err != nil {
// 			tx.Rollback()
// 			return 0, err
// 		}
// 	}
// 	res := tx.Commit()
// 	if res.Error != nil {
// 		return 0, res.Error
// 	}
// 	return newPost.ID, nil
// }
