package deque

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	Q []interface{} `json:"Q"`
	P []interface{} `json:"P"`
}

func (q *simpleDeque) Decode(rawData []byte) error {

	var d Data
	if err := json.Unmarshal(rawData, d); err != nil {
		return fmt.Errorf("decode: %v", err)
	}

	q.queue = append(q.queue, d.Q)
	for _, p := range d.P {
		q.processing.add(p)
	}

	return nil
}

func (q *simpleDeque) Encode() ([]byte, error) {

	d := Data{
		Q: q.queue,
	}

	for o, _ := range q.processing {
		d.P = append(d.P, o)
	}

	return json.Marshal(d)
}
