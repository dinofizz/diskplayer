package email

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v3"
	"log"
	"os"
	"time"
)

const MAILGUN_API_KEY = "MAILGUN_API_KEY"
const MAILGUN_DOMAIN = "MAILGUN_DOMAIN"
const EMAIL_ADDRESS = "EMAIL_ADDRESS"

func SendAuthenticationUrlEmail(url string) (string, error) {
	apiKey := os.Getenv(MAILGUN_API_KEY)
	domain := os.Getenv(MAILGUN_DOMAIN)
	to := os.Getenv(EMAIL_ADDRESS)

	if apiKey == "" {
		log.Fatalf("Environment variable %s is empty.", MAILGUN_API_KEY)
	}

	if domain == "" {
		log.Fatalf("Environment variable %s is empty.", MAILGUN_DOMAIN)
	}

	if to == "" {
		log.Fatalf("Environment variable %s is empty.", EMAIL_ADDRESS)
	}

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
