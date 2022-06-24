package order

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var orderList = []Order{
	{
		Region:        "Sub-Saharan Africa",
		Country:       "South Africa",
		ItemType:      "Fruits",
		SalesChannel:  "Offline",
		OrderPriority: "M",
		OrderDate:     time.Date(2012, time.July, 27, 0, 0, 0, 0, time.UTC),
		OrderID:       2345841,
		ShipDate:      time.Date(2012, time.July, 28, 0, 0, 0, 0, time.UTC),
		UnitsSold:     1593,
		UnitPrice:     9.33,
		UnitCost:      6.92,
		TotalRevenue:  14862.69,
		TotalCost:     11023.56,
		TotalProfit:   3839.13,
	},
	{
		Region:        "Sub-Saharan Africa",
		Country:       "South Africa",
		ItemType:      "Fruits",
		SalesChannel:  "Offline",
		OrderPriority: "L",
		OrderDate:     time.Date(2012, time.July, 27, 0, 0, 0, 0, time.UTC),
		OrderID:       2345841,
		ShipDate:      time.Date(2012, time.July, 28, 0, 0, 0, 0, time.UTC),
		UnitsSold:     1593,
		UnitPrice:     9.3339,
		UnitCost:      6.9229,
		TotalRevenue:  14862.6999,
		TotalCost:     11023.5699,
		TotalProfit:   3839.1339,
	},
}

type fakeStore struct {
	err error
}

func (s fakeStore) Save(Order) error {
	return s.err
}

type spyContext struct {
	order    Order
	code     int
	response map[string]string
	inputErr error
}

func (c *spyContext) Order() (Order, error) {
	return c.order, c.inputErr
}
func (c *spyContext) JSON(code int, v interface{}) {
	c.code = code
	c.response = v.(map[string]string)
}

func TestOrderNotAcceptOfflineChannel(t *testing.T) {
	h := &Handler{
		store:  fakeStore{},
		filter: "Online",
	}

	c := spyContext{order: orderList[0]}
	h.Order(&c)

	want := "Offline is not accept"

	if want != c.response["message"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["message"])
	}
}

func TestOrderInputError(t *testing.T) {
	h := &Handler{
		store:  fakeStore{},
		filter: "Online",
	}

	c := spyContext{inputErr: errors.New("input error")}
	h.Order(&c)

	want := "input error"

	if want != c.response["error"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["error"])
	}
}

func TestOrderDBError(t *testing.T) {
	h := &Handler{
		store:  fakeStore{err: errors.New("db error")},
		filter: "Offline",
	}

	c := spyContext{order: orderList[0]}
	h.Order(&c)

	want := "db error"

	if want != c.response["error"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["error"])
	}
}

func TestRandomOrderID(t *testing.T) {
	h := &Handler{
		store:  fakeStore{},
		filter: "Offline",
	}

	c := spyContext{order: orderList[0]}
	for i := 1; i <= 10; i++ {
		c.order.OrderID = uint(i)
		h.Order(&c)

		want := fmt.Sprintf("%d is saved", c.order.OrderID)

		if want != c.response["message"] {
			t.Errorf("%q is expected but got %q\n", want, c.response["message"])
		}
		t.Logf("round %d is pass", i)
	}
}

func TestRandomData(t *testing.T) {
	h := &Handler{
		store:  fakeStore{},
		filter: "Offline",
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var orderList []Order
	for i := 0; i < 10; i++ {
		data := Order{
			Region:        "Sub-Saharan Africa",
			Country:       "South Africa",
			ItemType:      "Fruits",
			SalesChannel:  "Offline",
			OrderPriority: "M",
			OrderDate:     time.Date(2012, time.July, 27, 0, 0, 0, 0, time.UTC),
			OrderID:       uint(r1.Uint64()),
			ShipDate:      time.Date(2012, time.July, 28, 0, 0, 0, 0, time.UTC),
			UnitsSold:     uint(r1.Uint64()),
			UnitPrice:     r1.Float64(),
			UnitCost:      r1.Float64(),
			TotalRevenue:  r1.Float64(),
			TotalCost:     r1.Float64(),
			TotalProfit:   r1.Float64(),
		}
		orderList = append(orderList, data)
	}

	for i, v := range orderList {
		c := spyContext{order: v}
		h.Order(&c)

		want := fmt.Sprintf("%d is saved", c.order.OrderID)

		if want != c.response["message"] {
			t.Errorf("%q is expected but got %q\n", want, c.response["message"])
		}
		t.Logf("round: %d, orderID: %d, unitSold: %d, is pass", i+1, c.order.OrderID, c.order.UnitsSold)
	}
}

func TestOrderSalesChannelIsEmpty(t *testing.T) {
	h := &Handler{
		store:  fakeStore{},
		filter: "Online",
	}

	c := spyContext{order: orderList[0]}
	c.order.SalesChannel = ""
	h.Order(&c)

	want := "sales_channel is not empty"

	if want != c.response["error"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["error"])
	}
}

func TestSliceDataOfOrder(t *testing.T) {
	h := &Handler{
		store:  fakeStore{},
		filter: "Offline",
	}

	for i, v := range orderList {
		c := spyContext{order: v}
		h.Order(&c)

		want := fmt.Sprintf("%d is saved", c.order.OrderID)

		if want != c.response["message"] {
			t.Errorf("%q is expected but got %q\n", want, c.response["message"])
		}
		t.Logf("round: %d, orderID: %d, unitSold: %d, is pass", i+1, c.order.OrderID, c.order.UnitsSold)
	}
}
