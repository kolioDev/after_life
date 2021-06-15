package models

import (
	"github.com/gobuffalo/uuid"
)

func (ms *ModelSuite) Test_Will_Create() {
	ms.LoadFixture("user")

	u := &User{}
	ms.NoError(ms.DB.First(u))
	ms.NotEqual(u.ID, uuid.Nil)

	w := &Will{
		Title: "My first will",
	}

	ms.DBDelta(1, "wills", func() {
		verrs, err := w.Create(ms.DB, u)
		ms.NoError(err)
		ms.Falsef(verrs.HasAny(), "Should not have verss but got %v", verrs.Errors)
	})

	//To short title
	w.Title = "1"
	ms.DBDelta(0, "wills", func() {
		ms.Error(w.Create(ms.DB, u))
	})

	//To long title
	w.Title = "1"
	ms.DBDelta(0, "wills", func() {
		ms.Error(w.Create(ms.DB, u))
	})

	ms.Fail("Test creation of instructions")
	ms.Fail("Implement fully")
}
