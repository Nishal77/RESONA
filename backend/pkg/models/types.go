package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringArray handles PostgreSQL varchar[] columns
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	quoted := make([]string, len(s))
	for i, v := range s {
		quoted[i] = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(quoted, ",") + "}", nil
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	str, ok := value.(string)
	if !ok {
		b, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("cannot scan type %T into StringArray", value)
		}
		str = string(b)
	}
	str = strings.Trim(str, "{}")
	if str == "" {
		*s = StringArray{}
		return nil
	}
	parts := strings.Split(str, ",")
	result := make(StringArray, len(parts))
	for i, p := range parts {
		result[i] = strings.Trim(strings.TrimSpace(p), `"`)
	}
	*s = result
	return nil
}

// APIResponse is the standard response envelope for all endpoints
type APIResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
	Error      string      `json:"error,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
}

type PaginationMeta struct {
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	Total   int64 `json:"total"`
	HasMore bool  `json:"has_more"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}
