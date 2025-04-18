package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"quiz_app/internal/config"
	api "quiz_app/pkg/api/v1"
)

var CreateReqBody = api.CreateOrderRequest{
	Item:     "apple",
	Quantity: 42,
}

var GetReqBody = api.GetOrderRequest{
	Id: "1",
}

func HTTP() {

	cfg, err := config.New()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load config: %v", err))

	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(CreateReqBody)
	if err != nil {
		fmt.Println(fmt.Errorf("encoding error: %v", err))
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://0.0.0.0:%d/v1/OrderService", cfg.HTTPPort), &buf)
	if err != nil {
		fmt.Println(fmt.Errorf("request error: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Errorf("response error: %v", err))
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bytes))

	err = json.NewEncoder(&buf).Encode(GetReqBody)
	if err != nil {
		fmt.Println(fmt.Errorf("encoding error: %v", err))
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("http://0.0.0.0:%d/v1/OrderService/%s", cfg.HTTPPort, GetReqBody.Id), &buf)
	if err != nil {
		fmt.Println(fmt.Errorf("request error: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(fmt.Errorf("response error: %v", err))
	}
	defer resp.Body.Close()

	bytes, _ = io.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}
