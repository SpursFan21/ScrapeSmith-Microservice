package utils

import (
	"os"

	"github.com/stripe/stripe-go/v75"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
}
