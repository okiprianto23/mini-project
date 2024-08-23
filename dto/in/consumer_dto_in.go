package in

import "time"

type ConsumerMultipart struct {
	Consumer    ConsumerRequest `json:"consumer" multipart:"json" required:"insert,update"`
	KTPPhoto    MultipartFile   `json:"KTPPhoto" multipart:"file" required:"insert,update" ext:"jpg,jpeg,png" min:"" max:"3000000"`
	SelfiePhoto MultipartFile   `json:"selfiePhoto" multipart:"file" required:"insert,update" ext:"jpg,jpeg,png" min:"" max:"3000000"`
}

type ConsumerRequest struct {
	ID           int64     `json:"id" min:"1" empty:"allowed"`
	UserID       int64     `json:"userID"`
	NIK          string    `json:"NIK" regex:"nik" required:"insert,update"`
	FullName     string    `json:"fullName" required:"insert,update" min:"1" max:"100"`
	LegalName    string    `json:"legalName" required:"insert,update" min:"1" max:"100"`
	BirthPlace   string    `json:"birthPlace" required:"insert,update"`
	BirthDateStr string    `json:"birthDate" required:"insert,update"`
	BirthDate    time.Time `required:"insert,update" dateFormat:"date_only"`
	Salary       float64   `json:"salary" required:"insert,update"`
	KTPPhoto     string    `json:"KTPPhoto"  required:"insert,update" empty:"allowed"`
	SelfiePhoto  string    `json:"selfiePhoto"  required:"insert,update" empty:"allowed"`
	UpdatedBy    int64     `json:"updated_by"`
	UpdatedAtStr string    `json:"updated_at" required:"update,delete"`
	UpdatedAt    time.Time `required:"update,delete"`
}
