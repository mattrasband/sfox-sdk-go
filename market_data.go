package sfox

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/shopspring/decimal"
)

type Price struct {
	Fees     decimal.Decimal `json:"fees"`
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
	Total    decimal.Decimal `json:"total"`
	VWAP     decimal.Decimal `json:"vwap"`
}

func (c *Client) BestPrice(side Side, amount decimal.Decimal, pair string) (Price, error) {
	var res Price
	return res, c.doGet(fmt.Sprintf("/v1/offer/%s", side), url.Values{
		"amount": {amount.String()},
		"pair":   {pair},
	}, &res)
}

type BidAsk struct {
	Price    decimal.Decimal
	Size     decimal.Decimal
	Exchange string
}

func (ba *BidAsk) UnmarshalJSON(b []byte) error {
	var arr []interface{}
	if err := json.Unmarshal(b, &arr); err != nil {
		return err
	}

	ba.Price = decimal.NewFromFloat(arr[0].(float64))
	ba.Size = decimal.NewFromFloat(arr[1].(float64))
	ba.Exchange = arr[2].(string)
	return nil
}

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch val := v.(type) {
	case float64:
		*t = Time(time.Unix(0, int64(val)))
	case string:
		if val != "" {
			formats := []string{
				"2006-01-02T15:04:05Z",
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05.999",
			}

			var (
				err        error
				parsedTime time.Time
			)
			for _, format := range formats {
				parsedTime, err = time.Parse(format, val)
				if err == nil {
					break
				}
			}
			if err != nil {
				return err
			}
			*t = Time(parsedTime)
		}
	default:
		return fmt.Errorf("Error parsing time value type %T (%v)", val, val)
	}

	return nil
}

type Orderbook struct {
	Bids         []BidAsk
	Asks         []BidAsk
	MarketMaking []struct {
		Bids []BidAsk
		Asks []BidAsk
	}
	// TODO: we should clarify what the two entries in the time slice are
	Timestamps    map[string][]Time
	Pair          string
	Currency      string
	LastPublished Time
	LastUpdated   Time
}

func (c *Client) Orderbook(pair string) (Orderbook, error) {
	var res Orderbook
	return res, c.doGet(
		fmt.Sprintf("/v1/markets/orderbook/%s", pair),
		url.Values{},
		&res,
	)
}
