package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSON json.RawMessage

// Scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result json.RawMessage
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)

	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 || string(j) == "null" {
		return nil, nil
	}

	return json.RawMessage(j).MarshalJSON()
}

func (j *JSON) UnmarshalJson(data []byte) error {
	if j == nil {
		return errors.New("null pointer exception")
	}
	*j = append((*j)[0:0], data...)

	return nil
}

func (j JSON) MarshalJson() []byte {
	if j == nil {
		return []byte("null")
	}

	return j
}
