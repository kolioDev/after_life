package models

func (ms *ModelSuite) Test_Session_Create() {
	ms.LoadFixture("user")
	u := &User{}
	s := &Session{}

	ms.DBDelta(0, "sessions", func() {
		ms.Error(s.Create(ms.DB, u, false))
	})

	ms.NoError(ms.DB.First(u))

	ms.DBDelta(1, "sessions", func() {
		ms.NoError(s.Create(ms.DB, u, false))
	})

	ms.NotEqual("", s.UniqueToken.String)
	ms.NotEqual("", s.RefreshToken)
	ms.Equal(s.UserID, u.ID)
	ms.NotEqual(UUIDNil(), u.ID)
}

func (ms *ModelSuite) Test_Session_ResetUniqueToken() {
	ms.LoadFixture("session")

	s := &Session{}

	ms.Error(s.ResetUniqueToken(ms.DB))

	ms.NoError(ms.DB.Where("unique_token=?", "singe_session_123").First(s))
	ms.True(s.UniqueToken.Valid)
	ms.NotEqual("", s.UniqueToken.String)

	s.ResetUniqueToken(ms.DB)
	ms.False(s.UniqueToken.Valid)
	ms.Equal("", s.UniqueToken.String)
	ms.NoError(ms.DB.Find(s, s.ID))
	ms.False(s.UniqueToken.Valid)
	ms.Equal("", s.UniqueToken.String)

}
