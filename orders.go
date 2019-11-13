package sfox

import (
	"fmt"
	"net/url"

	"github.com/shopspring/decimal"
)

type AlgoID int

const (
	SmartAlgo     AlgoID = 200
	GorillaAlgo   AlgoID = 301
	TortoiseAlgo  AlgoID = 302
	HareAlgo      AlgoID = 303
	StopLossAlgo  AlgoID = 304
	PolarBearAlgo AlgoID = 305
	SniperAlgo    AlgoID = 306
	TWAPAlgo      AlgoID = 307
)

type AssetPair struct {
	FormattedSymbol string `json:"formatted_symbol"`
	Symbol          string `json:"symbol"`
}

type AssetPairs map[string]AssetPair

func (c *Client) AssetPairs() (AssetPairs, error) {
	var resp AssetPairs
	return resp, c.doGet("/v1/markets/currency-pairs", url.Values{}, &resp)
}

type NewOrder struct {
	Quantity      decimal.Decimal `json:"quantity"`
	Pair          string          `json:"currency_pair"`
	Price         decimal.Decimal `json:"price"`
	AlgoID        AlgoID          `json:"algo_id"`
	ClientOrderID string          `json:"client_order_id"`
	// seconds
	Interval int64 `json:"interval"`
	// seconds
	TotalTime int64 `json:"total_time"`
}

type Order struct {
	ID            int64           `json:"id"`
	Quantity      decimal.Decimal `json:"quantity"`
	Price         decimal.Decimal `json:"price"`
	Action        string          `json:"o_action"`
	Pair          string          `json:"pair"`
	Type          string          `json:"type"`
	VWAP          decimal.Decimal `json:"vwap"`
	Filled        decimal.Decimal `json:"filled"`
	Status        string          `json:"status"`
	StatusCode    int             `json:"status_code"`
	ClientOrderID string          `json:"client_order_id"`
	Updated       Time            `json:"dateupdated"`
	Expires       *Time           `json:""expires"`
	Fees          decimal.Decimal `json:"fees"`
	NetProceeds   decimal.Decimal `json:"net_proceeds"`
	Proceeds      decimal.Decimal `json:"proceeds"`
}

// Place an order.
// Example of an "Instant" buy: NewOrder{Quantity: decimal.NewFromFloat(0.001)}
// Smart buy (no limit): NewOrder{Quantity: decimal.NewFromFloat(0.04), AlgoID: SmartAlgo, Pair: "ethusd"}
func (c *Client) PlaceOrder(side Side, newOrder NewOrder) (Order, error) {
	order := map[string]interface{}{}

	if !newOrder.Quantity.IsZero() {
		order["quantity"] = newOrder.Quantity
	}
	if newOrder.Pair != "" {
		order["currency_pair"] = newOrder.Pair
	}
	if !newOrder.Price.IsZero() {
		order["price"] = newOrder.Price
	}
	if newOrder.AlgoID != 0 {
		order["algo_id"] = newOrder.AlgoID
	}
	if newOrder.ClientOrderID != "" {
		order["client_order_id"] = newOrder.ClientOrderID
	}
	if newOrder.AlgoID == TWAPAlgo {
		if newOrder.Interval != 0 {
			order["interval"] = newOrder.Interval
		}

		if newOrder.TotalTime != 0 {
			order["total_time"] = newOrder.TotalTime
		}
	}

	var resp Order
	return resp, c.doPost(fmt.Sprintf("/v1/orders/%s", side), order, &resp)
}

func (c *Client) CancelOrder(orderID int64) error {
	return c.doDelete(fmt.Sprintf("/v1/orders/%d", orderID))
}

func (c *Client) OpenOrders() ([]Order, error) {
	var res []Order
	return res, c.doGet("/v1/orders", url.Values{}, &res)
}

func (c *Client) GetOrder(orderID int64) (Order, error) {
	var res Order
	return res, c.doGet(fmt.Sprintf("/v1/orders/%d", orderID), url.Values{}, &res)
}
