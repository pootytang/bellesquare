package models

import (
	"database/sql/driver"
	"fmt"
)

type PaymentStatus int

const (
	Unpaid  PaymentStatus = 0
	Pending PaymentStatus = 1
	Paid    PaymentStatus = 2
)

func (p PaymentStatus) String() string {
	return [...]string{"unpaid", "pending", "paid"}[p]
}

// Scan implements the sql.Scanner interface for database retrieval
func (p *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		*p = Unpaid
		return nil
	}
	intVal, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid type for PaymentStatus: %T", value)
	}
	*p = PaymentStatus(intVal)
	return nil
}

// Value implements the driver.Valuer interface for database storage
func (p PaymentStatus) Value() (driver.Value, error) {
	return int64(p), nil
}
