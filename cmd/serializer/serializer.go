package serializer

import (
	"encoding/gob"
	"os"
)

type Serializer struct {
	path string
}

func NewSerializer(path string) *Serializer {
	return &Serializer{path: path}
}

func (se *Serializer) Serialize(label string, obj any) error {
	file, err := os.Create(se.path + label + ".bin")
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(obj)
}

func (se *Serializer) Deserialize(label string, obj any) error {
	file, err := os.Open(se.path + label + ".bin")
	if err != nil {
		return err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	return dec.Decode(obj)
}
