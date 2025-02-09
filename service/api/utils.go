package api

import (
	"errors"
	"io"
	"mime/multipart"
	"regexp"
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

func checkUserName(name string) (bool, error) {

	if len(name) > 30 || len(name) < 3 {
		return false, errors.New("name must be between 3 and 30 characters")
	}

	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, name)
	if err != nil {
		return false, err
	}
	if !match {
		return false, errors.New("name must contain only letters and numbers")
	}

	return match, nil

}

func checkConversationName(name string) (bool, error) {

	if len(name) > 50 || len(name) < 3 {
		return false, errors.New("name must be between 3 and 50 characters")
	}

	match, err := regexp.MatchString(`^[a-zA-Z0-9 ]+$`, name)
	if err != nil {
		return false, err
	}
	if !match {
		return false, errors.New("name must contain only letters and numbers and sapces")
	}
	return match, nil

}

func checkMessageText(text string) (bool, error) {

	if len(text) > 200 {
		return false, errors.New("message must be between 0 and 200 characters")
	}

	match, err := regexp.MatchString(`^.+$`, text)
	if err != nil {
		return false, err
	}
	if !match {
		return false, errors.New("message must respect the pattern: ^.+$")
	}
	return match, nil
}

func checkReactionText(text string) (bool, error) {

	match, err := regexp.MatchString(`^(😂|🗿|😐|👍|❤️|🔥|🎉|😢|😡)$`, text)
	if err != nil {
		return false, err
	}
	if !match {
		return false, errors.New("message must respect the pattern: ^.+$")
	}
	return match, nil
}
