package entityTest

import (
	"entityRepo/entityRepo"
	"time"
)

var cardPaymentRepo entityRepo.RepoForEntity[CardPayment]
var CardPaymentRepo entityRepo.InheritRepoForOther[CardPayment]

type CardPayment interface {
	Payment
	SetCardAccountNumber(int)
	SetExpiryDate(time.Time)
	GetCardAccountNumber() int
	GetExpiryDate() time.Time
}

type CardPaymentEntity struct {
	PaymentEntity
	entityRepo.FieldChange
	ExpiryDate        time.Time `db:"expiry_date"`
	CardAccountNumber int       `db:"card_account_number"`
}

func (p *CardPaymentEntity) GetParentEntity() entityRepo.EntityForInheritRepo {
	return &p.PaymentEntity
}

func (p *CardPaymentEntity) SetExpiryDate(t time.Time) {
	p.ExpiryDate = t
	p.AddBasicFieldChange("expiry_date")
}

func (p *CardPaymentEntity) SetCardAccountNumber(cardAccountNumber int) {
	p.CardAccountNumber = cardAccountNumber
	p.AddBasicFieldChange("card_account_number")
}

func (p *CardPaymentEntity) GetCardAccountNumber() int {
	return p.CardAccountNumber
}

func (p *CardPaymentEntity) GetExpiryDate() time.Time {
	return p.ExpiryDate
}
