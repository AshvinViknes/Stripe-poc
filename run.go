package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
)

func main() {
	// Set your Stripe secret key.
	stripe.Key = "sk_test_tR3PYbcVNZZ796tH88S4VQ2u"

	// Initialize the router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/home", HomePage)
	router.HandleFunc("/success", SuccessPage)
	router.HandleFunc("/cancel", CancelPage)
	router.HandleFunc("/create-checkout-session", CreateCheckoutSession).Methods("POST")

	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	port := "9080"

	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, router)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func SuccessPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/sucess.html"))
	err := t.Execute(w, nil)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, _ := io.ReadAll(r.Body)
	log.Println("Payment Sucess", string(body))
}

func CancelPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/cancel.html"))
	err := t.Execute(w, nil)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, _ := io.ReadAll(r.Body)
	log.Println("Payment cancel", string(body))
}

func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:9080/success"),
		CancelURL:  stripe.String("http://localhost:9080/cancel"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("INR"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("Example Item"),
						Description: stripe.String("An example item for testing Stripe."),
					},
					UnitAmount: stripe.Int64(50 * 100),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	session, err := session.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("SessionDetails: ", session)
	log.Println("Link expiry at", session.ExpiresAt)
	log.Println("Amount", session.AmountTotal)
	log.Println("Amount subtotal", session.AmountSubtotal)
	log.Println("Amount with tax", session.AutomaticTax)
	log.Println("Session created", session.Created)
	log.Println("Stripe checkout session ID", session.ID)
	log.Println("Sucess Url", session.SuccessURL)
	log.Println("Cancel Url", session.CancelURL)
	log.Println("Session URL", session.URL)
	log.Println("session Status", session.Status)
	log.Println("Payment Status", session.PaymentStatus)

	http.Redirect(w, r, session.URL, http.StatusSeeOther)
}
