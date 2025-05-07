package enums

import (
	"errors"
	"fmt"
)

type Priority string

const (
	Low      Priority = "Low"
	Medium   Priority = "Medium"
	High     Priority = "High"
	Critical Priority = "Critical"
)

func ValidatePriority(p Priority) error {
	switch p {
	case Low, Medium, High, Critical:
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported priority: %v", p))
	}
}
