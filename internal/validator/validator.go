package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator struct will be used to hold the validation errors
type Validator struct {
	Errors map[string]string
}

// New function will be used to create a new instance of Validator struct
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid function will be used to check if the validation errors map is empty or not (if it is empty then return true)
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError function will be used to add a new error message to the validation errors map
func (v *Validator) AddError(field, message string) {
	if _, ok := v.Errors[field]; !ok {
		v.Errors[field] = message
	}
}

// Check function will be used to check if the given value is empty or not (if it is empty then add an error message to the validation errors map)
func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// In function will be used to check if the given value is in the list of valid values or not (if it is not in the list then add an error message to the validation errors map)
func (v *Validator) In(value string, values ...string) bool {
	for i := range values {
		if value == values[i] {
			return true
		}
	}
	return false
}

// Matches function will be used to check if the given value matches the regular expression or not (if it does not match then add an error message to the validation errors map)
func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique function will be used to check if the given value is unique in the list of values or not (if it is not unique then add an error message to the validation errors map)
func (v *Validator) Unique(value string, values ...string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(uniqueValues) == len(values)
}
