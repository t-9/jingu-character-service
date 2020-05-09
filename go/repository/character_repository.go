package repository

import (
	"github.com/t-9/jingu-character-service/go/db"
	"github.com/t-9/jingu-character-service/go/entity"
)

// FindByID finds a character by id.
func FindByID(id uint) (c entity.Character, err error) {
	db, err2 := db.OpenGorm()
	if err2 != nil {
		return c, err2
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	err = db.First(&c, id).Error
	return
}

// RetrieveAll retrieve all character records.
func RetrieveAll() (chars []entity.Character, err error) {
	db, err2 := db.OpenGorm()
	if err2 != nil {
		return chars, err2
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	err = db.Find(&chars).Error
	return
}
