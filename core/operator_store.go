package core

import (
	"errors"
	"fmt"
)

// ErrExistsDriver can not register same driver_default twice.
var ErrExistsDriver = errors.New("driver_default already exists")

// OperatorStore store all supported drivers.
type OperatorStore struct {
	operator map[string]Operator
}

// Get driver_default by name.
func (s *OperatorStore) Get(name string) (Operator, error) {
	if s.operator == nil {
		return nil, errors.New("no valid driver_default")
	}

	storage, ok := s.operator[name]
	if !ok {
		return nil, fmt.Errorf("unsupported driver_default: %s", name)
	}

	return storage, nil
}

// Register a new driver_default.
func (s *OperatorStore) Register(name string, storage Operator) error {
	if s.operator == nil {
		s.operator = make(map[string]Operator)
	}

	if _, ok := s.operator[name]; ok {
		return ErrExistsDriver
	}

	s.operator[name] = storage

	return nil
}
