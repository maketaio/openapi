package fields

import (
	"bytes"
	"encoding/json"
)

type State[T any] struct {
	v     *T
	isSet bool
}

func (s *State[T]) Set(x T) {
	s.v = &x
	s.isSet = true
}

func (s *State[T]) SetNull() {
	s.v = nil
	s.isSet = true
}

func (s *State[T]) Unset() {
	s.v = nil
	s.isSet = false
}

func (s State[T]) IsPresent() bool {
	return s.isSet
}

func (s State[T]) IsNull() bool {
	return s.isSet && s.v == nil
}

func (s State[T]) Value() (T, bool) {
	if !s.isSet || s.v == nil {
		var z T
		return z, false
	}

	return *s.v, true
}

func (s *State[T]) UnmarshalJSON(data []byte) error {
	s.isSet = true

	trim := bytes.TrimSpace(data)
	if bytes.Compare(trim, []byte("null")) == 0 {
		s.v = nil
		return nil
	}

	var v T
	if err := json.Unmarshal(trim, &v); err != nil {
		return err
	}

	s.v = &v
	return nil
}

func (s State[T]) MarshalJSON() ([]byte, error) {
	// If the field is tagged with `omitempty` and IsZero() on the wrapper returns true,
	// encoding/json will *omit* the field entirely and NEVER call this method.
	if !s.isSet {
		return []byte("null"), nil
	}

	if s.v == nil {
		return []byte("null"), nil
	}

	return json.Marshal(*s.v)
}
