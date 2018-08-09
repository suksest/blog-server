package postgres

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

func OpenDB() *gorm.DB {
	db, err := gorm.Open("postgres",
		`host=localhost
		user=sukma password=openpgpwd
		dbname=bookstore
		sslmode=disable`)
	if err != nil {
		panic(err)
	}

	return db
}

func IsUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" && pqErr.Constraint == constraintName
	}
	return false
}
