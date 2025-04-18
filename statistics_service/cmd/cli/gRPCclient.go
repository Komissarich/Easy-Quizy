package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"quiz_app/internal/config"
	api "quiz_app/pkg/api/v1"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func gRPC() {

	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load config: %v", err))

	}

	grpcConn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println(fmt.Errorf("connection error: %v", err))
	}
	defer grpcConn.Close()

	cli := api.NewStatisticsClient(grpcConn)

	input := bufio.NewScanner(os.Stdin)
LOOP:
	for {
		var command string
		var opts []string
		fmt.Print("Client> ")
		input.Scan()
		command = strings.Split(input.Text(), " ")[0]
		opts = strings.Split(input.Text(), " ")[1:]
		switch command {

		case "create":
			if len(opts) == 0 {
				continue LOOP
			}
			item := strings.Join(opts[0:len(opts)-1], " ")
			quantity, err := strconv.Atoi(opts[len(opts)-1])
			if err != nil {
				fmt.Println("Server> ", err.Error())
				break LOOP
			}
			resp, err := cli.UpdateStats(ctx, &api.UpdateStatsRequest{Item: item, Quantity: int32(quantity)})
			if err != nil {
				fmt.Println("Server> ", err.Error())
				break LOOP
			}
			fmt.Printf("Server> Order {Id: %s, Item: %s, Quantity: %d} created\n", resp.Id, item, quantity)

		case "get":
			fmt.Println("Server> Found orders:")
			if len(opts) == 0 {
				continue LOOP
			}
			for _, id := range opts {
				resp, err := cli.GetOrder(ctx, &api.GetOrderRequest{Id: id})
				if err != nil {
					fmt.Println("Server> ", err.Error())
					break LOOP
				}
				fmt.Printf("{Id: %s, Item: %s, Quantity: %d}\n",
					resp.Order.Id, resp.Order.Item, resp.Order.Quantity)
			}

		case "update":
			if len(opts) == 0 {
				continue LOOP
			}
			id := opts[0]
			item := strings.Join(opts[1:len(opts)-1], " ")
			quantity, err := strconv.Atoi(opts[len(opts)-1])
			if err != nil {
				fmt.Println("Client> ", err.Error())
				break LOOP
			}
			resp, err := cli.UpdateOrder(ctx, &api.UpdateOrderRequest{Id: id, Item: item, Quantity: int32(quantity)})
			if err != nil {
				fmt.Println("Server> ", err.Error())
				break LOOP
			}
			fmt.Printf("Server> Order %s updated: {Id: %s, Item: %s, Quantity: %d}\n", resp.Order.Id, resp.Order.Id, resp.Order.Item, resp.Order.Quantity)
		case "delete":
			if len(opts) == 0 {
				continue LOOP
			}
			id := opts[0]
			for {
				resp, err := cli.DeleteOrder(ctx, &api.DeleteOrderRequest{Id: id})
				if err != nil {
					fmt.Println("Server> ", err.Error())
					break LOOP
				}
				if resp.Success {
					fmt.Printf("Server> Order %s deleted successfuly\n", id)
					break
				} else {
					fmt.Printf("Server> Order %s was not deleted, try again?[Y/n]\n", id)
					input.Scan()
					ans := input.Text()
					if ans == "n" {
						break
					}
					continue
				}
			}

		case "ls":
			resp, err := cli.ListOrders(ctx, &api.ListOrdersRequest{})
			if err != nil {
				fmt.Println("Server> ", err.Error())
				break LOOP
			}
			fmt.Println("Server> List of all orders:")
			for _, order := range resp.Orders {
				fmt.Printf("{Id: %s, Item: %s, Quantity: %d}\n",
					order.Id, order.Item, order.Quantity)
			}
		case "exit":
			fmt.Println("Connection closed")
			break LOOP
		default:
			continue LOOP
		}

	}

}
