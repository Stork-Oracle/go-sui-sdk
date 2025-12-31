package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFaucetFundAccount(t *testing.T) {
	addr := Address.String()
	res, err := FaucetFundAccount(addr, "http://localhost:9123/gas")
	require.NoError(t, err)
	t.Log("hash = ", res)
}
