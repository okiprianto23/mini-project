package consumer

import (
	"database/sql"
	"main-xyz/config"
	"main-xyz/context"
	"main-xyz/dao"
	"main-xyz/dto/in"
	"main-xyz/error/bundles"
	"main-xyz/repository"
	"main-xyz/tx_helper"
	"math"
	"time"
)

func NewConsumerService(
	logger *config.LoggerCustom,
	consumerDAO dao.ConsumerDAO,
	userDAO dao.UserDAO,
	creditLimitDAO dao.CreditLimitDAO,
	bundles bundles.Bundles,
	txHelper tx_helper.TXHelper,
) ConsumerService {
	service := consumerService{
		logger:         logger,
		consumerDAO:    consumerDAO,
		userDAO:        userDAO,
		creditLimitDAO: creditLimitDAO,
		bundles:        bundles,
		txHelper:       txHelper,
	}

	return &service
}

type consumerService struct {
	logger         *config.LoggerCustom
	consumerDAO    dao.ConsumerDAO
	userDAO        dao.UserDAO
	creditLimitDAO dao.CreditLimitDAO

	bundles  bundles.Bundles
	txHelper tx_helper.TXHelper
}

func (c *consumerService) GetDTO() interface{} {
	return &in.ConsumerRequest{}
}

func (c *consumerService) GetMultipartDTO() interface{} {
	return &in.ConsumerMultipart{}
}

func (c *consumerService) parseDTO(cr interface{}) *in.ConsumerRequest {
	return cr.(*in.ConsumerRequest)
}

func (c *consumerService) parseMultipartDTO(cr interface{}) *in.ConsumerMultipart {
	return cr.(*in.ConsumerMultipart)
}

func (c consumerService) consumerConvertToRepo(ctx *context.ContextModel, inputStruct in.ConsumerRequest, now time.Time) repository.ConsumerModel {
	return repository.ConsumerModel{
		UserID:      sql.NullInt64{Int64: int64(inputStruct.UserID)},
		NIK:         sql.NullString{String: inputStruct.NIK},
		FullName:    sql.NullString{String: inputStruct.FullName},
		LegalName:   sql.NullString{String: inputStruct.LegalName},
		BirthPlace:  sql.NullString{String: inputStruct.BirthPlace},
		BirthDate:   sql.NullTime{Time: inputStruct.BirthDate},
		Salary:      sql.NullFloat64{Float64: inputStruct.Salary},
		KTPPhoto:    sql.NullString{String: inputStruct.KTPPhoto},
		SelfiePhoto: sql.NullString{String: inputStruct.SelfiePhoto},
		DefaultCreatedUpdated: repository.DefaultCreatedUpdated{
			UpdatedBy:     sql.NullInt64{Int64: ctx.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:     sql.NullTime{Time: now},
			UpdatedClient: sql.NullString{String: ctx.AuthAccessTokenModel.ClientID},
			CreatedBy:     sql.NullInt64{Int64: ctx.AuthAccessTokenModel.ResourceUserID},
			CreatedAt:     sql.NullTime{Time: now},
			CreatedClient: sql.NullString{String: ctx.AuthAccessTokenModel.ClientID},
		},
	}
}

/*
logicForCreditLimit

	Untuk perhitungan berdasarkan mendapatkan sebuah cicilan berdasarkan salary
	C=(P⋅r⋅(1+r)n)/((1+r)n-1)
	Dimana:

	𝐶 = cicilan bulanan (Rp 2.000.000)
	𝑟 = suku bunga bulanan (0,01)
	𝑛 = tenor dalam bulan
	𝑃 = jumlah pinjaman (limit yang akan dihitung)
*/
func (c consumerService) logicForCreditLimit(ctx *context.ContextModel, tx *sql.Tx, consumerModel repository.ConsumerModel, now time.Time) error {
	monthlyInstallments := 0.4 * consumerModel.Salary.Float64 // 40% dari gaji bulanan
	interestRate := 12.0                                      // 12% per tahun

	ccL := newCalculateCreditLimit(consumerModel.ID.Int64)
	ccL.setMonthlyInstallments(monthlyInstallments)
	ccL.setInterestRate(interestRate)
	ccL.setContextModel(ctx)
	ccL.setTimeNow(now)

	// setelah dapat limitnya masukan kedalam db
	var creditModels []repository.CreditLimitModel
	creditModels = append(creditModels,
		ccL.setTenor(1).calculateLoanLimit(),
		ccL.setTenor(2).calculateLoanLimit(),
		ccL.setTenor(3).calculateLoanLimit(),
		ccL.setTenor(6).calculateLoanLimit(),
	)

	for _, cm := range creditModels {
		_, err := c.creditLimitDAO.InsertCreditLimit(ctx, tx, cm)
		if err != nil {
			return err
		}
	}

	return nil
}

// newCalculateCreditLimit untuk membuat initiate function struct
func newCalculateCreditLimit(consumerID int64) *calculateCreditLimit {
	return &calculateCreditLimit{
		consumerID: consumerID,
	}
}

type calculateCreditLimit struct {
	ctx                 *context.ContextModel
	now                 time.Time
	consumerID          int64
	monthlyInstallments float64
	interestRate        float64
	tenor               int64
	limitAmount         float64
}

func (cc *calculateCreditLimit) setContextModel(ctx *context.ContextModel) *calculateCreditLimit {
	cc.ctx = ctx
	return cc
}

func (cc *calculateCreditLimit) setTimeNow(now time.Time) *calculateCreditLimit {
	cc.now = now
	return cc
}

func (cc *calculateCreditLimit) setMonthlyInstallments(mi float64) *calculateCreditLimit {
	cc.monthlyInstallments = mi
	return cc
}

func (cc *calculateCreditLimit) setInterestRate(ir float64) *calculateCreditLimit {
	cc.interestRate = ir
	return cc
}

func (cc *calculateCreditLimit) setTenor(t int64) *calculateCreditLimit {
	cc.tenor = t
	return cc
}

// calculateLoanLimit Fungsi untuk menghitung limit pinjaman
func (cc *calculateCreditLimit) calculateLoanLimit() repository.CreditLimitModel {
	// Konversi suku bunga tahunan ke suku bunga bulanan
	r := cc.interestRate / 100 / 12
	// Hitung limit pinjaman menggunakan rumus
	P := cc.monthlyInstallments * (math.Pow(1+r, float64(cc.tenor)) - 1) / (r * math.Pow(1+r, float64(cc.tenor)))
	cc.limitAmount = P
	return repository.CreditLimitModel{
		ConsumerID:          sql.NullInt64{Int64: cc.consumerID},
		MonthlyInstallments: sql.NullFloat64{Float64: cc.monthlyInstallments},
		InterestRate:        sql.NullFloat64{Float64: cc.interestRate},
		Tenor:               sql.NullInt64{Int64: cc.tenor},
		Limit:               sql.NullFloat64{Float64: cc.limitAmount},
		RemainingLimit:      sql.NullFloat64{Float64: cc.limitAmount},
		DefaultCreatedUpdated: repository.DefaultCreatedUpdated{
			UpdatedBy:     sql.NullInt64{Int64: cc.ctx.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:     sql.NullTime{Time: cc.now},
			UpdatedClient: sql.NullString{String: cc.ctx.AuthAccessTokenModel.ClientID},
			CreatedBy:     sql.NullInt64{Int64: cc.ctx.AuthAccessTokenModel.ResourceUserID},
			CreatedAt:     sql.NullTime{Time: cc.now},
			CreatedClient: sql.NullString{String: cc.ctx.AuthAccessTokenModel.ClientID},
		},
	}
}
