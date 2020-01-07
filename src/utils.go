package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// NewRequestID supposes to generate request id by the same rule as Tablestore.
// But now we just use uuid instead.
func NewRequestID() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Sprintf("Error on generate uuid: %s", err))
	}
	return id.String()
}
