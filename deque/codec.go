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

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

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

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	d := Data{
		Q: q.queue,
		P: q.processing,
	}

	return json.Marshal(d)
}
