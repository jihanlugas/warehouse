package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func UploadImage(path string, file *multipart.FileHeader) error {
	var err error
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return err
}

func RemoveImage(path string) error {
	var err error

	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return err
}
