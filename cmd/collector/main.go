// Package main implements a WebSocket-based trade collector for Bitget exchange.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

// Constants for WebSocket connection and reconnection delay.
const (
	wsURL          = "wss://ws.bitget.com/spot/v1/stream"
	reconnectDelay = 5 * time.Second
)

// Command-line flags.
var (
	pair  = flag.String("pair", "", "Trading pair (e.g., BTCUSDT)")
	debug = flag.Bool("debug", false, "Enable debug logging")
)

// Trade represents a single trade from the Bitget WebSocket stream.
type Trade struct {
	TradeID   string `json:"tradeId"`
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"`
	Timestamp int64  `json:"ts"`
}

// main is the entry point of the application.
func main() {
	flag.Parse()

	// Validate the trading pair argument.
	if *pair == "" {
		log.Fatal("Error: trading pair must be specified with -pair")
	}

	// Construct the database file path based on the trading pair.
	dbFile := fmt.Sprintf("/var/lib/bitget-collector/trades_%s.db", *pair)

	// Initialize the database.
	db, err := initDB(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run the collector in a loop with reconnection on failure.
	for {
		err := runCollector(db, *pair)
		if err != nil {
			log.Printf("Error: %v. Reconnecting in %v", err, reconnectDelay)
			time.Sleep(reconnectDelay)
		}
	}
}

// initDB initializes the SQLite database with WAL mode and creates the trades table.
func initDB(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Enable Write-Ahead Logging (WAL) for non-blocking database operations.
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, err
	}

	// Create the trades table if it doesn't exist.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS trades (
			trade_id TEXT,
			symbol TEXT,
			price REAL,
			size REAL,
			side TEXT,
			timestamp INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// runCollector connects to the Bitget WebSocket and processes trade messages.
func runCollector(db *sql.DB, pair string) error {
	// Establish WebSocket connection.
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	// Subscribe to the trade channel for the specified pair.
	subscribeMsg := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{"instType": "sp", "channel": "trade", "instId": pair},
		},
	}
	err = c.WriteJSON(subscribeMsg)
	if err != nil {
		return err
	}

	// Process incoming messages.
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return err
		}

		if *debug {
			log.Printf("Received message: %s", message)
		}

		var msg map[string]interface{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			if *debug {
				log.Printf("Unmarshal error: %v", err)
			}
			continue
		}

		// Extract the "data" field, which contains an array of trades.
		if data, ok := msg["data"].([]interface{}); ok {
			for _, item := range data {
				// Each trade is an array of strings: [timestamp, price, size, side].
				tradeArray, ok := item.([]interface{})
				if !ok || len(tradeArray) != 4 {
					if *debug {
						log.Printf("Invalid trade format: %v", item)
					}
					continue
				}

				// Extract trade fields.
				timestampStr, _ := tradeArray[0].(string)
				priceStr, _ := tradeArray[1].(string)
				sizeStr, _ := tradeArray[2].(string)
				side, _ := tradeArray[3].(string)

				// Convert fields to appropriate types.
				timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
				if err != nil {
					if *debug {
						log.Printf("Timestamp parse error: %v", err)
					}
					continue
				}
				price, err := strconv.ParseFloat(priceStr, 64)
				if err != nil {
					if *debug {
						log.Printf("Price parse error: %v", err)
					}
					continue
				}
				size, err := strconv.ParseFloat(sizeStr, 64)
				if err != nil {
					if *debug {
						log.Printf("Size parse error: %v", err)
					}
					continue
				}

				// Extract symbol from the "arg" field.
				arg, ok := msg["arg"].(map[string]interface{})
				if !ok {
					if *debug {
						log.Printf("Invalid arg format: %v", msg["arg"])
					}
					continue
				}
				symbol, _ := arg["instId"].(string)

				// Generate a unique trade ID.
				tradeID := fmt.Sprintf("%d-%d", timestamp, rand.Intn(1000))

				// Insert trade into the database.
				_, err = db.Exec(`
					INSERT INTO trades (trade_id, symbol, price, size, side, timestamp)
					VALUES (?, ?, ?, ?, ?, ?)
				`, tradeID, symbol, price, size, side, timestamp)
				if err != nil {
					if *debug {
						log.Printf("Database insert error: %v", err)
					}
				} else if *debug {
					log.Printf("Trade saved: %s, %s, %.2f, %.4f, %s, %d", tradeID, symbol, price, size, side, timestamp)
				}
			}
		}
	}
}
