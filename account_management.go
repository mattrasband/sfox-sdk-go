package sfox

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	Currency  string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance"`
	Available decimal.Decimal `json:"available"`
	Held      decimal.Decimal `json"held"`
}

func (c *Client) AccountBalance() ([]Balance, error) {
	var res []Balance
	return res, c.doGet("/v1/user/balance", url.Values{}, &res)
}

type TimeBasedPagination struct {
	From time.Time
	To   time.Time
}

type Transaction struct {
	ID              int64           `json:"id"`
	OrderID         string          `json:"order_id"`
	ClientOrderID   string          `json:"client_order_id"`
	Day             Time            `json:"day"`
	Action          string          `json:"action"` // TODO: make a type for these
	Currency        string          `json:"currency"`
	Amount          decimal.Decimal `json:"amount"`
	NetProceeds     decimal.Decimal `json:"net_proceeds"`
	Price           decimal.Decimal `json:"price"`
	Fees            decimal.Decimal `json:"fees"`
	Status          string          `json:"status"`
	HoldExpires     Time            `json:"hold_expires"`
	TxHash          string          `json:"tx_hash"`
	AlgoName        string          `json:"algo_name"`
	AlgoID          string          `json:"algo_id"`
	AccountBalance  decimal.Decimal `json:"account_balance"`
	WalletDisplayID string          `json:"wallet_display_id"`
}

// Transactions retrieves all the user's transactions, by default this will
// be all transactions inception to date.
func (c *Client) Transactions(t ...TimeBasedPagination) ([]Transaction, error) {
	params := url.Values{}
	if len(t) > 0 {
		page := t[0]
		params["from"] = []string{
			strconv.FormatInt(page.From.UnixNano()/1000, 10),
		}
		params["to"] = []string{
			strconv.FormatInt(page.To.UnixNano()/1000, 10),
		}
	}

	var res []Transaction
	return res, c.doGet("/v1/account/transactions", params, &res)
}

type ACHResponse struct {
	TXStatus int    `json:"tx_status"`
	Success  bool   `json:"success"`
	Error    string `json:"error"`
}

func (c *Client) ACHDeposit(amount decimal.Decimal) (ACHResponse, error) {
	var deposit struct {
		Amount decimal.Decimal `json:"amount"`
	}
	deposit.Amount = amount

	var res ACHResponse
	return res, c.doPost("/v1/user/bank/deposit", deposit, &res)
}

type Address struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
}

func (c *Client) CryptoAddresses(currency string) ([]Address, error) {
	var result map[string][]Address
	err := c.doGet(fmt.Sprintf("/v1/user/deposit/address/%s", currency), url.Values{}, &result)
	if err != nil {
		return []Address{}, err
	}
	return result[currency], nil
}

func (c *Client) CreateAddress(currency string) (Address, error) {
	var result Address
	return result, c.doPost(fmt.Sprintf("/v1/user/deposit/address/%s", currency), nil, &result)
}

type WithdrawResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// TODO: untested (getting 504s)
func (c *Client) Withdraw(amount decimal.Decimal, currency, address string) (WithdrawResponse, error) {
	var request struct {
		Amount   decimal.Decimal `json:"amount"`
		Address  string          `json:"address"`
		Currency string          `json:"currency"`
	}
	request.Amount = amount
	request.Address = address
	request.Currency = currency

	var result WithdrawResponse
	return result, c.doPost("/v1/user/withdraw", request, &result)
}
