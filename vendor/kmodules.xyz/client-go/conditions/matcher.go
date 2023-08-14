/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conditions

import (
	"fmt"

	kmapi "kmodules.xyz/client-go/api/v1"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// MatchConditions returns a custom matcher to check equality of conditionsapi.Conditions.
func MatchConditions(expected kmapi.Conditions) types.GomegaMatcher {
	return &matchConditions{
		expected: expected,
	}
}

type matchConditions struct {
	expected kmapi.Conditions
}

func (m matchConditions) Match(actual interface{}) (success bool, err error) {
	elems := []interface{}{}
	for _, condition := range m.expected {
		elems = append(elems, MatchCondition(condition))
	}

	return gomega.ConsistOf(elems).Match(actual)
}

func (m matchConditions) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected\n\t%#v\nto match\n\t%#v\n", actual, m.expected)
}

func (m matchConditions) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected\n\t%#v\nto not match\n\t%#v\n", actual, m.expected)
}

// MatchCondition returns a custom matcher to check equality of conditionsapi.Condition.
func MatchCondition(expected kmapi.Condition) types.GomegaMatcher {
	return &matchCondition{
		expected: expected,
	}
}

type matchCondition struct {
	expected kmapi.Condition
}

func (m matchCondition) Match(actual interface{}) (success bool, err error) {
	actualCondition, ok := actual.(kmapi.Condition)
	if !ok {
		return false, fmt.Errorf("actual should be of type Condition")
	}

	ok, err = gomega.Equal(m.expected.Type).Match(actualCondition.Type)
	if !ok {
		return ok, err
	}
	ok, err = gomega.Equal(m.expected.Status).Match(actualCondition.Status)
	if !ok {
		return ok, err
	}
	ok, err = gomega.Equal(m.expected.Severity).Match(actualCondition.Severity)
	if !ok {
		return ok, err
	}
	ok, err = gomega.Equal(m.expected.Reason).Match(actualCondition.Reason)
	if !ok {
		return ok, err
	}
	ok, err = gomega.Equal(m.expected.Message).Match(actualCondition.Message)
	if !ok {
		return ok, err
	}

	return ok, err
}

func (m matchCondition) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected\n\t%#v\nto match\n\t%#v\n", actual, m.expected)
}

func (m matchCondition) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("expected\n\t%#v\nto not match\n\t%#v\n", actual, m.expected)
}
