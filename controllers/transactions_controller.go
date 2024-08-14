package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"eth-transaction-api/etherscan"
	"eth-transaction-api/models"
	"eth-transaction-api/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetAccounts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accounts []models.Account
		if err := db.Find(&accounts).Error; err != nil {
			http.Error(w, "Failed to retrieve accounts", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"data":  accounts,
			"count": len(accounts),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}

// GetAccountTransactions handles the request to get an account's transactions
func GetAccountTransactions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountUuid := vars["accountUuid"]

		// Fetch the account from the database using accountUuid
		var account models.Account
		if err := db.Where("account_uuid = ?", accountUuid).First(&account).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Account not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to retrieve account", http.StatusInternalServerError)
			return
		}

		// Fetch transactions from Etherscan
		etherscanTransactions, err := etherscan.GetTransactions(account.Address)
		if err != nil {
			http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
			return
		}

		// Convert and save new transactions
		transactions := make([]map[string]interface{}, 0, len(etherscanTransactions.Result))
		for _, tx := range etherscanTransactions.Result {
			amountInEth := utils.WeiToEth(tx.Value, 18)

			txType := "withdrawal"
			if tx.To == account.Address {
				txType = "deposit"
			}

			timestamp, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
			utcTime := time.Unix(timestamp, 0).UTC()

			// Check if the transaction already exists
			var existingTransaction models.Transaction
			if err := db.Where("tx_hash = ?", tx.Hash).First(&existingTransaction).Error; err == gorm.ErrRecordNotFound {
				// Transaction does not exist, so create a new one
				amount, err := strconv.ParseFloat(amountInEth, 64)
				if err != nil {
					http.Error(w, "Failed to parse amount", http.StatusInternalServerError)
					return
				}

				newTransaction := models.Transaction{
					Amount:    amount,
					Token:     "ETH",
					Timestamp: utcTime,
					Sender:    tx.From,
					Receiver:  tx.To,
					Type:      txType,
					TxHash:    tx.Hash,
					AccountID: account.ID,
				}

				if err := newTransaction.Save(db); err != nil {
					http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
					return
				}
			}

			// Append to the response list
			transaction := map[string]interface{}{
				"id":          tx.Hash,
				"accountUuid": accountUuid,
				"toAddress":   tx.To,
				"fromAddress": tx.From,
				"type":        txType,
				"amount":      amountInEth,
				"symbol":      "ETH",
				"decimal":     18,
				"timestamp":   utcTime.Format(time.RFC3339),
				"txnHash":     tx.Hash,
			}

			transactions = append(transactions, transaction)
		}

		// sort by timestamp
		// Sort transactions by timestamp (newest to oldest)
		sort.Slice(transactions, func(i, j int) bool {
			// Parse the timestamp strings to time.Time objects
			timestampI, err := time.Parse(time.RFC3339, transactions[i]["timestamp"].(string))
			if err != nil {
				log.Println("Error parsing timestamp:", err)
				return false
			}

			timestampJ, err := time.Parse(time.RFC3339, transactions[j]["timestamp"].(string))
			if err != nil {
				log.Println("Error parsing timestamp:", err)
				return false
			}

			// Compare the two time.Time values
			return timestampI.After(timestampJ)
		})

		// Create response with data and count
		response := map[string]interface{}{
			"data":  transactions,
			"count": len(transactions),
		}

		// Return the formatted response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// not finish....
// stETH contract details
const stETHAddress = "0xae7ab96520de3a18e5e111b5eaab095312d7fe84"
const stETHABI = `[ ... ABI definition here ... ]` // Replace with actual ABI

// GetPooledETHAndShares returns the total pooled ETH and total shares from the stETH token contract
func GetPooledETHAndShares(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + os.Getenv("INFURA_PROJECT_ID"))
		if err != nil {
			http.Error(w, "Failed to connect to Ethereum client", http.StatusInternalServerError)
			return
		}

		contractAddress := common.HexToAddress(stETHAddress)
		parsedABI, err := abi.JSON(strings.NewReader(stETHABI))
		if err != nil {
			http.Error(w, "Failed to parse contract ABI", http.StatusInternalServerError)
			return
		}

		// Create a callOpts object to make a read-only call to the Ethereum network
		// Define the variables to store the returned values
		var totalPooledETH, totalShares big.Int

		// Pack the method name and arguments (if any) into the input data
		data, err := parsedABI.Pack("getPooledEth")
		if err != nil {
			http.Error(w, "Failed to pack input data", http.StatusInternalServerError)
			return
		}

		// Make the call to the Ethereum network
		bPool, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		}, nil)
		if err != nil {
			http.Error(w, "Failed to call contract", http.StatusInternalServerError)
			return
		}

		// Unpack the returned data into the totalPooledETH variable
		err = parsedABI.UnpackIntoInterface(&totalPooledETH, "getPooledEth", data)
		if err != nil {
			http.Error(w, "Failed to unpack output data", http.StatusInternalServerError)
			return
		}

		// Pack and call the getTotalShares method
		data, err = parsedABI.Pack("getTotalShares")
		if err != nil {
			http.Error(w, "Failed to pack input data", http.StatusInternalServerError)
			return
		}

		bSh, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		}, nil)
		if err != nil {
			http.Error(w, "Failed to call contract", http.StatusInternalServerError)
			return
		}

		fmt.Println("bPool", bPool)
		fmt.Println("bSh", bSh)

		err = parsedABI.UnpackIntoInterface(&totalShares, "getTotalShares", data)
		if err != nil {
			http.Error(w, "Failed to unpack output data", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"total_pooled_eth": totalPooledETH.String(),
			"total_shares":     totalShares.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// GetLast5Depositors returns the last 5 addresses that deposited into the stETH pool
func GetLast5Depositors(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + os.Getenv("INFURA_PROJECT_ID"))
		if err != nil {
			http.Error(w, "Failed to connect to Ethereum client", http.StatusInternalServerError)
			return
		}

		contractAddress := common.HexToAddress(stETHAddress)
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(0),
			ToBlock:   nil,
			Addresses: []common.Address{contractAddress},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
			return
		}

		var last5Depositors []string
		for _, vLog := range logs {
			if len(last5Depositors) < 5 {
				last5Depositors = append(last5Depositors, vLog.Address.Hex())
			} else {
				break
			}
		}

		response := map[string]interface{}{
			"last_5_depositors": last5Depositors,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
