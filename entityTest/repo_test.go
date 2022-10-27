package entityTest

import (
	"entityRepo/entityRepo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInterfaceRepo(t *testing.T) {
	item := ItemRepo.New()
	item2 := ItemRepo.New()
	item.SetPrice(123.0123)
	item2.SetPrice(456.0123)
	item.SetBarcode(114514)
	item2.SetBarcode(1919810)
	item.SetBelongedItem(item2)
	item.AddContainedItem(item2)
	item2.AddContainedItem(item)
	item2.SetBelongedItem(item)

	err := entityRepo.Saver.Save()
	require.NoError(t, err)

	item3 := item.GetBelongedItem()
	require.NotNil(t, item3)
	itemArr := item.GetContainedItem()
	require.NotNil(t, itemArr)
}
