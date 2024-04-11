package db

import (
	e "BannerFlow/internal/domain/errors"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Attrs map[string]interface{}

func (a *Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("%w: type assertion to []byte failed", e.ErrorInternal)
	}

	return json.Unmarshal(b, &a)
}
