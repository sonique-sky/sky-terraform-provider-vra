package api

import "fmt"

type Error struct {
	Errors []struct {
		Code          int    `json:"code"`
		Message       string `json:"message"`
		SystemMessage string `json:"systemMessage"`
	} `json:"errors"`
}

func (e Error) Error() error {
	return fmt.Errorf("vRealize API: %+v", e.Errors)
}

func (e Error) isEmpty() bool {
	return len(e.Errors) == 0
}
