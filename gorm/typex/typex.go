package typex

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
)

type StringList []string

func (p StringList) Value() (driver.Value, error) {
	if p == nil || len(p) == 0 {
		return "[]", nil
	}
	return json.Marshal(p)
}

func (p *StringList) Scan(input interface{}) error {
	if input == nil {
		*p = make(StringList, 0)
		return nil
	}

	var data []byte
	switch v := input.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", input)
	}

	if len(data) == 0 || string(data) == "[]" {
		*p = make(StringList, 0)
		return nil
	}

	if err := json.Unmarshal(data, p); err != nil {
		return fmt.Errorf("failed to unmarshal StringList: %v", err)
	}
	return nil
}

func (p StringList) GetOne() string {
	if len(p) == 0 {
		return ""
	}
	return p[0]
}

func (p StringList) Contains(t string) bool {
	for _, v := range p {
		if t == v {
			return true
		}
	}
	return false
}

type Uint64List []uint64

func (p Uint64List) Value() (driver.Value, error) {
	if p == nil || len(p) == 0 {
		return "[]", nil
	}
	return json.Marshal(p)
}

func (p *Uint64List) Scan(input interface{}) error {
	err := json.Unmarshal(input.([]byte), p)
	if err != nil {
		*p = make([]uint64, 0)
	}
	return nil
}

func (p Uint64List) GetOne() uint64 {
	if len(p) == 0 {
		return 0
	}
	return p[0]
}

func (p Uint64List) Contains(t uint64) bool {
	for _, v := range p {
		if t == v {
			return true
		}
	}
	return false
}

type Decimal struct {
	decimal.Decimal
}

func (d *Decimal) Scan(value interface{}) error {
	if value == nil {
		d.Decimal = decimal.Zero
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return d.Decimal.UnmarshalText(v)
	case string:
		return d.Decimal.UnmarshalText([]byte(v))
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}
