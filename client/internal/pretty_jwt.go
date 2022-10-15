package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
)

func GetJSONString(token string) string {
	parts := strings.Split(token, ".")
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Println(err)
	}

	buf := &bytes.Buffer{}
	if err := json.Indent(buf, decoded, "", "\t"); err != nil {
		log.Println(err)
	}
	return buf.String()
}
