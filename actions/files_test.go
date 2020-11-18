package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/uuid"
	"github.com/kolioDev/after_life/models"
	"os"
)

const testUploadFileName = "test.jpg"
const testUploadPath = "images/" + testUploadFileName

var storagePath = envy.Get("LOCAL_TEST_STORAGE_PATH", "storage_test/items")

var FILE_TYPES = []string{
	"images", "videos", "audios",
}

func (as *ActionSuite) Test_Files_SaveFile() {
	as.LoadFixture("session")
	clearStorage(as)

	s := &models.Session{}
	as.NoError(as.DB.First(s))
	token := getJWT(as, s)
	f := saveFile(as, token)

	//check if all needed params are returned
	as.Equal("images", f.Path)
	as.Contains(f.Filename, testUploadFileName)

	//check if f.Url works
	res := as.HTML("%s?_token=%s", f.Url, token).Get()
	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_Files_DeleteFile() {
	clearStorage(as)
	as.LoadFixture("sessions")

	//User 1
	s1 := &models.Session{}
	as.NoError(as.DB.Where("unique_token=?", "singe_session_123_1").First(s1))
	token1 := getJWT(as, s1)
	f1 := saveFile(as, token1)

	//User 2
	s2 := &models.Session{}
	as.NoError(as.DB.Where("unique_token=?", "singe_session_123_2").First(s2))
	token2 := getJWT(as, s2)
	f2 := saveFile(as, token2)

	//Users 1 should not be able to delete file owned by user 2
	res := as.JSON("/file/%s?_token=%s", f2.ID, token1).Delete()
	as.Equal(406, res.Code)

	//Non existing file
	res = as.JSON("/file/%s?_token=%s", uuid.FromStringOrNil("a7f6c910-14c8-4314-91d1-74f069f4558d"), token1).Delete()
	as.Equal(404, res.Code)

	//All ok
	res = as.JSON("/file/%s?_token=%s", f1.ID, token1).Delete()
	as.Equal(200, res.Code)
	//check if the file WAS deleted from both DB and storage
	err := as.DB.Find(f1, f1.ID)
	as.Error(err)
	as.Equal(sql.ErrNoRows, err)

	_, err = os.Open(fmt.Sprintf("%s/%s/%s", storagePath, "images/", f1.Filename))
	as.Error(err)
}

func saveFile(as *ActionSuite, token string) *models.File {
	f := &models.File{}

	r, err := os.Open(fmt.Sprintf("%s/%s", storagePath, testUploadFileName))
	as.NoError(err)

	httpF := httptest.File{
		ParamName: "file",
		FileName:  testUploadPath,
		Reader:    r,
	}

	res, err := as.HTML("/file?_token=%s", token).MultiPartPost(f, httpF)
	as.NoError(err)
	as.Equal(200, res.Code)
	as.NoError(json.Unmarshal(res.Body.Bytes(), &f))

	//check if file is saved to disk
	_, err = os.Open(fmt.Sprintf("%s/%s/%s", storagePath, "images/", f.Filename))
	as.NoError(err)

	return f
}

func clearStorage(as *ActionSuite) {
	for _, fileType := range FILE_TYPES {
		as.NoError(os.RemoveAll(fmt.Sprintf("%s/%s", storagePath, fileType)))
		as.NoError(os.Mkdir(fmt.Sprintf("%s/%s", storagePath, fileType), 0755))
	}
}
