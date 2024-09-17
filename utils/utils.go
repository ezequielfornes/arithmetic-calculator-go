package utils

import (
	"errors"
	"io"
	"net/http"
)

func GetRandomString() (string, error) {
	resp, err := http.Get("https://www.random.org/strings/?num=1&len=10&digits=on&upperalpha=on&loweralpha=on&unique=on&format=plain&rnd=new")
	if err != nil {
		return "", errors.New("error making request to random.org")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error reading response body")
	}

	return string(body), nil
}
