package uploader

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func genName() string {
	now := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(now, 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func checkPath(userId uint, fname string, dir string) (string, string, error) {
	userPath := fmt.Sprintf("%s/%d/", uploadDir, userId)
	if err := createIfNotExistDir(userPath); err != nil {
		return "", "", err
	}
	userPath += dir
	if err := createIfNotExistDir(userPath); err != nil {
		return "", "", err
	}
	if fname == "" {
		return "", "", errors.New("No file name")
	}
	path := userPath + "/" + fname
	if _, err := os.Stat(path); err == nil {
		ext := filepath.Ext(fname)
		basename := strings.TrimSuffix(fname, ext)
		fname = basename + "_" + genName() + ext
		path = userPath + "/" + fname
	}
	return path, fname, nil
}

func createIfNotExistDir(path string) error {
	err := os.Mkdir(path, 0777)
	if os.IsExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func getListFiles(dir string, userId uint) []string {
	path := fmt.Sprintf("%s/%d/%s/", uploadDir, userId, dir)
	files, _ := ioutil.ReadDir(path)
	f := make([]string, len(files))
	for i, file := range files {
		f[i] = file.Name()
	}
	return f
}
