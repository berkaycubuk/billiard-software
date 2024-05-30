package order

import (
	"testing"

	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/shopspring/decimal"
)

func TestOrderCreate(t *testing.T) {
	order, err := Create(nil, decimal.Zero, models.ORDER_STATUS_WAITING)
	if err != nil {
		t.Error(err.Error())
	}

	if order.Status != models.ORDER_STATUS_WAITING {
		t.Errorf(`Expected order status "WAITING" but got %o`, order.Status)
	}

	if !order.Price.Equal(decimal.Zero) {
		t.Errorf(`Expected order price "0" but got %o`, order.Price)
	}
}
