package databasetest

import (
	"bytes"
	"encoding/json"
)

// record はレコードを表す
// SELECT文の列の並びを維持してJSONオブジェクトにエンコードするために使用する
type record struct {
	columns []string
	values  []any
}

func (r record) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	buf.WriteString("{")
	for i, column := range r.columns {
		if i > 0 {
			buf.WriteString(",")
		}
		if err := enc.Encode(column); err != nil {
			return nil, err
		}
		buf.WriteString(":")
		if err := enc.Encode(r.values[i]); err != nil {
			return nil, err
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
