package testcase

import (
	"encoding/json"
)

type Element struct {
	In  string
	Out string
}

type Testcase struct {
	// metadata should be here
	Data []Element
}

func (t *Testcase) Count() int {
	return len(t.Data)
}

func (t Testcase) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Testcase) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
