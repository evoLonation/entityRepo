package entityTest

import (
	"entityRepo/entityRepo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInherit(t *testing.T) {
	//cardPayment, err := CardPaymentRepo.Get(1)
	//payment, err := PaymentRepo.Get(1)
	//log.Fatal(err)
	//print(cardPayment)
	//p := CardPaymentRepo.New()
	//
	//p.GoenId = 1
	//cp := CardPaymentRepo.New()
	//cp.SetAmountTendered(0.111222)
	//cp.SetCardAccountNumber(123124)
	//CardPaymentRepo.AddInAllInstance(cp)
	//err := entityRepo.Saver.Save()
	//require.NoError(t, err)

	//cp2, err := CardPaymentRepo.GetFromAllInstanceBy("amount_tendered", 0.111222)
	cp2 := CardPaymentRepo.GetFromAllInstanceBy("card_account_number", 123124)
	cp2.SetAmountTendered(0.333444)
	cp2.SetCardAccountNumber(3215125)
	err := entityRepo.Saver.Save()
	require.NoError(t, err)
	//CardPaymentRepo.GetFromAllInstanceBy("expiry_date", 123)
	//CardPaymentRepo.GetFromAllInstanceBy("amount_tendered", 123)
	//CardPaymentRepo.RemoveFromAllInstance(p)
	//CardPaymentRepo.FindFromAllInstanceBy("expiry_date", 123)
	//CardPaymentRepo.FindFromAllInstanceBy("amount_tendered", 123)

}
func TestInherit2(t *testing.T) {
	p := PaymentRepo.GetFromAllInstanceBy("goen_id", 3)
	prr := PaymentRepo.FindFromAllInstanceBy("goen_id", 3)
	require.NotNil(t, prr)
	var cp CardPayment
	if PaymentRepo.GetRealType(p) == CardPaymentInheritType {
		var err error
		cp, err = CardPaymentRepo.CastFrom(p)
		require.NoError(t, err)
	}
	require.NotNil(t, cp)
}

func TestInherit3(t *testing.T) {
	item := ItemRepo.New()
	cp := CardPaymentRepo.New()
	cp.SetCardAccountNumber(1838940019)
	item.SetPrice(123.0123)
	item.SetBarcode(114514)
	item.SetBelongedPayment(cp)

	err := entityRepo.Saver.Save()
	require.NoError(t, err)

	p := item.GetBelongedPayment()
	require.NoError(t, err)
	require.Equal(t, PaymentRepo.GetRealType(p), CardPaymentInheritType)
	cp2, err := CardPaymentRepo.CastFrom(p)
	require.NoError(t, err)
	require.NotNil(t, cp2)
}

func TestPayment(t *testing.T) {
	p := PaymentRepo.New()
	cp := CardPaymentRepo.New()
	cp.SetCardAccountNumber(1838940019)
	p.SetAmountTendered(123)
	cp.SetAmountTendered(456)
	cp.SetExpiryDate(time.Now())
	PaymentRepo.AddInAllInstance(p)
	CardPaymentRepo.AddInAllInstance(cp)

	err := entityRepo.Saver.Save()
	require.NoError(t, err)

	p2 := PaymentRepo.GetFromAllInstanceBy("amount_tendered", 123)
	require.NoError(t, err)
	require.NotNil(t, p2)
	cp2 := CardPaymentRepo.GetFromAllInstanceBy("amount_tendered", 456)
	require.NoError(t, err)
	require.NotNil(t, cp2)
	p3 := PaymentRepo.GetFromAllInstanceBy("amount_tendered", 456)
	require.NoError(t, err)
	require.NotNil(t, p3)

	require.Equal(t, PaymentRepo.GetRealType(p3), CardPaymentInheritType)
	cp3, err := CardPaymentRepo.CastFrom(p3)
	require.NoError(t, err)
	require.NotNil(t, cp3)

}
