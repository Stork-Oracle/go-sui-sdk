package client

import (
	"context"
	"strconv"
	"testing"

	"github.com/Stork-Oracle/go-sui-sdk/v2/sui_types"
	"github.com/Stork-Oracle/go-sui-sdk/v2/types"
	"github.com/shopspring/decimal"

	"github.com/Stork-Oracle/go-sui-sdk/v2/account"
	"github.com/stretchr/testify/require"
)

var (
	M1Mnemonic = "auction dose outer sorry interest daring marine tent element curious warm penalty"
	Address, _ = sui_types.NewAddressFromHex("0xf0ac79b51afa31bf3f6ca6cde5ab78b05b08362a0cfdb2923ae4c48115426589")
)

// func MainnetClient(t *testing.T) *Client {
// 	c, err := Dial(types.MainnetRpcUrl)
// 	require.NoError(t, err)
// 	return c
// }

// func TestnetClient(t *testing.T) *Client {
// 	c, err := Dial(types.TestnetRpcUrl)
// 	require.NoError(t, err)
// 	return c
// }

// func DevnetClient(t *testing.T) *Client {
// 	c, err := Dial(types.DevNetRpcUrl)
// 	require.NoError(t, err)

// 	balance, err := c.GetBalance(context.Background(), *Address, types.SUI_COIN_TYPE)
// 	require.NoError(t, err)
// 	if balance.TotalBalance.BigInt().Uint64() < SUI(0.3).Uint64() {
// 		_, err = FaucetFundAccount(Address.String(), DevNetFaucetUrl)
// 		require.NoError(t, err)
// 	}
// 	return c
// }

func LocalFundedClient(t *testing.T) *Client {
	c, err := Dial("http://localhost:9000")
	require.NoError(t, err)
	balance, err := c.GetBalance(context.Background(), *Address, types.SUI_COIN_TYPE)
	require.NoError(t, err)
	for balance.TotalBalance.BigInt().Uint64() < SUI(0.3).Uint64() {
		_, err = FaucetFundAccount(Address.String(), "http://localhost:9123/gas")
		require.NoError(t, err)
		balance, err = c.GetBalance(context.Background(), *Address, types.SUI_COIN_TYPE)
		require.NoError(t, err)
	}
	return c
}

func M1Account(t *testing.T) *account.Account {
	a, err := account.NewAccountWithMnemonic(M1Mnemonic)
	require.NoError(t, err)
	return a
}

func M1Address(t *testing.T) *suiAddress {
	return Address
}

func Signer(t *testing.T) *account.Account {
	return M1Account(t)
}

type SUI float64

func (s SUI) Int64() int64 {
	return int64(s * 1e9)
}
func (s SUI) Uint64() uint64 {
	return uint64(s * 1e9)
}
func (s SUI) Decimal() decimal.Decimal {
	return decimal.NewFromInt(s.Int64())
}
func (s SUI) String() string {
	return strconv.FormatInt(s.Int64(), 10)
}

func SuiAddressNoErr(str string) *suiAddress {
	s, _ := sui_types.NewAddressFromHex(str)
	return s
}

func ValidatorAddress(t *testing.T) *suiAddress {
	cli := LocalFundedClient(t)
	state, err := cli.GetLatestSuiSystemState(context.Background())
	require.NoError(t, err)
	return &state.ActiveValidators[0].SuiAddress
}
