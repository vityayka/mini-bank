package mail

import (
	"bank/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	gmail := NewGmailSender(config.GmailName, config.GmailFrom, config.GmailAccPassword)

	err = gmail.Send(
		"Test",
		`
			<h1>Test mail</h1>
		`,
		[]string{"gegviktor@yandex.ru"},
		nil,
		nil,
		[]string{"../README.md"},
	)

	require.NoError(t, err)
}
