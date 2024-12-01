package config_test

import (
	"testing"

	"github.com/ichimei0125/gotradecrypto/internal/config"
)

func TestGetConfig(t *testing.T) {
	c := config.GetConfig()
	t.Log(c.Trade.InvestMoney)

}
