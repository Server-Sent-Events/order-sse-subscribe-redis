package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// OrderStatusType is an exported
type OrderStatusType int64

// DRAFT OrderStatusType is an exported
const (
	DRAFT OrderStatusType = iota
	ENTERED
	PAID
	CLOSED
)

var orderStatus = []string{
	"DRAFT",
	"ENTERED",
	"PAID",
	"CLOSED",
}

func (o OrderStatusType) name() string {
	return orderStatus[o]
}

func (o OrderStatusType) ordinal() int {
	return int(o)
}

func (o OrderStatusType) String() string {
	return orderStatus[o]
}

func (o OrderStatusType) values() *[]string {
	return &orderStatus
}

// OrderStatusToInt function will return name
func OrderStatusToInt(value string) int {
	for i, b := range orderStatus {
		if b == value {
			return i
		}
	}
	return 0
}

func (status *OrderStatusType) UnmarshalJSON(b []byte) error {
	var value string
	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}
	*status = OrderStatusType(OrderStatusToInt(value))
	return nil
}

func (status OrderStatusType) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf(orderStatus[status]))
}

// Merchant is an exported POSTGRESL
type Merchant struct {
	ID              uint      `json:"-"`
	Name            string    `json:"name,omitempty" validate:"required"`
	Number          string    `json:"number,omitempty" validate:"required"`
	Email           string    `json:"email,omitempty"`
	NotificationURL string    `json:"notification_url,omitempty"`
	UUID            string    `json:"id" validate:"required"`
	CreatedAt       time.Time `json:"created_at,omitempty" validate:"required"`
	UpdatedAt       time.Time `json:"updated_at,omitempty" validate:"required"`
}

// Terminal is an exported
type Terminal struct {
	ID        uint      `json:"-"`
	UUID      string    `json:"id" validate:"required"`
	Number    string    `json:"-" validate:"required"`
	Merchant  Merchant  `json:"merchant,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Sub       *redis.PubSub
}

// Order is an exported
type Order struct {
	UUID        string          `json:"id"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty"`
	Number      string          `json:"number,omitempty"`
	Reference   string          `json:"reference,omitempty"`
	Status      OrderStatusType `json:"status"`
	Notes       string          `json:"notes,omitempty"`
	Price       int             `json:"price,omitempty"`
	MerchantID  string          `json:"merchant_id"`
	LogicNumber string          `json:"logic_number"`
	//Terminal     Terminal             `json:"terminal"`
	Items        []Item               `json:"items,omitempty"`
	Transactions []PaymentTransaction `json:"transactions,omitempty"`
}

// Item is an exported
type Item struct {
	UUID          string    `json:"id" validate:"required"`
	CreatedAt     time.Time `json:"-" validate:"required"`
	UpdatedAt     time.Time `json:"-" validate:"required"`
	Sku           string    `json:"sku,omitempty"`
	SkuType       string    `json:"sku_type,omitempty"`
	Name          string    `json:"name,omitempty" validate:"required"`
	Description   string    `json:"description,omitempty"`
	UnitPrice     string    `json:"unit_price,omitempty" validate:"required"`
	Quantity      int       `json:"quantity,omitempty"`
	UnitOfMeasure string    `json:"unit_of_measure,omitempty"`
	Details       string    `json:"details,omitempty"`
}

// PaymentTransaction is an exported
type PaymentTransaction struct {
	UUID              string         `json:"id"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	TransactionType   string         `json:"-"`
	ExternalID        string         `json:"external_id,omitempty"`
	Status            string         `json:"status,omitempty"`
	Description       string         `json:"description,omitempty"`
	TerminalNumber    string         `json:"terminal_number,omitempty"`
	Number            int            `json:"number,omitempty"`
	AuthorizationCode int            `json:"authorization_code,omitempty"`
	Amount            int            `json:"amount,omitempty"`
	Card              Card           `json:"card,omitempty"`
	PaymentProduct    PaymentProduct `json:"payment_product,omitempty"`
}

// Card is an exported
type Card struct {
	UUID      string    `json:"id"`
	Brand     string    `json:"brand,omitempty"`
	Mask      string    `json:"mask,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// PaymentProduct is an exported
type PaymentProduct struct {
	UUID           string         `json:"id"`
	Number         int            `json:"number,omitempty"`
	Name           string         `json:"name,omitempty"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	PaymentService PaymentService `json:"payment_service,omitempty"`
}

// PaymentService is an exported
type PaymentService struct {
	UUID      string    `json:"id"`
	Number    int       `json:"number,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
