package model

import (
	"time"

	"github.com/google/uuid"
)

type Terminal struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Gate struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Code       string    `json:"code" db:"code"`
	Name       string    `json:"name" db:"name"`
	TerminalID uuid.UUID `json:"terminal_id" db:"terminal_id"`
	GateType   string    `json:"gate_type" db:"gate_type"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Card struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CardNumber string    `json:"card_number" db:"card_number"`
	Balance    float64   `json:"balance" db:"balance"`
	Status     string    `json:"status" db:"status"`
	IssuedDate time.Time `json:"issued_date" db:"issued_date"`
	ExpiryDate time.Time `json:"expiry_date" db:"expiry_date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID              int64     `json:"id" db:"id"`
	CardID          uuid.UUID `json:"card_id" db:"card_id"`
	GateID          uuid.UUID `json:"gate_id" db:"gate_id"`
	TerminalID      uuid.UUID `json:"terminal_id" db:"terminal_id"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	Amount          *float64  `json:"amount" db:"amount"`
	BalanceAfter    float64   `json:"balance_after" db:"balance_after"`
	TransactionTime time.Time `json:"transaction_time" db:"transaction_time"`
}

type FareMatrix struct {
	OriginTerminalID      uuid.UUID `json:"origin_terminal_id" db:"origin_terminal_id"`
	DestinationTerminalID uuid.UUID `json:"destination_terminal_id" db:"destination_terminal_id"`
	FareAmount            float64   `json:"fare_amount" db:"fare_amount"`
	IsActive              bool      `json:"is_active" db:"is_active"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

type Admin struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
}
