package imgur

import (
	"errors"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

// ClientAccount describe authontification
type ClientAccount struct {
	clientID    string // client ID received after registration
	accessToken string // is your secret key used to access the user's data
}

// Client used to for go-imgur
type Client struct {
	Log          zerolog.Logger
	httpClient   *http.Client
	imgurAccount ClientAccount
	rapidAPIKey  string
}

// NewClient simply creates an imgur client. RapidAPIKEY is "" if you are using the free API.
func NewClient(httpClient *http.Client, clientID string, rapidAPIKey string) (*Client, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if len(clientID) == 0 {
		msg := "imgur client ID is empty"
		logger.Error().Msg(msg)
		return nil, errors.New(msg)
	}

	if len(rapidAPIKey) == 0 {
		logger.Info().Msg("rapid api key is empty")
	}

	return &Client{
		httpClient:  httpClient,
		Log:         logger,
		rapidAPIKey: rapidAPIKey,
		imgurAccount: ClientAccount{
			clientID: clientID,
		},
	}, nil
}
