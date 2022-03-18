package utils

import (
	"encoding/gob"
	"os"
)

func SaveObj(name string, obj interface{}) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	return enc.Encode(obj)
}

func LoadObj(name string) (interface{}, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)
	var u interface{}
	err2 := dec.Decode(&u)
	return u, err2
}
