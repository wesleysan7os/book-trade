package entity

type TradeAction string
type TradeStatus string

const (
	BUY  TradeAction = "BUY"
	SELL TradeAction = "SELL"
)

const (
	OPEN   TradeStatus = "OPEN"
	CLOSED TradeStatus = "CLOSED"
)

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	OrderType     TradeAction
	Status        TradeStatus
	Transactions  []*Transaction
}

func NewOrder(orderID string, investor *Investor, asset *Asset, shares int, price float64, orderType TradeAction) *Order {
	return &Order{
		ID:            orderID,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        OPEN,
		Transactions:  []*Transaction{},
	}
}
