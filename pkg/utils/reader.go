package utils

import (
	"encoding/json"
	"io"
	"log"
	"strings"
)

func StructToReader(data any) io.Reader {

	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf("Failed to marshall json %v", err)
	}

	return io.NopCloser(strings.NewReader(string(jsonData)))

}
