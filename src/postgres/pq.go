package postgres

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//OpenDB open conncection to postgresql database
func OpenDB() *gorm.DB {
	db, err := gorm.Open("postgres",
		`host=localhost
<<<<<<< Updated upstream
		user=airin password=password
=======
		user=sukma password=openpgpwd
>>>>>>> Stashed changes
		dbname=bookstore
		sslmode=disable`)
	if err != nil {
		panic(err)
	}

	return db
}
