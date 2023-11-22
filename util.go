package sqlgo

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Support utilities
func combineSQL(firstSQL string, secondSQL string) string {
	if secondSQL == "" {
		return firstSQL
	} else if firstSQL != "" {
		firstSQL = fmt.Sprintf("%s ", firstSQL)
	}
	return fmt.Sprintf("%s%s", firstSQL, secondSQL)
}

func hash(values ...interface{}) string {
	var sums string
	for _, value := range values {
		if jsonVal, err := json.Marshal(value); err == nil {
			sums += string(jsonVal)
		}
	}
	h := sha1.New()
	h.Write([]byte(sums))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
