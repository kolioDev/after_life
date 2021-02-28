package models

func (ms *ModelSuite) Test_UserConfirmation_Create() {
	ms.LoadFixture("user")
	u := &User{}
	uc := &UserConfirmation{}

	ms.DBDelta(0, "user_confirmations", func() {
		ms.Error(uc.Create(ms.DB, u))
	})

	ms.NoError(ms.DB.First(u))

	ms.DBDelta(1, "user_confirmations", func() {
		ms.NoError(uc.Create(ms.DB, u))
	})

	ms.Equal(false, uc.KeysSent)
	ms.Equal(false, uc.Confirmed)
	ms.Equal(false, uc.KeysSeenConfirmation)
	ms.Equal(u.ID, uc.UserID)
	ms.NotEqual(0, len(uc.Keys))
}

func (ms *ModelSuite) Test_UserConfirmation_SendKeys() {
	ms.LoadFixture("user")
	u := &User{}
	uc := &UserConfirmation{}

	ms.NoError(ms.DB.First(u))
	ms.NoError(uc.Create(ms.DB, u))
	uc = &UserConfirmation{}

	ms.NoError(ms.DB.Where("user_id=?", u.ID).First(uc))
	ms.False(uc.KeysSent)
	ms.Equal("", uc.Keys)

	uc.Keys = "go6o"
	uc.Confirmed = true
	ms.NoError(uc.SendKeys(ms.DB))
	ms.NotEqual("", uc.Keys)
	ms.NotEqual("go6o", uc.Keys)
	ms.True(uc.KeysSent)

	uc = &UserConfirmation{}
	ms.NoError(ms.DB.Where("user_id=?", u.ID).First(uc))
	ms.False(uc.Confirmed)
	ms.Equal("", uc.Keys)

}
