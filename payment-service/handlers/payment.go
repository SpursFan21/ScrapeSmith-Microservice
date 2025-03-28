package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentintent"
	"github.com/stripe/stripe-go/v75/price"
	"github.com/stripe/stripe-go/v75/product"
)

// CreatePaymentIntent handles creating a payment intent
func CreatePaymentIntent(c *fiber.Ctx) error {
	productID := c.Params("id") // Get the product ID from URL parameter

	// Retrieve the product from Stripe
	stripeProduct, err := product.Get(productID, nil)
	if err != nil {
		log.Printf("Error fetching product details: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product details",
		})
	}

	// Retrieve the price for the product
	params := &stripe.PriceListParams{
		Product: stripe.String(productID), // Use the product ID to find its price
	}
	priceIterator := price.List(params)

	var productPrice *stripe.Price
	for priceIterator.Next() {
		productPrice = priceIterator.Price()
		break // Take the first available price
	}

	if productPrice == nil {
		log.Printf("No price found for product: %s", productID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product price not found",
		})
	}

	// Create Payment Intent with the retrieved price
	stripeAmount := productPrice.UnitAmount

	paramsIntent := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(stripeAmount),
		Currency: stripe.String(string(productPrice.Currency)),
	}

	intent, err := paymentintent.New(paramsIntent)
	if err != nil {
		log.Printf("Stripe payment intent error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment intent",
		})
	}

	return c.JSON(fiber.Map{
		"clientSecret": intent.ClientSecret,
		"productName":  stripeProduct.Name,
		"price":        float64(stripeAmount) / 100,
		"currency":     productPrice.Currency,
	})
}

// GetProductDetails handles retrieving product details by ID
func GetProductDetails(c *fiber.Ctx) error {
	productID := c.Params("id") // Get the product ID from URL parameter

	// Retrieve the product from Stripe
	stripeProduct, err := product.Get(productID, nil)
	if err != nil {
		log.Printf("Error fetching product details: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product details",
		})
	}

	// Retrieve the price for the product
	params := &stripe.PriceListParams{
		Product: stripe.String(productID),
	}
	priceIterator := price.List(params)

	var productPrice *stripe.Price
	for priceIterator.Next() {
		productPrice = priceIterator.Price()
		break // Take the first available price
	}

	if productPrice == nil {
		log.Printf("No price found for product: %s", productID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product price not found",
		})
	}

	return c.JSON(fiber.Map{
		"id":          stripeProduct.ID,
		"name":        stripeProduct.Name,
		"description": stripeProduct.Description,
		"price":       float64(productPrice.UnitAmount) / 100,
		"currency":    productPrice.Currency,
	})
}
