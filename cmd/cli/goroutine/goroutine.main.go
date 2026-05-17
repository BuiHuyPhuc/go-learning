package main

/* Singleton instance khi sử dụng nhiều luồng đồng thời có khả năng bị lỗi tạo nhiều lần
import (
	"fmt"
	"sync"
)

type Singleton struct{}

var singleton *Singleton

// -----> Fix new singleton instance 1 lần duy nhất sử dụng sync.Once
var once sync.Once

func GetInstance() *Singleton {
	if singleton == nil {
		fmt.Println("New instance...")
		singleton = &Singleton{}
	}
	// once.Do(func() {
	// 	fmt.Println("New instance...")
	// 	singleton = &Singleton{}
	// })
	return singleton
}

// -----> Nhược điểm sync.Once có thể xảy ra deadlock
// once.Do(func() {
//   once.Do(func() {
//     fmt.Println("New...")
//   })
// })

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := GetInstance()
			fmt.Printf("Singleton instance address: %p\n", s)
		}()
	}

	wg.Wait()
}
*/

/* Data race là quả bom hẹn giờ vì sử dụng biến bị tranh chấp bởi nhiều luồng
import (
	"fmt"
	"sync"
)

var counter int

// -----> Fix data race sử dụng mutex
// var mu sync.Mutex

func increment() {
	// mu.Lock()
	// defer mu.Unlock()
	counter++
}

func main() {
	// Đảm bảo tất cả luồng chạy hoàn thành thì với kết thúc luồng main
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	wg.Wait()
	fmt.Println("Final counter:", counter)

	// go run -race main.go
}
*/

/* Multiple channel and goroutines
type Message struct {
	OrderId string
	Title   string
	Price   int
}

func buyTicket(ch chan<- Message, orders []Message) {
	for _, order := range orders {
		time.Sleep(time.Second * 1)
		fmt.Printf("Send buy:::%s\n", order.OrderId)
		ch <- order
	}
	close(ch)
}
func cancelTicket(ch chan<- string, orderIds []string) {
	for _, orderId := range orderIds {
		time.Sleep(time.Second * 2)
		fmt.Printf("Send cancel:::%s\n", orderId)
		ch <- orderId
	}
	close(ch)
}
func handlerOrder(orderCh <-chan Message, cancelCh <-chan string) {
	for {
		// Nếu không dung select thì channel sẽ chờ nhau hoàn thành thì mới tiếp tục gây ra việc chậm trễ
		select {
		// delay 1s
		case order, orderOK := <-orderCh:
			if orderOK {
				fmt.Printf("Handler buy:::%s, %s, %d\n", order.OrderId, order.Title, order.Price)
			} else {
				fmt.Println("Order channel closed...")
				orderCh = nil
			}

		// delay 2s
		case orderId, cancelOK := <-cancelCh:
			if cancelOK {
				fmt.Printf("Handler cancel:::%s\n", orderId)
			} else {
				fmt.Println("Cancel channel closed...")
				cancelCh = nil
			}
		}
	}
}

func main() {
	buyCh := make(chan Message)
	cancelCh := make(chan string)

	buyOrders := []Message{
		{OrderId: "Order-1", Title: "Go learning 1", Price: 30},
		{OrderId: "Order-2", Title: "Go learning 2", Price: 40},
		{OrderId: "Order-3", Title: "Go learning 3", Price: 50},
	}

	cancelOrderIds := []string{"Order-1", "Order-3"}

	go buyTicket(buyCh, buyOrders)
	go cancelTicket(cancelCh, cancelOrderIds)

	// handler
	go handlerOrder(buyCh, cancelCh)

	time.Sleep(time.Second * 10)
	fmt.Println("End buying and cancelling...")
}
*/

/* Pub/Sub use channel and goroutines
func publisher(ch chan<- Message, orders []Message) {
	for _, order := range orders {
		fmt.Printf("Pub:::%s\n", order.OrderId)
		ch <- order
		time.Sleep(time.Second * 1)
	}
	close(ch)
}

func subscriber(ch <-chan Message, userName string) {
	for msg := range ch {
		fmt.Printf("Sub:::%s, %s, %s, %d\n", userName, msg.OrderId, msg.Title, msg.Price)
		time.Sleep(time.Second * 1)
	}
}

func main() {
	// 1. channel order
	orderCh := make(chan Message)

	// 2. create orders
	orders := []Message{
		{OrderId: "Order-1", Title: "Go learning 1", Price: 30},
		{OrderId: "Order-2", Title: "Go learning 2", Price: 40},
		{OrderId: "Order-3", Title: "Go learning 3", Price: 50},
	}

	// 3. send order to pub
	go publisher(orderCh, orders)
	go subscriber(orderCh, "PhucNgo")

	time.Sleep(time.Second * 5)
	fmt.Println("End pub sub...")
}
*/

/* Basic channel and goroutines
import "fmt"

type Course struct {
	Title string
	Price int
}

func main() {
	// 1. add channel
	ch := make(chan Course)

	// 2. create goroutine
	go func() {
		course := Course{Title: "Go learning", Price: 30}
		ch <- course // send data to channel
	}()

	c := <-ch // receive data from channel

	fmt.Printf("Receive course: %s, price: %d\n", c.Title, c.Price)
}
*/

/* Waitgroup to sync
import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("Starting...")
	var wg sync.WaitGroup

	ids := []int{1, 2, 3, 4, 5}

	start := time.Now()
	for _, id := range ids {
		wg.Add(1)
		go getProductById(id, &wg)
	}

	wg.Wait()
	fmt.Println("Finished...", time.Since(start))
}

func getProductById(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("https://fakestoreapi.com/products/%d", id)
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	fmt.Printf(">>> Data productId %d: %s\n", id, string(body))
}
*/
