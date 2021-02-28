package storage

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/envy"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"time"
)

func getLocalStorage() string {
	if envy.Get("GO_ENV", "development") == "test" {
		return envy.Get("LOCAL_TEST_STORAGE_PATH", "storage_test/items")
	}
	return envy.Get("LOCAL_STORAGE_PATH", "storage/items")
}

// Save stores the file on ftp server or local disk
// returns number of bytes, filename
func Save(filename string, buf *bytes.Buffer) (int, error) {

	if envy.Get("SAVE_FILES_ON_FTP_SERVER", "false") == "true" {
		return saveOnFtp(filename, buf)
	}

	filename = fmt.Sprintf("%s/%s", getLocalStorage(), filename)
	return saveOnDisk(filename, buf)
}

func Remove(filename string) error {
	if envy.Get("SAVE_FILES_ON_FTP_SERVER", "false") == "true" {
		return removeFromFtp(filename)
	}

	filename = fmt.Sprintf("%s/%s", getLocalStorage(), filename)
	return removeFromDisk(filename)
}

func Read(filename string) ([]byte, error) {
	if envy.Get("SAVE_FILES_ON_FTP_SERVER", "false") == "true" {
		return readFromFtp(filename)
	}

	filename = fmt.Sprintf("%s/%s", getLocalStorage(), filename)
	return readFromDisk(filename)
}

func saveOnDisk(filename string, buf *bytes.Buffer) (n int, err error) {
	fBytes := buf.Bytes()

	f, err := os.Create(filename)
	if err != nil {
		return
	}

	defer f.Close()

	n, err = f.Write(fBytes)
	return
}

func removeFromDisk(filename string) error {
	return os.Remove(filename)
}

func readFromDisk(filename string) (data []byte, err error) {

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return
}

func saveOnFtp(filename string, buf *bytes.Buffer) (n int, err error) {

	n = buf.Len()

	c, err := ftpConnect()
	if err != nil {
		return
	}
	defer c.Quit()

	err = c.Stor(filename, buf)
	return
}

func removeFromFtp(filename string) error {
	c, err := ftpConnect()
	if err != nil {
		return errors.WithStack(err)
	}
	defer c.Quit()

	return c.Delete(filename)
}

func readFromFtp(filename string) (data []byte, err error) {
	c, err := ftpConnect()
	if err != nil {
		return
	}

	res, err := c.Retr(filename)

	if err != nil {
		return
	}

	_, err = res.Read(data)
	return
}

func ftpConnect() (*ftp.ServerConn, error) {
	url := envy.Get("FTP_URL", "ftp.example.org")
	user := envy.Get("FTP_USER", "anonymous")
	pass := envy.Get("FTP_USER_PASSWORD", "anonymous")
	port := envy.Get("FTP_PORT", "21")

	c, err := ftp.Dial(fmt.Sprintf("%s:%s", url, port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return c, err
	}

	err = c.Login(user, pass)

	return c, err
}
