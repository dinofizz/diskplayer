package diskplayer

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v3"
	"time"
)


func SendAuthenticationUrlEmail(url string) (string, error) {
	apiKey := GetConfigString(MAILGUN_API_KEY)
	domain := GetConfigString(MAILGUN_DOMAIN)
	to := GetConfigString(TO_ADDRESS)

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
