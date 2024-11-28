package models

import "time"

type CustomerRevenue struct {
	CustomerID int64
	Email      string
	Revenue    float64
}

type QuantileStats struct {
	NumberOfCustomers int
	MaxRevenue        float64
	MinRevenue        float64
}

type CustomerEventData struct {
	EventID     int64
	CustomerID  int64
	EventTypeID int16
	ContentID   int
	Quantity    int16
	EventDate   time.Time
}

type ContentPrice struct {
	ContentID int
	Price     float64
}

type CustomerEmail struct {
	CustomerID   int64
	ChannelValue string
}

type ExportRecord struct {
	CustomerID int64
	Email      string
	Revenue    float64
}
