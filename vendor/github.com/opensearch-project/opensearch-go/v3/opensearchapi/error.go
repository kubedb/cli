// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.

package opensearchapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error vars
var (
	ErrUnexpectedEmptyBody    = errors.New("body is unexpectedly empty")
	ErrReadBody               = errors.New("failed to read body")
	ErrJSONUnmarshalBody      = errors.New("failed to json unmarshal body")
	ErrUnknownOpensearchError = errors.New("opensearch error response could not be parsed as error")
)

// Error represents the Opensearch API error response
type Error struct {
	Err    Err `json:"error"`
	Status int `json:"status"`
}

// Err represents the error of an API error response
type Err struct {
	RootCause []RootCause `json:"root_cause"`
	Type      string      `json:"type"`
	Reason    string      `json:"reason"`
	Index     string      `json:"index,omitempty"`
	IndexUUID string      `json:"index_uuid,omitempty"`
}

// RootCause represents the root_cause of an API error response
type RootCause struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Index     string `json:"index,omitempty"`
	IndexUUID string `json:"index_uuid,omitempty"`
}

// StringError represnets an Opensearch API Error with a string as error
type StringError struct {
	Err    string `json:"error"`
	Status int    `json:"status"`
}

// Error returns a string
func (e Error) Error() string {
	return fmt.Sprintf("status: %d, type: %s, reason: %s, root_cause: %s", e.Status, e.Err.Type, e.Err.Reason, e.Err.RootCause)
}

// Error returns a string
func (e StringError) Error() string {
	return fmt.Sprintf("status: %d, error: %s", e.Status, e.Err)
}

// UnmarshalJSON is a custom unmarshal function for Error returning custom errors in special cases
func (e *Error) UnmarshalJSON(b []byte) error {
	var dummy struct {
		Err    json.RawMessage `json:"error"`
		Status int             `json:"status"`
	}
	if err := json.Unmarshal(b, &dummy); err != nil {
		return fmt.Errorf("%w: %s", err, b)
	}

	if len(dummy.Err) == 0 {
		return fmt.Errorf("%w: %s", ErrUnknownOpensearchError, b)
	}

	var osErr Err
	if err := json.Unmarshal(dummy.Err, &osErr); err != nil {
		return StringError{Status: dummy.Status, Err: string(dummy.Err)}
	}
	if dummy.Status == 0 || (osErr.Type == "" && osErr.Reason == "") {
		return fmt.Errorf("%w: %s", ErrUnknownOpensearchError, b)
	}

	e.Err = osErr
	e.Status = dummy.Status

	return nil
}
