package post

import (
	"fmt"
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

func Delete(id uint) {

	db := postgres.OpenDB()
	defer db.Close()

	post := Post{
		ID: id,
	}
	res := db.Delete(&post)
	fmt.Println(res)

}

func Update(req *Post, id uint) (*Post, error) {

	db := postgres.OpenDB()
	defer db.Close()

	post := &Post{
		Title: req.Title,
		Body:  req.Body,
	}
	updatedPost := new(Post)
	err := db.First(&updatedPost, id).Updates(post).Error

	return updatedPost, err
}
