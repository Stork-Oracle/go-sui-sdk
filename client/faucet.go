package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/stork-oracle/go-sui-sdk/sui_types"
)

const (
	DevNetFaucetUrl  = "https://faucet.devnet.sui.io/gas"
	TestNetFaucetUrl = "https://faucet.testnet.sui.io/gas"
)

func FaucetFundAccount(address string, faucetUrl string) (string, error) {
	_, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return "", err
	}

	paramJson := fmt.Sprintf(`{"FixedAmountRequest":{"recipient":"%v"}}`, address)
	request, err := http.NewRequest(http.MethodPost, faucetUrl, bytes.NewBuffer([]byte(paramJson)))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return "", fmt.Errorf("post %v response code = %v", faucetUrl, res.Status)
	}
	defer res.Body.Close()

	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		// Old format (pre-v1.0)
		TransferredGasObjects []struct {
			Amount uint64 `json:"amount"`
			Id     string `json:"id"`
			Digest string `json:"transferTxDigest"`
		} `json:"transferredGasObjects,omitempty"`
		// New format (v1.0+)
		Status    string `json:"status,omitempty"`
		CoinsSent []struct {
			Amount uint64 `json:"amount"`
			Id     string `json:"id"`
			Digest string `json:"transferTxDigest"`
		} `json:"coins_sent,omitempty"`
		Error string `json:"error,omitempty"`
	}
	err = json.Unmarshal(resByte, &response)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(response.Error) != "" {
		return "", errors.New(response.Error)
	}

	// Try new format first
	if len(response.CoinsSent) > 0 {
		return response.CoinsSent[0].Digest, nil
	}
	// Fall back to old format
	if len(response.TransferredGasObjects) > 0 {
		return response.TransferredGasObjects[0].Digest, nil
	}

	return "", errors.New("transaction not found")
}
