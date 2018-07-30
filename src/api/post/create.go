package post

import (
	"api/postgres"

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
