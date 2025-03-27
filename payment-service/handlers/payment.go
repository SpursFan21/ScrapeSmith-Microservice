package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
)

func CreateCheckoutSession(c *fiber.Ctx) error {
	stripeAmount := int64(1000) // Amount in cents (e.g., $10.00)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("1 Job"),
					},
					UnitAmount: stripe.Int64(stripeAmount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String(os.Getenv("SUCCESS_URL")),
		CancelURL:  stripe.String(os.Getenv("CANCEL_URL")),
	}

	session, err := session.New(params)
	if err != nil {
		log.Printf("Stripe checkout session error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create checkout session",
		})
	}

	return c.JSON(fiber.Map{
		"sessionId": session.ID,
		"url":       session.URL,
	})
}
