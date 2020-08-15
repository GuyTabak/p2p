package main

import (
	"bytes"
	"encoding/gob"
)

func serialize(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	err := enc.Encode(data)
	if err == nil {
		return buff.Bytes()
	}

	return []byte{}
}

func deserialize(source []byte, dest interface{}) error {

	dec := gob.NewDecoder(bytes.NewBuffer(source))
	err := dec.Decode(dest)

	return err
}
