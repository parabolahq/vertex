package config

import (
	"errors"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"io"
	"log"
	"net/http"
	"os"
)

var PublicKey = new(jwk.Key)
var KeySet = jwk.NewSet()

func LoadKey(path string, isRemote bool, alg jwa.SignatureAlgorithm) (*jwk.Key, error) {
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
			return nil, errors.New("remote server returned Error")
		}
		keyRaw, _ = io.ReadAll(req.Body)
	}

	key, parseErr := jwk.ParseKey(keyRaw, jwk.WithPEM(true))
	key.Set(jwk.AlgorithmKey, alg)
	if parseErr != nil {
		return nil, parseErr
	}
	return &key, nil
}

func LoadKeys() {
	var err error
	PublicKey, err = LoadKey(Config.String("keys.public"), Config.Bool("keys.remote"), jwa.RS256)
	if err != nil {
		log.Fatal(err)
	}
	KeySet.Add(*PublicKey)
}
