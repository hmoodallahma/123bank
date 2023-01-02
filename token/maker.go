package token

import "time"

// Maker is an interface for maanging tokens
type Maker interface {
	//CreateToken creates a new token for specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken check input token valid
	VerifyToken(token string) (*Payload, error)
}
