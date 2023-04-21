package insurance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type Insurance struct {
	Name string
	Price int
	URL string
	Expect string
	Content string
	Purchased bool
	PaidOut bool
}

var API_KEY = "QVBJX0tFWTo0N2NmZDZmOGZlODRhYTNmMTljZjUzYmYxOWU3ZmFkOTo1MGM5N2VmNmE4Y2YwYjg5MDUxNDY1NWM1ZDRhNDVlOQ=="

var Master_wallet = "1015033390"

var Wallet_ID = "aa3ea49b-142d-5561-ada8-6dbcc224228f"
var Wallet_Address = "0x4e7c398f9eb8be615d33af76c7ca34c1c3f9d99a"


var Insurances = []*Insurance{
	{
		Name: "Service guarantee",
		Price: 0,
		Content: "If localhost:12345/test is not 'up', you get 5 USDC",
		URL : "http://localhost:12345/test",
		Expect: "up",
		Purchased: true,
		PaidOut: false,
	},
	{
		Name: "Delivery guarantee",
		Price: 0,
		URL: "http://localhost:12345/delivery",
		Expect: "delivered",
		Content: "If localhost:12345/delivery response is not 'delivered', you get 5 USDC",
		Purchased: true,
		PaidOut: false,
	},
	{
		Name: "Weather guarantee",
		Price: 10,
		URL: "http://localhost:12345/weather",
		Expect: "sunny",
		Content: "If localhost:12345/weather response is not 'sunny', you get 5 USDC",
		Purchased: false,
		PaidOut: false,
	},
}

type Source struct {
	Type string `json:"type"`
	ID string `json:"id"`
}
type Amount struct{
	Currency string `json:"currency"`
	Amount string `json:"amount"`
}
type Request struct {
	Source Source `json:"source"`
	Destination Source `json:"destination"`
	Amount Amount `json:"amount"`
	IdempotencyKey string `json:"idempotencyKey"`
}

type Response struct {
	Test bool
}


func Validate() {
	for _, insurance := range Insurances{
		if !insurance.Purchased || insurance.PaidOut {
			continue
		}
		resp, err := http.Get(insurance.URL)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		if string(body) != insurance.Expect {
			fmt.Println(insurance.URL)
			fmt.Println(string(body))
			fmt.Println(insurance.Expect)
			fmt.Println("need to pay money")
			PayMoney(context.Background())
			insurance.PaidOut = true
			continue
		}
	}

}


//encore:api public path=/validater
func PayMoney(ctx context.Context) (*Response, error) {
	url := "https://api-sandbox.circle.com/v1/payouts"
	reqData := Request{
		Source: Source{
			Type: "wallet",
			ID: Master_wallet,
		},
		Destination: Source{
			Type: "address_book",
			ID: "aa3ea49b-142d-5561-ada8-6dbcc224228f",
		},
		Amount: Amount{
			Amount: "0.1",
			Currency: "USD",
		},
		IdempotencyKey: uuid.NewString(),
	}
    data, _ := json.Marshal(reqData)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
    req.Header.Set("authorization", "Bearer QVBJX0tFWTo0N2NmZDZmOGZlODRhYTNmMTljZjUzYmYxOWU3ZmFkOTo1MGM5N2VmNmE4Y2YwYjg5MDUxNDY1NWM1ZDRhNDVlOQ==")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
	return nil, nil
}
