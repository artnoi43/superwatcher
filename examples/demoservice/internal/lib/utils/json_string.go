package utils

import "encoding/json"

func StringJSON(e any) string {
	b, err := json.Marshal(e)
	if err != nil {
		panic("failed to marshal ENSArtifact to JSON")
	}

	return string(b)
}
