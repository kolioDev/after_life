package models

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/validate/validators"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

const MAX_FILESIZE = 1000000000 //in bytes = 1 GB

var ACCEPTABLE_FILETYPES = map[string][]string{
	//Image
	"image": {
		"png",
		"jpg",
		"jp2",
		"jpeg",
		"gif",
		"bmp",
		"wbmp",
		"pgm",
	},

	//video
	"video": {
		"webm",
		"mkv",
		"flv",
		"vob",
		"ogv",
		"ogg",
		"gifv",
		"avi",
		"mng",
		"mov",
		"wmv",
		"mp4",
		"m4p",
		"m4v",
		"mpg",
		"mp2",
		"mpeg",
		"mpe",
		"mpv",
		"m2v",
		"m4v",
	},
	//audio
	"audio": {
		"3gp",
		"aa",
		"aac",
		"aiff",
		"dvf",
		"flac",
		"gsm",
		"m4a",
		"m4b",
		"m4p",
		"mmf",
		"mp3",
		"ogg",
		"oga",
		"mogg",
		"raw",
		"tta",
		"voc",
		"vox",
		"wav",
		"wma",
		"wv",
		"webm",
	},
}

type File struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	OwnerID  uuid.UUID `json:"owner_id" db:"owner_id"`
	Url      string    `json:"url" db:"url"`
	Filename string    `json:"filename" db:"filename"`
	Path     string    `json:"path" db:"path"`
	Type     string    `json:"type" db:"type"`
	FileSize int64     `json:"file_size" db:"file_size"`
}

// String is not required by pop and may be deleted
func (f File) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Files is not required by pop and may be deleted
type Files []File

// String is not required by pop and may be deleted
func (f Files) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *File) Validate(tx *pop.Connection) (*validate.Errors, error) {

	vals := []validate.Validator{&validators.StringIsPresent{Name: "url", Field: f.Url},
		&validators.StringLengthInRange{Name: "url", Min: 20, Max: 250, Field: f.Url},
		&validators.URLIsPresent{Name: "url", Field: f.Url},

		&validators.StringIsPresent{Name: "path", Field: f.Path},
		&validators.StringLengthInRange{Name: "path", Field: f.Path, Min: 1, Max: 64},
		&validators.RegexMatch{Name: "path", Field: f.Path, Expr: "(audios|videos|images)", Message: "invalid path"},
		&validators.RegexMatch{Name: "type", Field: f.Type, Expr: "(audio|video|image)", Message: "invalid type"},
		&validators.StringIsPresent{Name: "filename", Field: f.Filename},
		&validators.StringLengthInRange{Name: "filename", Min: 10, Max: 150, Field: f.Filename},

		&validators.RegexMatch{
			Name:    "filename",
			Field:   f.Filename,
			Expr:    fmt.Sprintf("\\S\\.(%s)", strings.Join(ACCEPTABLE_FILETYPES[f.Type], "|")),
			Message: "invalid filename",
		},

		&validators.StringsMatch{Name: "url", Field: filepath.Base(f.Url), Field2: f.Filename},

		&validators.IntIsPresent{Name: "file_size", Field: int(f.FileSize)},
		&validators.IntIsLessThan{Name: "file_size", Field: int(f.FileSize), Compared: MAX_FILESIZE},

		&validators.UUIDIsPresent{Name: "owner_id", Field: f.OwnerID},
		&validators.FuncValidator{Name: "owner_id", Field: f.OwnerID.String(), Fn: func() bool {
			if err := tx.Find(&User{}, f.OwnerID); err != nil {
				return false
			}
			return true
		}},
	}

	return validate.Validate(vals...), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *File) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *File) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//Creates file entry in the DB
func (f *File) Create(tx *pop.Connection, userID uuid.UUID) (*validate.Errors, error) {
	f.OwnerID = userID
	return tx.ValidateAndCreate(f)
}
