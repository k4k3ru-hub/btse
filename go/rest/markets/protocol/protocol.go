// Package protocol defines BTSE Markets v1 REST wire envelopes.
package protocol

import "fmt"

type Envelope[T any] struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Time    int64  `json:"time"`
	Data    T      `json:"data"`
}

// Validate validates a BTSE Markets response envelope.
//
// Version:
//   - 2026-08-26: Added.
func (e Envelope[T]) Validate() error {
	if !e.Success || e.Code != 1 {
		return fmt.Errorf("failed to validate markets response: API rejected request: code=%d message=%q", e.Code, e.Message)
	}
	return nil
}
