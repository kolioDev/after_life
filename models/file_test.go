package models

import (
	"fmt"
	"github.com/gobuffalo/uuid"
	"strings"
	"time"
)

func (ms *ModelSuite) Test_File_Create() {
	ms.LoadFixture("user")
	u := &User{}
	ms.NoError(ms.DB.Where("username=?", "test.single@test.com").First(u))

	filename := fmt.Sprintf("%d__%s", time.Now().UnixNano(), "test.png")
	f := &File{
		Url:      "http://127.0.0.1:3000/images/pictutes/" + filename,
		Filename: filename,
		Path:     "images",
		Type:     "image",
		FileSize: 300,
	}

	ms.DBDelta(1, "files", func() {
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Falsef(verrs.HasAny(), "Should not have verrs, but got %v", verrs.Errors)
		f.ID = uuid.Nil
	})

	savedF := &File{}
	ms.NoError(ms.DB.Where("filename=?", filename).First(savedF))
	ms.Equal(savedF.OwnerID, u.ID)
	ms.Equal(savedF.Url, f.Url)
	ms.Equal(savedF.Filename, f.Filename)
	ms.Equal(savedF.FileSize, f.FileSize)

	ms.DBDelta(0, "files", func() {
		verrs, err := f.Create(ms.DB, uuid.FromStringOrNil("7aecba58-d7b2-4fcb-aff7-598a6767e42b"))
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
		ms.ElementsMatchf([]string{"owner_id"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())

	})

	ms.DBDelta(0, "files", func() {
		f.Filename = strings.Replace(f.Filename, "png", "jpg", 1)
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
		ms.ElementsMatchf([]string{"url"}, verrs.Keys(), "arrays mismatched got %v %v", verrs.Keys(), verrs.Errors)
	})

	ms.DBDelta(0, "files", func() {
		f.Filename = strings.Replace(f.Filename, "jpg", "exe", 1)
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
		ms.ElementsMatchf([]string{"url", "filename"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())
		f.Filename = strings.Replace(f.Filename, "exe", "png", 1)
	})

	ms.DBDelta(0, "files", func() {
		f.Url = "127.0.0.1:3000/images/pictutes/" + filename
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
		ms.ElementsMatchf([]string{"url"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())
		f.Url = "http://127.0.0.1:3000/images/pictutes/" + filename
	})

	ms.DBDelta(0, "files", func() {
		f.FileSize = 1000000000
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
		ms.ElementsMatchf([]string{"file_size"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())
		f.FileSize = 10000
	})

	ms.DBDelta(0, "files", func() {
		f.Path = ""
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
		ms.ElementsMatchf([]string{"path"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())
		f.Path = "images"
	})

	ms.DBDelta(0, "files", func() {
		f.Type = "car"
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Truef(verrs.HasAny(), "Should have verrs, but got %v", verrs.Errors)
		ms.ElementsMatchf([]string{"type"}, verrs.Keys(), "arrays mismatched got %v", verrs.Keys())
		f.Type = "image"
	})

	ms.DBDelta(1, "files", func() {
		f.Type = "image"
		f.Filename = filename
		f.Url = "http://127.0.0.1:3000/images/" + filename
		f.FileSize = 30000
		verrs, err := f.Create(ms.DB, u.ID)
		ms.NoError(err)
		ms.Falsef(verrs.HasAny(), "Should NOT have verrs, but got %v", verrs.Errors)
	})

}
