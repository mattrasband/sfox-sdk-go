To Update in SFOX Docs:

* Add "held" to example response in User & Account Management > Get Account Balance
* Inconsistent timestamp formats: some are RFC3339/ISO8601 others are EPOCH
  {
    "id": 13328345,
    "order_id": "",
    "client_order_id": "",
    "day": "2019-11-07T00:00:00.000Z",
    "action": "Deposit",
    "currency": "usd",
    "memo": "",
    "amount": 25,
    "net_proceeds": 25,
    "price": 1,
    "fees": 0,
    "status": "done",
    "hold_expires": "2019-12-04 00:00:00",
    "tx_hash": "",
    "algo_name": "",
    "algo_id": "",
    "account_balance": 25.05191372,
    "AccountTransferFee": 0,
    "Description": "",
    "wallet_display_id": "5a3f1b1c-719d-11e9-b0be-0ea0e44d1000"
  }
* Interval is in seconds, not to be less than 15 minutes
* I can use the same client_id multiple times
