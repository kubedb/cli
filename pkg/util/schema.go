package util

import (
	"k8s.io/kubernetes/pkg/kubectl/validation"
)

type ConjunctiveSchema []validation.Schema

func (c *ConjunctiveSchema) ValidateBytes(data []byte) error {
	return nil
}

func Validator() validation.Schema {
	return validation.ConjunctiveSchema{
		validation.NoDoubleKeySchema{},
	}
}
