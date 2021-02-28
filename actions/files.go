package actions

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/kolioDev/after_life/models"
	"github.com/kolioDev/after_life/storage"
	"github.com/pkg/errors"
	"io"
	"path/filepath"
	"time"
)

// FilesSaveFile stores file on disk (or ftp) and makes data entry in the DB
func FilesSaveFile(c buffalo.Context) error {
	var folder2filetype = map[string]string{
		"videos": "video",
		"images": "image",
		"audios": "audio",
	}

	var buf bytes.Buffer
	tx := c.Value("tx").(*pop.Connection)
	u := c.Value("current_user").(*models.User)

	file, err := c.File("file")
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	io.Copy(&buf, file)

	filename := fmt.Sprintf("%s/%d___%s",
		filepath.Dir(file.FileHeader.Filename),
		time.Now().UnixNano(),
		filepath.Base(file.FileHeader.Filename))

	f := models.File{
		Url:      fmt.Sprintf("%s/%s", envy.Get("FILESERVER_URL", "http://127.0.0.1:3000/files"), filename),
		Filename: filepath.Base(filename),
		FileSize: file.FileHeader.Size,
		Path:     filepath.Dir(filename),
		Type:     folder2filetype[filepath.Dir(filename)],
	}

	verrs, err := f.Create(tx, u.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		return c.Render(406, r.JSON(verrs.Errors))
	}

	_, err = storage.Save(filename, &buf)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(f))
}

// FilesDeleteFile removes file on disk (or ftp) and deletes data entry from the DB
func FilesDeleteFile(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	u := c.Value("current_user").(*models.User)

	f := &models.File{}

	err := tx.Find(f, c.Param("id"))
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return c.Render(404, r.JSON("No such file"))
		}
		return errors.WithStack(err)
	}

	if f.OwnerID != u.ID {
		return c.Render(406, r.JSON("Not an owner"))
	}

	err = storage.Remove(filepath.Join(f.Path, f.Filename))
	if err != nil {
		return errors.WithStack(err)
	}

	if err = tx.Destroy(f); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON("Successfully deleted image"))
}

// FilesServeFile used as protected file server
func FilesServeFile(c buffalo.Context) error {

	//TODO::add some security
	//TODO::test

	filename := fmt.Sprintf("%s/%s", c.Param("folder"), c.Param("filename"))

	data, err := storage.Read(filename)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = c.Response().Write(data)
	return err
}
