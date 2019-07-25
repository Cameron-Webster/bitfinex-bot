package subscriptions

import (
	"bitfinex/bot/timescale"
	"context"
	"errors"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
	"log"
	"time"
)

type BitfinexSubscriptions struct {
	Tickers     []string
	Position    int
	Initialised bool
	Connection  *websocket.Client
}

var ErrEOF = errors.New("EOF")

func (b *BitfinexSubscriptions) Next() (int, string, error) {
	b.Position++
	if b.Position > len(b.Tickers) {
		return b.Position, "", ErrEOF
	}
	return b.Position, b.Tickers[b.Position-1], nil
}

func (b *BitfinexSubscriptions) listen() {
	for obj := range b.Connection.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
			b.Initialised = false
			b.NewConnection()
		case *websocket.InfoEvent:
			if b.Initialised == false {
				b.addTradeSubs()
			}
		case *bitfinex.Trade:
			val := obj.(*bitfinex.Trade)
			timescale.InsertTradeData(float32(val.Price), float32(val.Amount), val.Pair, "bitfinex")
		default:
		}
	}
}

func (b *BitfinexSubscriptions) NewConnection() {
	p := websocket.NewDefaultParameters()
	p.ManageOrderbook = true
	conn := websocket.NewWithParams(p)
	b.Connection = conn
	err := b.Connection.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
		time.Sleep(10 * time.Second)
		b.NewConnection()
	}
	b.listen()
}

func (b *BitfinexSubscriptions) addTradeSubs() {

	b.Initialised = true

	for {
		_, ticker, err := b.Next()

		if err == ErrEOF {
			log.Printf("Done")
			break
		}

		if err != nil {
			log.Fatalf("Unknown error: %s", err)
		}

		_, err = b.Connection.SubscribeTrades(context.Background(), ticker)

		if err != nil {
			log.Printf("could not subscribe to trades: %s", err.Error(), ticker)
		} else {
			log.Printf("subscribed %s", ticker)
		}
		time.Sleep(25 * time.Millisecond)
	}
}
