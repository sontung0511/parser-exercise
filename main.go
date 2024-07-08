package main

import (
  "context"
  "log"
  "net/http"
  "sync"

  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/ethclient"
  "github.com/gin-gonic/gin"
)

type Server struct {
  client         *ethclient.Client
  monitoredAddrs map[string]bool
  transactions   map[string][]Transaction
  mu             sync.Mutex
}

type Transaction struct {
  Hash        string `json:"hash"`
  From        string `json:"from"`
  To          string `json:"to"`
  Value       string `json:"value"`
  BlockNumber uint64 `json:"block_number"`
}

func main() {
  // assumption already exists
  client, err := ethclient.Dial("https://mainnet.infura.io/v3/123abc456def789ghi")
  if err != nil {
    log.Fatalf("Failed to connect to the Ethereum client: %v", err)
  }

  server := &Server{
    client:         client,
    monitoredAddrs: make(map[string]bool),
    transactions:   make(map[string][]Transaction),
  }

  go server.monitorTransactions() // Start monitoring transactions in a separate goroutine

  r := gin.Default()
  r.GET("/blockNumber", server.getBlockNumber)
  r.POST("/subscribe", server.subscribeAddress)
  r.GET("/transactions", server.getTransactions)
  r.Run(":8080")
}

func (s *Server) getBlockNumber(c *gin.Context) {
  blockNumber, err := s.client.BlockNumber(context.Background())
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, gin.H{"blockNumber": blockNumber})
}

func (s *Server) subscribeAddress(c *gin.Context) {
  var req struct {
    Address string `json:"address" binding:"required"`
  }

  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  s.mu.Lock()
  s.monitoredAddrs[req.Address] = true
  s.mu.Unlock()

  c.JSON(http.StatusOK, gin.H{"message": "Address subscribed successfully"})
}

func (s *Server) getTransactions(c *gin.Context) {
  address := c.Query("address")
  if address == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
    return
  }

  s.mu.Lock()
  transactions, exists := s.transactions[address]
  s.mu.Unlock()

  if !exists {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Address not subscribed"})
    return
  }

  c.JSON(http.StatusOK, transactions)
}

func (s *Server) monitorTransactions() {
  headers := make(chan *types.Header)
  sub, err := s.client.SubscribeNewHead(context.Background(), headers)
  if err != nil {
    log.Fatalf("Failed to subscribe to new headers: %v", err)
  }

  for {
    select {
    case err := <-sub.Err():
      log.Fatalf("Subscription error: %v", err)
    case header := <-headers:
      block, err := s.client.BlockByHash(context.Background(), header.Hash())
      if err != nil {
        log.Printf("Failed to retrieve block: %v", err)
        continue
      }

      s.mu.Lock()
      for _, tx := range block.Transactions() {
        from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
        if err != nil {
          log.Printf("Failed to get sender: %v", err)
          continue
        }

        if s.monitoredAddrs[from.Hex()] || (tx.To() != nil && s.monitoredAddrs[tx.To().Hex()]) {
          s.transactions[from.Hex()] = append(s.transactions[from.Hex()], Transaction{
            Hash:        tx.Hash().Hex(),
            From:        from.Hex(),
            To:          tx.To().Hex(),
            Value:       tx.Value().String(),
            BlockNumber: block.NumberU64(),
          })
        }
      }
      s.mu.Unlock()
    }
  }
}
