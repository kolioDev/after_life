package models

import (
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/uuid"
)

func (ms *ModelSuite) Test_Trustee_Create() {
	ms.LoadFixture("user")
	t := &Trustee{}
	u := &User{}
	ms.NoError(ms.DB.Where("username=?", "test.single@test.com").First(u))

	//Everything is invalid
	verrs, err := t.Create(ms.DB, User{}.ID)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.ElementsMatch([]string{"user_id", "email", "name", "phone", "relationship"}, verrs.Keys())

	t.UserID = u.ID
	t.Email = "test@gmail.com"
	t.Name = "Geogre Georgiev"
	t.Phone = "+35899999999"
	t.Relationship = "other_friend"
	t.FacebookLink = nulls.NewString("https://www.facebook.com/user.random1234")
	t.TwitterLink = nulls.NewString("https://twitter.com/user.random1234")

	ms.DBDelta(1, "trustees", func() {
		verrs, err := t.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Falsef(verrs.HasAny(), "Should not have validation errors but got %v", verrs.Errors)
		t.ID = uuid.Nil
	})

	t.FacebookLink = nulls.NewString("http://m.facebook.com/user.random321")
	t.TwitterLink = nulls.NewString("http://www.twitter.com/ueer.sssa")
	ms.DBDelta(1, "trustees", func() {
		verrs, err := t.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Falsef(verrs.HasAny(), "Should not have validation errors but got %v", verrs.Errors)
		t.ID = uuid.Nil
	})

	t.FacebookLink = nulls.String{}
	t.TwitterLink = nulls.String{}
	ms.DBDelta(1, "trustees", func() {
		verrs, err := t.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Falsef(verrs.HasAny(), "Should not have validation errors but got %v", verrs.Errors)
		t.ID = uuid.Nil
	})

	t.Email = "invalid.fmail.com"
	t.Name = "shrt"
	t.Phone = "go6o"
	t.Relationship = "qas"
	t.FacebookLink = nulls.NewString("http://www.twitter.com/ueer.sssa")

	ms.DBDelta(0, "trustees", func() {
		verrs, err := t.Create(ms.DB, uuid.FromStringOrNil("bafb92c4-8073-4a35-87a1-7e77bc2efc8a"))
		ms.NoError(err)
		ms.True(verrs.HasAny())
		ms.ElementsMatchf([]string{"user_id", "email", "name", "phone", "relationship", "facebook"}, verrs.Keys(), "arrays mismatcehd got %v", verrs.Keys())
		t.ID = uuid.Nil
	})
}

func (ms *ModelSuite) Test_Trustee_GetForUser() {
	ms.LoadFixture("trustees")

	u := &User{}
	ts := Trustees{}
	ms.NoError(ms.DB.Where("username=?", "trustees.owner1@test.com").First(u))

	ms.NoError(ts.GetForUser(ms.DB, u.ID, "", ""))
	ms.Equal(3, len(ts))

	ms.NoError(ms.DB.Where("username=?", "trustees.owner2@test.com").First(u))
	ms.NoError(ts.GetForUser(ms.DB, u.ID, "", ""))
 	ms.Equal(2, len(ts))
}

func (ms *ModelSuite) Test_Trustee_Update() {
	ms.LoadFixture("trustees")

	u := &User{}
	t := &Trustee{}
	ms.NoError(ms.DB.Where("username=?", "trustees.owner1@test.com").First(u))

	ms.NoError(ms.DB.Where("user_id=?", u.ID).First(t))

	t.UserID = uuid.FromStringOrNil("3eebfd6f-2d2a-4417-8eb8-f4f1a0920e16")
	ms.DBDelta(0, "trustees", func() {
		verrs, err := t.Update(ms.DB)
		ms.NoError(err)
		ms.True(verrs.HasAny())
		ms.ElementsMatchf([]string{"user_id"}, verrs.Keys(), "Expecting user_id verrss but got %v", verrs.Errors)
	})

	t.UserID = u.ID
	t.Name = "inval"
	ms.DBDelta(0, "trustees", func() {
		verrs, err := t.Update(ms.DB)
		ms.NoError(err)
		ms.True(verrs.HasAny())
		ms.ElementsMatchf([]string{"name"}, verrs.Keys(), "Expecting only name verrs but got %v", verrs.Errors)
	})

	t.Name = "Brand New Name"
	ms.DBDelta(0, "trustees", func() {
		verrs, err := t.Update(ms.DB)
		ms.NoError(err)
		ms.False(verrs.HasAny())
	})

	ms.NoError(ms.DB.Find(t, t.ID))
	ms.Equal("Brand New Name", t.Name)
}
