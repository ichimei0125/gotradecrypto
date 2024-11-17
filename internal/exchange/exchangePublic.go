package exchange

import "fmt"

type Exchange interface {
	FetchPrice(symbol Symbol) (float64, error)
}

// 主逻辑，接受 Exchange 接口
func FetchPriceFromExchange(exchange Exchange, symbol Symbol) {
	price, err := exchange.FetchPrice(symbol)
	if err != nil {
		fmt.Printf("Error fetching price: %v\n", err)
		return
	}
	fmt.Printf("Price: %.2f\n", price)
}
