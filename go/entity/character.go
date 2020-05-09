package entity

import (
	"errors"
	"time"

	"github.com/t-9/jingu-character-service/go/db"
)

// Character represents columns in the Character table.
type Character struct {
	ID        uint      `gorm:"primary_key"`
	Surname   string    `gorm:"column:surname"`
	GivenName string    `gorm:"column:given_name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// Fillable represents fillable column names.
var Fillable = [...]string{"surname", "givenName"}

// Create a character.
func Create(params map[string]string) (c Character, err error) {
	c = Character{
		Surname:   params["surname"],
		GivenName: params["givenName"],
	}

	db, err2 := db.OpenGorm()
	if err2 != nil {
		return c, err2
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	if !db.NewRecord(&c) {
		return c, errors.New("could not create new record")
	}

	if err2 := db.Create(&c).Error; err2 != nil {
		return c, err2
	}
	return
}

// Destroy a character.
func (c *Character) Destroy() (err error) {
	db, err2 := db.OpenGorm()
	if err2 != nil {
		return err2
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	return db.Delete(c).Error
}

// Update a character.
func (c *Character) Update() (err error) {
	db, err2 := db.OpenGorm()
	if err2 != nil {
		return err2
	}
	defer func() {
		if err2 := db.Close(); err2 != nil {
			err = err2
		}
	}()

	return db.Save(c).Error
}

// SetFieldByName sets a field by name.
func (c *Character) SetFieldByName(name string, value string) {
	switch name {
	case "surname":
		c.Surname = value
	case "givenName":
		c.GivenName = value
	}
}
