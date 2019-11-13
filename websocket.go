package sfox

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

const (
	WebsocketHost = "wss://ws.sfox.com/ws"
	subscribe     = "subscribe"
	unsubscribe   = "unsubscribe"
)

var (
	wsTypeRe = regexp.MustCompile(`recipient":\s*"([^.]+)`)
)

type Websocket struct {
	host string
}

type subscribeMsg struct {
	Type  string   `json:"type"`
	Feeds []string `json:"feeds"`
}

type WebsocketEnvelope struct {
	Type      string `json:"type"`
	Sequence  int64  `json:"sequence"`
	Timestamp int64  `json:"timestamp"`
	Recipient string `json:"recipient"`
}

type TickerPayload struct {
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

type WSTicker struct {
	WebsocketEnvelope
	Payload TickerPayload `json:"payload"`
}

type OrderbookPayload struct {
	Bids          []BidAsk          `json:"bids"`
	Asks          []BidAsk          `json:"asks"`
	Timestamps    map[string][]Time `json:"timestamps"`
	LastUpdated   Time              `json:"lastupdated"`
	LastPublished Time              `json:"lastpublished"`
	Pair          string            `json:"pair"`
	Currency      string            `json:"currency"`
}

type WSOrderbook struct {
	WebsocketEnvelope
	Payload OrderbookPayload `json:"payload"`
}

type TradePayload struct {
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

type WSTrade struct {
	WebsocketEnvelope
	Payload TradePayload `json:"payload"`
}

type Event struct {
	Err error
	Msg interface{}
}

func (w *Websocket) Listen(ctx context.Context, feeds []string) (<-chan Event, error) {
	conn, _, err := websocket.DefaultDialer.Dial(w.host, nil)
	feed := make(chan Event, 100)
	if err != nil {
		close(feed)
		return feed, err
	}

	go func() {
		defer close(feed)
		defer conn.Close()

		err := conn.WriteJSON(subscribeMsg{Type: subscribe, Feeds: feeds})
		if err != nil {
			feed <- Event{Err: err}
			return
		}

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				feed <- Event{Err: err}
			} else {
				match := wsTypeRe.FindSubmatch(data)
				if len(match) > 0 {
					evt := Event{}

					switch string(match[1]) {
					case "ticker":
						var msg WSTicker
						if err := json.Unmarshal(data, &msg); err != nil {
							evt.Err = err
						}
						evt.Msg = msg

					case "trades":
						var msg WSTrade
						if err := json.Unmarshal(data, &msg); err != nil {
							evt.Err = err
						}
						evt.Msg = msg

					case "orderbook":
						var msg WSOrderbook
						if err := json.Unmarshal(data, &msg); err != nil {
							evt.Err = err
						}
						evt.Msg = msg
					}

					feed <- evt
				}
			}

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
	return &Websocket{
		host: WebsocketHost,
	}
}
