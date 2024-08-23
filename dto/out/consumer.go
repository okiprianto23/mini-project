package out

type ConsumerOut struct {
	ID          int64   `json:"ID"`
	UUID        string  `json:"UUID"`
	NIK         string  `json:"NIK"`
	FullName    string  `json:"fullName"`
	LegalName   string  `json:"legalName"`
	BirthPlace  string  `json:"birthPlace"`
	BirthDate   string  `json:"birthDate"`
	Salary      float64 `json:"salary"`
	KTPPhoto    string  `json:"KTPPhoto"`
	SelfiePhoto string  `json:"selfiePhoto"`
}

type CreditLimitOut struct {
	ID                  int64   `json:"id"`
	UUID                string  `json:"UUID"`
	ConsumerID          int64   `json:"consumerID"`
	MonthlyInstallments float64 `json:"monthlyInstallments"`
	InterestRate        float64 `json:"interestRate"`
	Tenor               int64   `json:"tenor"`
	LimitAmount         float64 `json:"limitAmount"`
}
