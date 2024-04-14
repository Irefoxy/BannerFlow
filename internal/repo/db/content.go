package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Attrs map[string]interface{}

func (a *Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	switch data := value.(type) {
	case []byte:
		return json.Unmarshal(data, &a)
	case string:
		return json.Unmarshal([]byte(data), &a)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}
