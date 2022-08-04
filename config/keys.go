package config

import (
	"errors"
	"github.com/lestrrat-go/jwx/jwk"
	"io"
	"log"
	"net/http"
	"os"
)

var PublicKey = new(jwk.Key)

func LoadKey(path string, isRemote bool) (*jwk.Key, error) {
	var keyRaw []byte
	if !isRemote {
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
			return nil, errors.New("Remote server returned Error")
		}
		keyRaw, _ = io.ReadAll(req.Body)
	}

	key, parseErr := jwk.ParseKey(keyRaw, jwk.WithPEM(true))
	if parseErr != nil {
		return nil, parseErr
	}
	return &key, nil
}

func loadKeys() {
	var err error
	PublicKey, err = LoadKey(Config.String("keys.public"), Config.Bool("keys.remote"))
	if err != nil {
		log.Fatal(err)
	}
}
