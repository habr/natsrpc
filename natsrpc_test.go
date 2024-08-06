package natsrpc_test

import (
    "context"
    "testing"
    "time"
    "encoding/json"
    "io/ioutil"

    "github.com/LeKovr/natsrpc"
	"github.com/LeKovr/natsrpc/example"
	"github.com/nats-io/nats.go"
)

type Config struct {
    NATSURL      string `json:"nats_url"`
    NATSUser     string `json:"nats_user"`
    NATSPassword string `json:"nats_password"`
}

func TestErrorHandling(t *testing.T) {
	// ... настроить сервер и клиент с некорректными параметрами ...
	// ... вызвать RPC и проверить, что вернулась ошибка ...
}
  
func TestSerialization(t *testing.T) {
	// ... создать сложную структуру данных ...
	// ... сериализовать и десериализовать структуру ...
	// ... сравнить исходную и десериализованную структуры ...
}

func TestPerformance(t *testing.T) {
    data, err := ioutil.ReadFile("config.json")
    if err != nil {
        t.Fatal(err)
    }

    var config Config
    err = json.Unmarshal(data, &config)
    if err != nil {
        t.Fatal(err)
    }

    // Connect to NATS using the configuration
    opts := []nats.Option{}
    if config.NATSUser != "" && config.NATSPassword != "" {
        opts = append(opts, nats.UserInfo(config.NATSUser, config.NATSPassword))
    }
    nc, err := nats.Connect(config.NATSURL, opts...)
    if err != nil {
        t.Fatal(err)
    }
    defer nc.Close()

    // Create server and client using the NATS connection
    server, err := natsrpc.NewServer(nc)
    if err != nil {
        t.Fatal(err)
    }
    defer server.Close(context.Background())

    client := natsrpc.NewClient(nc)

    // Register the service
    svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})
    if err != nil {
        t.Fatal(err)
    }
    defer svc.Close()

    // Create the client
    cli := example.NewGreetingNRClient(client)

    // Number of requests to send
    numRequests := 10000

    // Start timer
    start := time.Now()

    // Send requests concurrently
    for i := 0; i < numRequests; i++ {
        go func() {
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()
            _, err := cli.Hello(ctx, &example.HelloRequest{Name: "bruce"})
            if err != nil {
                t.Errorf("Error sending request: %v", err)
            }
        }()
    }

    // Wait for all requests to complete
    time.Sleep(time.Second * 2) // Adjust as needed

    // Calculate elapsed time
    elapsed := time.Since(start)

    // Print performance metrics
    t.Logf("Sent %d requests in %v", numRequests, elapsed)
    t.Logf("Requests per second: %.2f", float64(numRequests)/elapsed.Seconds())
}

type HelloSvc struct{}

func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
    return &example.HelloReply{Message: "hello " + req.Name}, nil
}