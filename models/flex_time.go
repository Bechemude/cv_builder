package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// FlexTime is a nullable time.Time that accepts "", null, and {"T":"..."} from JSON.
type FlexTime struct {
	T *time.Time
}

func (ft *FlexTime) UnmarshalJSON(data []byte) error {
	s := string(data)

	// null or empty string
	if s == "null" || s == `""` || s == "" {
		ft.T = nil
		return nil
	}

	// plain string: "2024-01-01T00:00:00Z"
	if s[0] == '"' {
		inner := s[1 : len(s)-1]
		// handle "null", "none", "present", "current" — LLM sometimes returns these as strings
		if inner == "" || inner == "null" || inner == "none" || inner == "present" || inner == "current" {
			ft.T = nil
			return nil
		}
		t, err := time.Parse(time.RFC3339, inner)
		if err != nil {
			return fmt.Errorf("FlexTime: cannot parse %q: %w", inner, err)
		}
		ft.T = &t
		return nil
	}

	// object: {"T": "2024-01-01T00:00:00Z"} or {"T": null}
	if s[0] == '{' {
		var obj struct {
			T *string `json:"T"`
		}
		if err := json.Unmarshal(data, &obj); err != nil {
			return fmt.Errorf("FlexTime: cannot parse object: %w", err)
		}
		if obj.T == nil || *obj.T == "" {
			ft.T = nil
			return nil
		}
		t, err := time.Parse(time.RFC3339, *obj.T)
		if err != nil {
			return fmt.Errorf("FlexTime: cannot parse T field %q: %w", *obj.T, err)
		}
		ft.T = &t
		return nil
	}

	return fmt.Errorf("FlexTime: unexpected value %s", s)
}

func (ft FlexTime) MarshalJSON() ([]byte, error) {
	if ft.T == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + ft.T.Format(time.RFC3339) + `"`), nil
}

// GORM support

func (ft FlexTime) Value() (driver.Value, error) {
	if ft.T == nil {
		return nil, nil
	}
	return *ft.T, nil
}

func (ft *FlexTime) Scan(value interface{}) error {
	if value == nil {
		ft.T = nil
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("FlexTime: cannot scan %T", value)
	}
	ft.T = &t
	return nil
}
