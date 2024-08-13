package etherscan

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"eth-transaction-api/config"
)

type Transaction struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

type EtherscanResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

// stETH contract details
const stETHAddress = "0xae7ab96520de3a18e5e111b5eaab095312d7fe84"
const stETHABI = `[ ... ABI definition here ... ]` // Replace with actual ABI

func GetTransactions(address string) (*EtherscanResponse, error) {
	cfg := config.LoadConfig()

	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=10&sort=desc&apikey=%s",
		address, cfg.EtherscanApiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var etherscanResponse EtherscanResponse
	if err := json.Unmarshal(body, &etherscanResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	if etherscanResponse.Status != "1" {
		log.Printf("Etherscan API returned error: %s", etherscanResponse.Message)
	}

	return &etherscanResponse, nil
}
