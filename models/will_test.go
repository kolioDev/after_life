package models

func (ms *ModelSuite) Test_Will_Create() {
	w := &Will{
		Title: "My first will",
	}

	ms.DBDelta(1, "wills", func() {
		ms.NoError(w.Create(ms.DB))
	})

	ms.Fail("Implement fully")
}
