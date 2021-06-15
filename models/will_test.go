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
	w.ID = uuid.Nil

	//To short title
	w.Title = "1"
	ms.DBDelta(0, "wills", func() {
		verrs, err := w.Create(ms.DB, u)
		ms.NoError(err)
		ms.True(verrs.HasAny())
	})

	//To long title
	w.Title = "123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+123456789+"
	ms.DBDelta(0, "wills", func() {
		verrs, err := w.Create(ms.DB, u)
		ms.NoError(err)
		ms.True(verrs.HasAny())
	})

	//Will with instructions
	w.Title = "My second will"
	w.Instructions = &Instructions{
		{Text: "Do this first", Index: 0},
		{Text: "Do this second", Index: 1},
		{Text: "Do this third", Index: 2},
	}

	ms.DBDelta(3, "instructions", func() {
		ms.DBDelta(1, "wills", func() {
			verrs, err := w.Create(ms.DB, u)
			ms.NoError(err)
			ms.Falsef(verrs.HasAny(), "Should not have verss but got %v", verrs.Errors)
		})
	})
	w.ID = uuid.Nil

	//Will with instructions - invalid instruction
	w.Title = "My second will"
	w.Instructions = &Instructions{
		{Text: "Do this first", Index: 0},
		{Text: "Do this second", Index: 1},
		{Text: "", Index: 2},
	}

	ms.DBDelta(0, "instructions", func() {
		ms.DBDelta(0, "wills", func() {
			verrs, err := w.Create(ms.DB, u)
			ms.NoError(err)
			ms.True(verrs.HasAny())
		})
	})
	w.ID = uuid.Nil

	//Will with instructions - invalid will title
	w.Title = "1"
	w.Instructions = &Instructions{
		{Text: "Do this first", Index: 0},
		{Text: "Do this second", Index: 1},
		{Text: "Do this third", Index: 2},
	}

	ms.DBDelta(0, "instructions", func() {
		ms.DBDelta(0, "wills", func() {
			verrs, err := w.Create(ms.DB, u)
			ms.NoError(err)
			ms.True(verrs.HasAny())
		})
	})
	w.ID = uuid.Nil

	//TODO:implement files
}
