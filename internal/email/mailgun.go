package email

import (
	"context"
	"fmt"
	"github.com/dinofizz/diskplayer/internal/config"
	"github.com/mailgun/mailgun-go/v3"
	"time"
)

const MAILGUN_API_KEY = "mailgun.api_key"
const MAILGUN_DOMAIN = "mailgun.domain"
const TO_ADDRESS = "mailgun.to_address"

func SendAuthenticationUrlEmail(url string) (string, error) {
	apiKey := config.GetConfigString(MAILGUN_API_KEY)
	domain := config.GetConfigString(MAILGUN_DOMAIN)
	to := config.GetConfigString(TO_ADDRESS)

	from := fmt.Sprintf("Diskplayer <diskplayer@%s>", domain)

	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		from,
		"Spotify Authentication URL",
		url,
		to,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}
