package entityTest

import (
	"entityRepo/entityRepo"
	"log"
)

const (
	CardPaymentInheritType entityRepo.GoenInheritType = iota + 1
)

func init() {

	tmpItemRepo, err := entityRepo.NewRepo[ItemEntity, Item]("item")
	if err != nil {
		log.Fatal(err)
	}
	itemRepo = tmpItemRepo
	ItemRepo = tmpItemRepo

	tmpPaymentRepo, err := entityRepo.NewRepo[PaymentEntity, Payment]("payment")
	if err != nil {
		log.Fatal(err)
	}
	paymentRepo = tmpPaymentRepo
	PaymentRepo = tmpPaymentRepo

	tmpCardPaymentRepo, err := entityRepo.NewInheritRepo[CardPaymentEntity, CardPayment]("card_payment", tmpPaymentRepo, CardPaymentInheritType)
	if err != nil {
		log.Fatal(err)
	}
	cardPaymentRepo = tmpCardPaymentRepo
	CardPaymentRepo = tmpCardPaymentRepo

}
