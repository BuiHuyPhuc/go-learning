package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func placeOrderWithContext(ctx context.Context, orderId string) error {
	fmt.Printf("Start handle order: %s\n", orderId)

	select {
	case <-time.After(time.Second * 3):
		fmt.Printf("Finish handle order %s after 3s\n", orderId)
		return nil

	case <-ctx.Done():
		fmt.Printf("Cancelled order %s: %v\n", orderId, ctx.Err())
		return ctx.Err()
	}
}

func OrderHandlerWithContext(w http.ResponseWriter, r *http.Request) {
	orderId := "Order-123"

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*4)
	defer cancel()

	err := placeOrderWithContext(ctx, orderId)
	if err != nil {
		fmt.Printf("Handle order %s failed: %v\n", orderId, err)
		http.Error(w, "Failed to handle order", http.StatusGatewayTimeout)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order success!"))
}

// func placeOrderWithoutContext(orderId string) error {
// 	fmt.Printf("Start handle order: %s\n", orderId)

// 	time.Sleep(time.Second * 3)

// 	fmt.Printf("Finish handle order %s after 3s\n", orderId)
// 	return nil
// }

// func OrderHandlerSelect(w http.ResponseWriter, r *http.Request) {
// 	orderId := "Order-123"
// 	ch := make(chan error, 1)

// 	go func() {
// 		err := placeOrderWithoutContext(orderId)
// 		ch <- err
// 	}()

// 	select {
// 	case err := <-ch:
// 		if err != nil {
// 			fmt.Printf("Handle order %s failed\n", orderId)
// 			http.Error(w, "Failed to handle order", http.StatusInternalServerError)
// 		  return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Order success!"))

// 	case <-time.After(time.Second * 2):
// 		fmt.Printf("Handle order %s timeout 2s\n", orderId)
// 		http.Error(w, "Timeout to handle order, please again", http.StatusGatewayTimeout)
// 	}
// }

func main() {
	// http.HandleFunc("/order", OrderHandlerWithContext)
	// fmt.Println("Server run to url http://locahost:8890")
	// log.Fatal(http.ListenAndServe(":8890", nil))

	// context root
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	orderId := "Order-123"
	err := placeOrderWithContext(ctx, orderId)
	if err != nil {
		fmt.Printf("Handle order %s failed: %v\n", orderId, err)
	} else {
		fmt.Printf("Handle order %s success\n", orderId)
	}
}
