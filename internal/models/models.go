package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Order struct {
	OrderUID          string `json:"order_uid"          db:"order_uid"          validate:"required"`
	TrackNumber       string `json:"track_number"       db:"track_number"       validate:"required"`
	Entry             string `json:"entry"              db:"entry"              validate:"required"`
	Delivery          `json:"delivery"  validate:"required"`
	Payment           `json:"payment"  validate:"required"`
	Items             []OrderItem `json:"items"              db:"order_items"        validate:"required"`
	Locale            string      `json:"locale"             db:"locale"`
	InternalSignature string      `json:"internal_signature" db:"internal_signature"`
	CustomerID        string      `json:"customer_id"        db:"customer_id"        validate:"required"`
	DeliveryService   string      `json:"delivery_service"   db:"delivery_service"   validate:"required"`
	Shardkey          string      `json:"shardkey"           db:"shardkey"`
	SMID              int         `json:"sm_id"              db:"sm_id"`
	DateCreated       time.Time   `json:"date_created"       db:"date_created"       validate:"required"`
	OOFShard          string      `json:"oof_shard"          db:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"    db:"name"    validate:"required"`
	Phone   string `json:"phone"   db:"phone"   validate:"required,e164"`
	Zip     string `json:"zip"     db:"zip"     validate:"required"`
	City    string `json:"city"    db:"city"    validate:"required"`
	Address string `json:"address" db:"address" validate:"required"`
	Region  string `json:"region"  db:"region"  validate:"required"`
	Email   string `json:"email"   db:"email"   validate:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"   db:"payment_transaction"   validate:"required"`
	RequestID    string `json:"request_id"    db:"payment_request_id"`
	Currency     string `json:"currency"      db:"payment_currency"      validate:"required"`
	Provider     string `json:"provider"      db:"payment_provider"      validate:"required"`
	Amount       int    `json:"amount"        db:"payment_amount"        validate:"required"`
	PaymentDT    int64  `json:"payment_dt"    db:"payment_dt"            validate:"required"`
	Bank         string `json:"bank"          db:"payment_bank"          validate:"required"`
	DeliveryCost int    `json:"delivery_cost" db:"payment_delivery_cost" validate:"required"`
	GoodsTotal   int    `json:"goods_total"   db:"payment_goods_total"   validate:"required"`
	CustomFee    int    `json:"custom_fee"    db:"payment_custom_fee"`
}

type OrderItem struct {
	ChrtID      int    `json:"chrt_id"      db:"chrt_id"      validate:"required"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required"`
	Price       int    `json:"price"        db:"price"        validate:"required"`
	RID         string `json:"rid"          db:"rid"          validate:"required"`
	Name        string `json:"name"         db:"name"         validate:"required"`
	Sale        int    `json:"sale"         db:"sale"         validate:"required"`
	Size        string `json:"size"         db:"size"`
	TotalPrice  int    `json:"total_price"  db:"total_price"  validate:"required"`
	NMID        int    `json:"nm_id"        db:"nm_id"        validate:"required"`
	Brand       string `json:"brand"        db:"brand"        validate:"required"`
	Status      int    `json:"status"       db:"status"       validate:"required"`
}

type OrderData struct {
	OrderUID string `db:"order_uid"`
	Data     []byte `db:"data"`
}

func ValidationErrors(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid email", err.Field()))
		case "e164":
			errMsgs = append(
				errMsgs,
				fmt.Sprintf("field %s is not valid phone number", err.Field()),
			)
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return strings.Join(errMsgs, "\n")
}
