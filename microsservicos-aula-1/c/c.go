package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/hashicorp/go-retryablehttp"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid := coupons.Check(coupon)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

	if(result.Status == "valid") {
		result := sendEmail()
		if(result.Status != "ok") {
			fmt.Fprintf(w, result.Status)
		}
	}

}

func sendEmail() Result {
	result := Result{Status: "ok"}
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 2

	res, err := retryClient.PostForm("http://localhost:9093", nil)
	if err != nil {
		result := Result{Status: "Servidor de e-mail fora do ar!"}
		fmt.Println(result.Status)
		addOnQueue("email@domain.com")
		return result
	}

	defer res.Body.Close()

	return result
}

func addOnQueue(email string) {
	fmt.Println(email)
}
