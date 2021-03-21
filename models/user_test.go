package models

func (ms *ModelSuite) Test_User_Create() {
	u := &User{}

	//Everything is invalid
	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.ElementsMatch([]string{"provider_id", "provider", "username"}, verrs.Keys())

	//Invalid provider
	u.ProviderID = "123321"
	u.Provider = "ivan"
	u.Name = "test@gmail.com"
	verrs, err = u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.ElementsMatch([]string{"provider"}, verrs.Keys())

	//All valid
	u.Provider = "google"
	ms.DBDelta(1, "users", func() {
		verrs, err = u.Create(ms.DB)
		ms.NoError(err)
		ms.False(verrs != nil && verrs.HasAny())
	})

}
