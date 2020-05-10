package entity

import (
	"strconv"
	"testing"
)

func TestSetFieldByNameSuccess(t *testing.T) {
	const expectID = 32
	const expectSurname = "💮♨大好き"
	const expectGivenName = "﷽꧅"
	const unconvertedID = "test"
	c := &Character{}

	c.SetFieldByName("id", strconv.Itoa(expectID))
	if c.ID != expectID {
		t.Errorf("expected ID: %d, actual ID: %d", expectID, c.ID)
	}

	c.SetFieldByName("id", unconvertedID)
	if c.ID != expectID {
		t.Errorf("expected ID: %d, actual ID: %d", expectID, c.ID)
	}

	c.SetFieldByName("surname", expectSurname)
	if c.ID != expectID {
		t.Errorf(
			"expected Surname: %s, actual Surname: %s", expectSurname, c.Surname)
	}

	c.SetFieldByName("givenName", expectGivenName)
	if c.ID != expectID {
		t.Errorf(
			"expected GivenName: %s, actual GivenName: %s",
			expectGivenName, c.GivenName)
	}
}

func TestGetIDSuccess(t *testing.T) {
	const expect = 15
	c := &Character{
		ID: expect,
	}
	if c.GetID() != expect {
		t.Errorf("expected: %d, actual: %d", expect, c.GetID())
	}
}
