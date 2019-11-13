# sfox-sdk-go

An unofficial client to trade on [SFOX](https://sfox.com).

# Usage

```golang
package main

import (
    "fmt"

    "github.com/mrasband/sfox-sdk-go"
    "github.com/shopspring/decimal"
)

func main() {
    c := sfox.New(os.Getenv("SFOX_API_KEY"))

    fmt.Println("Placing order...")
    order, err := c.PlaceOrder(sfox.Buy, sfox.NewOrder{
        Quantity: decimal.NewFromFloat(0.04),
        AlgoID:   sfox.SmartAlgo,
        Pair:     "ethusd",
        Price:    decimal.NewFromFloat(125),
    })
    if err != nil {
        fmt.Printf("Err: %s\n", err)
        return
    }
    fmt.Printf("Order: %+v\n", order)

    fmt.Println("Getting open orders...")
    orders, err := c.OpenOrders()
    if err != nil {
        fmt.Printf("Err: %s\n", err)
        return
    }
    fmt.Printf("Orders: %+v\n", orders)

    fmt.Println("Getting order...")
    order2, err := c.GetOrder(order.ID)
    if err != nil {
        fmt.Printf("Err: %s\n", err)
        return
    }
    fmt.Printf("Order: %+v\n", order2)

    time.Sleep(10 * time.Second)

    fmt.Println("Cancelling order...")
    err = c.CancelOrder(order.ID)
    if err != nil {
        fmt.Printf("Err: %s\n", err)
        return
    }
    fmt.Println("OK")

    fmt.Println("Getting cancelled order...")
    order2, err = c.GetOrder(order.ID)
    if err != nil {
        fmt.Printf("Err: %s\n", err)
        return
    }
    fmt.Printf("Order: %+v\n", order2)
}
```
