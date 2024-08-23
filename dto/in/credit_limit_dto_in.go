package in

import "time"

type CreditLimitRequest struct {
	ID           int64     `json:"ID" empty:"allowed"`
	ConsumerID   int64     `json:"consumerID" required:"insert,update"`
	Tenor        int16     `json:"tenor" required:"insert,update" enum:"tenor"`
	LimitAmount  float64   `json:"limitAmount" required:"insert,update"`
	UpdatedBy    int64     `json:"updated_by"`
	UpdatedAtStr string    `json:"updated_at" required:"update,delete"`
	UpdatedAt    time.Time `required:"update,delete"`
}
