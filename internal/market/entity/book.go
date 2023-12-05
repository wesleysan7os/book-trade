package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrdersChan    chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChan {
		switch order.OrderType {
		case BUY:
			buyOrders.Push(order)
			b.processTransactions(sellOrders, order)
		case SELL:
			sellOrders.Push(order)
			b.processTransactions(buyOrders, order)
		}
	}
}

func (b *Book) processTransactions(oppositeOrders *OrderQueue, order *Order) {
	if oppositeOrders.Len() > 0 && oppositeOrders.Orders[0].Price <= order.Price {
		oppositeOrder := oppositeOrders.Pop().(*Order)
		if oppositeOrder.PendingShares > 0 {
			transaction := NewTransaction(oppositeOrder, order, order.Shares, oppositeOrder.Price)
			b.AddTransaction(transaction, b.Wg)
			oppositeOrder.Transactions = append(oppositeOrder.Transactions, transaction)
			order.Transactions = append(order.Transactions, transaction)
			b.OrdersChanOut <- oppositeOrder
			b.OrdersChanOut <- order
			if oppositeOrder.PendingShares > 0 {
				oppositeOrders.Push(oppositeOrder)
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)
	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()
	b.Transactions = append(b.Transactions, transaction)
}
