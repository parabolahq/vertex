package config

import (
	"errors"
	"github.com/lestrrat-go/jwx/jwk"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var PublicKey = new(jwk.Key)
var KeySet = jwk.NewSet()

func LoadKey(path string) (*jwk.Key, error) {
	var keyRaw []byte
	if !strings.HasPrefix(path, "http") {
		var readErr error
		keyRaw, readErr = os.ReadFile(path)
		if readErr != nil {
			return nil, readErr
		}
	} else {
		req, requestError := http.Get(path)
		if requestError != nil {
			return nil, requestError
		} else if req.StatusCode != http.StatusOK {
			return nil, errors.New("remote server returned Error")
		}
		keyRaw, _ = io.ReadAll(req.Body)
	}

	key, parseErr := jwk.ParseKey(keyRaw)
	if parseErr != nil {
		return nil, parseErr
	}
	return &key, nil
}

func LoadKeys() {
	var err error
	for _, s := range Config.Strings("keys") {
		PublicKey, err = LoadKey(s)
		if err != nil {
			log.Fatal(err)
		}
		KeySet.Add(*PublicKey)
	}
}
