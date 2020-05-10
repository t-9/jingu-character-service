package repository

import (
	"github.com/t-9/jingu-character-service/go/db"
	"github.com/t-9/jingu-character-service/go/entity"
)

// FindByID finds a entity record by id.
func FindByID(id uint) (e entity.Character, err error) {
	db, err2 := db.OpenGorm()
	if err2 != nil {
		return e, err2
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	err = db.First(&e, id).Error
	return
}

// RetrieveAll retrieves all entity records.
func RetrieveAll() (e []entity.Character, err error) {
	db, err2 := db.OpenGorm()
	if err2 != nil {
		return e, err2
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	err = db.Find(&e).Error
	return
}
