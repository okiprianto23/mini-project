package in

import (
	"time"
)

type TransactionRequest struct {
	ID                int64     `json:"id" min:"1" empty:"allowed"`
	ConsumerID        int64     `json:"consumerID"`
	ContractNumber    string    `json:"contract_number" empty:"allowed"`
	Otr               float64   `json:"otr" required:"insert" min:"1"`
	AdminFee          float64   `json:"adminFee" required:"insert" min:"1"`
	InstallmentAmount float64   `json:"installmentAmount" empty:"allowed"`
	InterestAmount    float64   `json:"interestAmount" empty:"allowed"`
	AssetName         string    `json:"assetName" required:"insert"`
	TransactionDate   time.Time `json:"transaction_date"  empty:"allowed"`
	UpdatedBy         int64     `json:"updated_by"`
	UpdatedAtStr      string    `json:"updated_at" required:"update,delete"`
	UpdatedAt         time.Time `required:"update,delete"`
}
