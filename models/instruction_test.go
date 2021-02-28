package models

import (
	"github.com/gobuffalo/uuid"
	"log"
)

func (ms *ModelSuite) Test_Instruction_Create() {
	ms.LoadFixture("user")

	u := &User{}
	ms.NoError(ms.DB.First(u))

	w := &Will{
		Title:    "Do something for my will 1",
		Priority: 1,
		UserID:   u.ID,
	}

	ms.NoError(ms.DB.Create(w))

	//All ok
	i := &Instruction{
		Index: 0,
		Text:  "Do that #1",
	}

	verrs, err := i.Create(ms.DB, *w)
	ms.NoError(err)
	ms.Falsef(verrs.HasAny(), "Should not have verrs, but got %v", verrs.Errors)

	//All ok
	ms.DBDelta(9, "instructions", func() {
		var k uint = 1
		for ; k < 10; k++ {
			i.ID = uuid.Nil
			i.Index = k
			i.Text = "Do that #2"
			log.Println("#K", k)
			verrs, err = i.Create(ms.DB, *w)
			ms.NoError(err)
			ms.Falsef(verrs.HasAny(), "Should not have verrs, but got %v", verrs.Errors)
		}
	})

	//invalid id
	i.ID = uuid.Nil
	i.Index = 2

	verrs, err = i.Create(ms.DB, *w)
	ms.NoError(err)
	ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
	ms.ElementsMatchf([]string{"index"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())

	//invalid id
	i.Index = 20

	verrs, err = i.Create(ms.DB, *w)
	ms.NoError(err)
	ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
	ms.ElementsMatchf([]string{"index"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())

	//invalid will_id
	i.Index = 0
	verrs, err = i.Create(ms.DB, Will{ID: uuid.FromStringOrNil("166c2bbe-6717-404d-b60f-68128cbcc416")})
	ms.NoError(err)
	ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
	ms.ElementsMatchf([]string{"will_id"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())

}

func (ms *ModelSuite) Test_Instructions_Create() {
	ms.LoadFixture("user")

	u := &User{}
	ms.NoError(ms.DB.First(u))

	w := &Will{
		Title:    "Do something for my will 1",
		Priority: 1,
		UserID:   u.ID,
	}

	ms.NoError(ms.DB.Create(w))

	//All ok
	ii := &Instructions{
		//Valid
		{
			Index: 0,
			Text:  "Do that #1",
		},
		//Invalid
		{
			Index: 12,
			Text:  "Do that #1",
		},
	}

	ms.DBDelta(0, "instructions", func() {
		verrs, err := ii.Create(ms.DB, *w)
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should not have verrs, but got %v", verrs.Errors)
	})

}
