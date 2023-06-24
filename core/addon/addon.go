package addon

import (
	"fmt"
	"main/core/whats"
	"reflect"
)

// AddOnList List of add-ons that registered
var AddOnList = []*AddOn{}

// NewAddOn Create new Addon with name and drescription
func NewAddOn(n, d string) AddOn {
	return AddOn{
		Name:        n,
		Description: d,
		Features:    []*Feature{},
	}
}

// NewAddOnRegistered Create new Addon with name and drescription also register it
func NewAddOnRegistered(n, d string) *AddOn {
	var p = NewAddOn(n, d)

	AddOnList = append(AddOnList, &p)

	return &p
}

// Calls Trigering execution by the addon
func Calls(i interface{}, gwb *whats.GoWhatsBot) map[string][]error {
	var res = map[string][]error{}

	for _, add := range AddOnList {
		if add.Disable {
			continue
		}

		for key, errs := range add.ExecuteFeatures(i, gwb) {
			res[key] = errs
		}
	}

	return res
}

// AOValidator Type for Add-on Validator
type AOValidator func(*AddOn, interface{}, *whats.GoWhatsBot) error

// AddOn Type for AddOn
type AddOn struct {
	Name        string
	Description string
	Disable     bool
	Tags        []string
	Features    []*Feature
	Validator   AOValidator
}

// ExecuteFeatures Execure every feature by it own addon
func (p *AddOn) ExecuteFeatures(i interface{}, gwb *whats.GoWhatsBot) map[string][]error {
	var resp = map[string][]error{}

	if p.Validator != nil {
		if err := p.Validator(p, i, gwb); err != nil {
			return map[string][]error{p.Name: {err}}
		}
	}

	for _, c := range p.Features {

		if c.Disable {
			continue
		}

		var key = fmt.Sprintf("%s>%s", p.Name, c.Name)
		var typeGranted bool

		if c.Expecting != nil {
			typeGranted = reflect.TypeOf(i) == reflect.TypeOf(c.Expecting)
		}

		if typeGranted {
			var validGranted bool

			if c.Validator != nil {
				if err := c.Validator(c, i, gwb); err == nil {
					validGranted = true
				} else {
					resp[key] = append(resp[key], err)
				}
			} else {
				validGranted = true
			}

			if validGranted {

				if err := c.Execute(p, c, i, gwb); err != nil {
					resp[key] = append(resp[key], err)
				}
			}
		}
	}

	return resp
}

// AddFeatures Add features
func (p *AddOn) AddFeatures(c ...*Feature) {
	p.Features = append(p.Features, c...)
}

// FeatureValidator Type for feature validator
type FeatureValidator func(*Feature, interface{}, *whats.GoWhatsBot) error

// FeatureExecute Type for feature executor
type FeatureExecute func(*AddOn, *Feature, interface{}, *whats.GoWhatsBot) error

// Feature Type for a feature
type Feature struct {
	Name        string
	Description string
	Usage       string
	Patterns    []string
	Tags        []string
	Disable     bool
	Expecting   interface{}
	Execute     FeatureExecute
	Validator   FeatureValidator
}

// IsMy Check if pattern that given is exists in current feature
func (f *Feature) IsMy(p string) bool {
	for _, v := range f.Patterns {
		if v == p {
			return true
		}
	}
	return false
}
