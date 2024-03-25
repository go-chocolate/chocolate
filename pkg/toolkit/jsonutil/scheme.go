package jsonutil

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

const (
	INTEGER = "integer"
	FLOAT   = "float"
	STRING  = "string"
	BOOLEAN = "boolean"
	OBJECT  = "object"
	ARRAY   = "array"
)

type Scheme struct {
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
	Regexp      string    `json:"regexp"`
	Hint        string    `json:"hint"`
	Type        string    `json:"type"`
	Schemes     []*Scheme `json:"schema"`
}

func (s *Scheme) Validate() error {
	if s.Regexp != "" {
		matched, err := regexp.MatchString(s.Regexp, s.Value)
		if err != nil {
			return err
		}
		if !matched {
			if s.Hint != "" {
				return fmt.Errorf("%s: %s", s.Name, s.Hint)
			} else {
				return fmt.Errorf("invalid value for %s", s.Name)
			}
		} else {
			return nil
		}
	}
	if s.Value == "" {
		return nil
	}
	switch s.Type {
	case INTEGER:
		if _, err := strconv.Atoi(s.Value); err != nil {
			return fmt.Errorf("invalid integer value %s", s.Value)
		}
	case FLOAT:
		if _, err := strconv.ParseFloat(s.Value, 64); err != nil {
			return fmt.Errorf("invalid double value %s", s.Value)
		}
	case BOOLEAN:
		if _, err := strconv.ParseBool(s.Value); err != nil {
			return fmt.Errorf("invalid boolean value %s", s.Value)
		}
	case OBJECT:
		if len(s.Schemes) > 0 {
			if err := s.Schemes[0].Validate(); err != nil {
				return err
			}
		}
	case ARRAY:
		for _, item := range s.Schemes {
			if err := item.Validate(); err != nil {
				return err
			}
		}
	case STRING:

	default:
		return fmt.Errorf("invalid type %s", s.Type)
	}
	return nil
}

func (s *Scheme) Build() []byte {
	buf := bytes.NewBuffer(nil)
	s.build(buf)
	return buf.Bytes()
}

func (s *Scheme) build(buf *bytes.Buffer) {
	switch s.Type {
	case INTEGER:
		if s.Value == "" {
			buf.WriteString("0")
		} else {
			buf.WriteString(s.Value)
		}
	case FLOAT:
		if s.Value == "" {
			buf.WriteString("0.0")
		} else {
			buf.WriteString(s.Value)
		}
	case STRING:
		buf.WriteString(fmt.Sprintf("\"%s\"", s.Value))
	case BOOLEAN:
		if s.Value == "" {
			buf.WriteString("false")
		} else {
			buf.WriteString(s.Value)
		}
	case OBJECT:
		buf.WriteString("{")
		for i, field := range s.Schemes {
			buf.WriteString(fmt.Sprintf("\"%s\":", field.Name))
			field.build(buf)
			if i < len(s.Schemes)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("}")
	case ARRAY:
		buf.WriteString("[")
		for i, field := range s.Schemes {
			field.build(buf)
			if i < len(s.Schemes)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("]")
	}
}
