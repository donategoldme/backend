package uploader

import (
	"io"
	"os"

	"fmt"

	"donategold.me/components/users"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Get("/img", getListImgs)
	api.Get("/sound", getListSounds)
	api.Post("/img", saveImg)
	api.Post("/sound", saveSound)
	api.Delete("/img", deleteImg)
	api.Delete("/sound", deleteSound)
}

func getListImgs(c *iris.Context) {
	userId := c.Get("token").(users.AccessToken).User.ID
	c.JSON(200, getListFiles(picDir, userId))
}

func saveImg(c *iris.Context) {
	file, info, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, "Need file")
		return
	}
	ct := info.Header.Get("Content-Type")
	if ct[:5] != "image" {
		c.JSON(400, "Need image file")
		return
	}
	//if info.Header.Get("Content-Type")[:6]
	fname := info.Filename
	userID := c.Get("token").(users.AccessToken).User.ID
	path, fname, err := checkPath(userID, fname, picDir)
	if err != nil {
		c.EmitError(400)
		return
	}
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	defer out.Close()
	fo := file
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	defer fo.Close()
	io.Copy(out, fo)
	c.JSON(200, fname)
}

func getListSounds(c *iris.Context) {
	userId := c.Get("token").(users.AccessToken).User.ID
	c.JSON(200, getListFiles(soundDir, userId))
}

func saveSound(c *iris.Context) {
	file, info, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	ct := info.Header.Get("Content-Type")
	if ct[:5] != "audio" {
		c.JSON(400, "Need audio file")
		return
	}
	//if info.Header.Get("Content-Type")[:6]
	fname := info.Filename
	userId := c.Get("token").(users.AccessToken).User.ID
	path, fname, err := checkPath(userId, fname, soundDir)
	if err != nil {
		c.EmitError(400)
		return
	}
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	defer out.Close()
	fo := file
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	defer fo.Close()
	io.Copy(out, fo)
	c.JSON(200, fname)
}
func deleteImg(c *iris.Context) {
	userId := c.Get("token").(users.AccessToken).User.ID
	filename := c.URLParam("file")
	if filename == "" {
		c.JSON(400, "Filename required")
		return
	}
	path := fmt.Sprintf("%s/%d/%s/%s", uploadDir, userId, picDir, filename)
	os.Remove(path)
	c.JSON(200, filename)
}

func deleteSound(c *iris.Context) {
	userId := c.Get("token").(users.AccessToken).User.ID
	filename := c.URLParam("file")
	if filename == "" {
		c.JSON(400, "Filename required")
		return
	}
	path := fmt.Sprintf("%s/%d/%s/%s", uploadDir, userId, soundDir, filename)
	os.Remove(path)
	c.JSON(200, filename)
}
