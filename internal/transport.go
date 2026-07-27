package whitearchive

import (
	"bufio"
	"bytes"
	"encoding/json"
	"time"
)

type TransportObject struct {
	Path       string    `json:"path"`
	Hash       []byte    `json:"hash"`
	LastUpdate time.Time `json:"update"`
}

func UnmarshalSnapshot(data []byte) (Snapshot, error) {
	var objects []TransportObject
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var obj TransportObject
		if err := json.Unmarshal(scanner.Bytes(), &obj); err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return trasnsportsToSnapshot(objects)
}

func MarshalSnapshot(snapshot Snapshot) ([]byte, error) {
	objects, err := snapshotToTransports(snapshot)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, obj := range objects {
		if err := enc.Encode(obj); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func trasnsportsToSnapshot(objects []TransportObject) (Snapshot, error) {
	snapshot := Snapshot{}
	for _, obj := range objects {
		snapshot[obj.Path] = Data{
			Update: obj.LastUpdate,
			Hash:   obj.Hash,
		}
	}

	return snapshot, nil
}

func snapshotToTransports(snapshot Snapshot) ([]TransportObject, error) {
	transport := make([]TransportObject, 0, len(snapshot))
	for k, v := range snapshot {
		transport = append(transport, TransportObject{
			Path:       k,
			Hash:       v.Hash,
			LastUpdate: v.Update,
		})
	}
	return transport, nil
}
