package sfox

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

const (
	ProductionWebsocketHost = "wss://ws.sfox.com/ws"
)

type Websocket struct {
	host string
}

type WebsocketEnvelope struct {
	Type       string          `json:"type"`
	Sequence   int64           `json:"sequence"`
	Timestamp  int64           `json:"timestamp"`
	Recipient  string          `json:"recipient"`
	RawPayload json.RawMessage `json:"payload"` // delay parsing
}

func (wse *WebsocketEnvelope) Payload() (interface{}, error) {
	const (
		tickerPrefix    = "ticker"
		orderbookPrefix = "orderbook"
		tradePrefix     = "trade"
	)

	var reified interface{}
	if strings.HasPrefix(wse.Recipient, tickerPrefix) {
		reified = &TickerMsg{}
	} else if strings.HasPrefix(wse.Recipient, orderbookPrefix) {
		reified = &OrderbookMsg{}
	} else if strings.HasPrefix(wse.Recipient, tradePrefix) {
		reified = &TradeMsg{}
	} else {
		return nil, ErrUnknownPayload
	}

	return reified, json.Unmarshal(wse.RawPayload, reified)
}

type TickerMsg struct {
	Amount   decimal.Decimal `json:"amount"`
	Exchange string          `json:"exchange"`
	Open     decimal.Decimal `json:"open"`
	High     decimal.Decimal `json:"high"`
	Low      decimal.Decimal `json:"low"`
	Last     decimal.Decimal `json:"last"`
	Volume   decimal.Decimal `json:"volume"`
	VWAP     decimal.Decimal `json:"vwap"`
	Pair     string          `json:"pair"`
	Route    string          `json:"route"`
	Source   string          `json:"source"`
	Time     Time            `json:"timestamp"`
}

type OrderbookMsg struct {
	Bids          []BidAsk          `json:"bids"`
	Asks          []BidAsk          `json:"asks"`
	Timestamps    map[string][]Time `json:"timestamps"`
	LastUpdated   Time              `json:"lastupdated"`
	LastPublished Time              `json:"lastpublished"`
	Pair          string            `json:"pair"`
	Currency      string            `json:"currency"`
}

type TradeMsg struct {
	ID          string          `json:"id"`
	Pair        string          `json:"pair"`
	Price       decimal.Decimal `json:"price,string"`
	Quantity    decimal.Decimal `json:"quantity,string"`
	Side        string          `json:"side"`
	BuyOrderID  string          `json:"buy_order_id"`
	SellOrderID string          `json:"sell_order_id"`
	Exchange    string          `json:"exchange"`
	ExchangeID  int             `json:"exchange_id"`
	Timestamp   Time            `json:"timestamp"`
}

type Event struct {
	Msg WebsocketEnvelope
	Err error
}

func (w *Websocket) Listen(ctx context.Context, feeds []string) (<-chan Event, error) {
	conn, _, err := websocket.DefaultDialer.Dial(w.host, nil)
	feed := make(chan Event, 100)
	if err != nil {
		conn.Close()
		close(feed)
		return feed, err
	}

	go func() {
		defer close(feed)
		defer conn.Close()

		err := conn.WriteJSON(struct {
			Type  string   `json:"type"`
			Feeds []string `json:"feeds"`
		}{
			Type:  "subscribe",
			Feeds: feeds,
		})
		if err != nil {
			feed <- Event{Err: err}
			return
		}

		for {
			evt := Event{}
			if err := conn.ReadJSON(&evt.Msg); err != nil {
				evt.Err = err
			}
			feed <- evt

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return feed, nil
}

func NewWebsocket() *Websocket {
	return NewWebsocketWithHost(ProductionWebsocketHost)
}

func NewWebsocketWithHost(host string) *Websocket {
	return &Websocket{
		host: host,
	}
}
