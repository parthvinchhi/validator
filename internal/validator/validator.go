// Package validator holds the pure validation logic. It has no HTTP knowledge.
package validator

import (
	"errors"
	"sort"
)

// ErrUnknownType is returned when no validator is registered for a type.
var ErrUnknownType = errors.New("unknown validator type")

// Validator validates one kind of value.
type Validator interface {
	Type() string
	Validate(input string) Result
}

// Registry maps a type name to its validator.
type Registry struct {
	validators map[string]Validator
}

func NewRegistry(vs ...Validator) *Registry {
	r := &Registry{validators: make(map[string]Validator, len(vs))}
	for _, v := range vs {
		r.validators[v.Type()] = v
	}
	return r
}

// DefaultRegistry returns a registry with every built-in validator.
// To add a new validator: create the file, then add it to this list.
func DefaultRegistry() *Registry {
	return NewRegistry(
		Name{}, Mobile{}, Email{},
		PAN{}, Aadhaar{}, GSTIN{}, IFSC{},
	)
}

func (r *Registry) Get(t string) (Validator, bool) {
	v, ok := r.validators[t]
	return v, ok
}

// Validate looks up the validator for typ and runs it.
func (r *Registry) Validate(typ, value string) (Result, error) {
	v, ok := r.validators[typ]
	if !ok {
		return Result{}, ErrUnknownType
	}
	return v.Validate(value), nil
}

// Types returns the sorted list of registered type names.
func (r *Registry) Types() []string {
	out := make([]string, 0, len(r.validators))
	for k := range r.validators {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
