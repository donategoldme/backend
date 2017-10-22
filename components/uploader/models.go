package uploader

import "time"

const uploadDir = "/var/uploads"
const picDir = "pic"
const soundDir = "sound"

// type Category struct {
// 	ID     uint   `json:"id"`
// 	UserID uint   `json:"user_id"`
// 	Name   string `json:"name"`
// }

type Img struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Path   string `json:"path"`
	Name   string `json:"name"`

	DeletedAt *time.Time `json:"-"`
	// Category  uint       `json:"category"`
}

type Sound struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Path   string `json:"path"`
	Name   string `json:"name"`

	DeletedAt *time.Time `json:"-"`
	// Category  uint       `json:"category"`
}
