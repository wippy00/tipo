package api

import (
	"io"
	"mime/multipart"
)

func validateFile(photo_multipart multipart.File, handler *multipart.FileHeader, err error) ([]byte, error) {
	var photo []byte

	if err != nil {
		return photo, err
	}

	if handler.Size > 1024*1024 {
		return photo, err
	}

	var erro error
	photo, erro = io.ReadAll(photo_multipart)
	if erro != nil {

		return photo, erro
	}
	return photo, nil
}
