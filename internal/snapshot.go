package whitearchive

import (
	"bytes"
	"time"
)

type Snapshot map[string]Data

type Data struct {
	Hash   []byte
	Update time.Time
}

func diffs(primary, secondary Snapshot) Snapshot {
	result := Snapshot{}
	for path, data := range secondary {
		if p, exists := primary[path]; !exists || !bytes.Equal(p.Hash, data.Hash) {
			result[path] = data
		}
	}
	return result
}
