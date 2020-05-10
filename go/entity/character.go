package entity

import (
	"errors"
	"strconv"
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
func (c *Character) Create() (err error) {
	d, err2 := db.OpenGorm()
	if err2 != nil {
		return err2
	}
	defer func() {
		if err2 := d.Close(); err2 != nil {
			err = err2
		}
	}()

	if !d.NewRecord(c) {
		return errors.New("could not create new record")
	}

	return d.Create(c).Error
}

// Destroy a character.
func (c *Character) Destroy() (err error) {
	d, err2 := db.OpenGorm()
	if err2 != nil {
		return err2
	}
	defer func() {
		if err2 := d.Close(); err2 != nil {
			err = err2
		}
	}()

	return d.Delete(c).Error
}

// Update a character.
func (c *Character) Update() (err error) {
	d, err2 := db.OpenGorm()
	if err2 != nil {
		return err2
	}
	defer func() {
		if err2 := d.Close(); err2 != nil {
			err = err2
		}
	}()

	return d.Save(c).Error
}

// SetFieldByName sets a field by name.
func (c *Character) SetFieldByName(name string, value string) {
	switch name {
	case "id":
		u64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return
		}
		c.ID = uint(u64)
	case "surname":
		c.Surname = value
	case "givenName":
		c.GivenName = value
	}
}

// GetID returns a id.
func (c *Character) GetID() uint {
	return c.ID
}
