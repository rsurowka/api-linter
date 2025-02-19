// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// RuleName is an identifier for a rule. Allowed characters include a-z, 0-9, -.
//
// The namespace separator :: is allowed between RuleName segments
// (for example, my-namespace::my-rule).
type RuleName string

const nameSeparator string = "::"

var ruleNameValidator = regexp.MustCompile("^([a-z0-9][a-z0-9-]*(::[a-z0-9][a-z0-9-]*)?)+$")

// getRuleGroup takes an AIP number and returns the appropriate group.
// It panics if no group is found.
func getRuleGroup(aip int) string {
	if aip > 0 && aip < 1000 {
		return "core"
	}
	panic(fmt.Sprintf("Invalid AIP number %d: no available group.", aip))
}

// NewRuleName creates a RuleName from an AIP number and a unique name within
// that AIP.
func NewRuleName(aip int, name string) RuleName {
	return RuleName(strings.Join([]string{
		getRuleGroup(aip),
		fmt.Sprintf("%04d", aip),
		name,
	}, nameSeparator))
}

// IsValid checks if a RuleName is syntactically valid.
func (r RuleName) IsValid() bool {
	return r != "" && ruleNameValidator.Match([]byte(r))
}

// HasPrefix returns true if r contains prefix as a namespace. prefix parameters can be "::" delimited
// or specified as independent parameters.
// For example:
//
// r := NewRuleName("foo", "bar", "baz")   // string(r) == "foo::bar::baz"
//
// r.HasPrefix("foo::bar")          == true
// r.HasPrefix("foo", "bar")        == true
// r.HasPrefix("foo", "bar", "baz") == true   // matches the entire string
// r.HasPrefix("foo", "ba")         == false  // prefix must end on a delimiter
func (r RuleName) HasPrefix(prefix ...string) bool {
	s := strings.Join(prefix, nameSeparator)
	return string(r) == s || strings.HasPrefix(string(r), s+nameSeparator)
}
