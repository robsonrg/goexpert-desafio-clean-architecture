package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/robsonrg/goexpert-desafio-clean-architecture/configs"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/event"
	event_handler "github.com/robsonrg/goexpert-desafio-clean-architecture/internal/event/handler"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/infra/database"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/infra/graph"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/infra/grpc/pb"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/infra/grpc/service"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/internal/infra/webserver"
	usecase "github.com/robsonrg/goexpert-desafio-clean-architecture/internal/usecase/orders"
	"github.com/robsonrg/goexpert-desafio-clean-architecture/pkg/events"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &event_handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := setupCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := setupListOrderUseCase(db)

	startWebServer(configs.WebServerPort, createOrderUseCase, listOrderUseCase)
	startGRPCServer(configs.GRPCServerPort, *createOrderUseCase, *listOrderUseCase)
	startGraphQLServer(configs.GraphQLServerPort, *createOrderUseCase, *listOrderUseCase)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func setupCreateOrderUseCase(db *sql.DB, eventDispatcher *events.EventDispatcher) *usecase.CreateOrderUseCase {
	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	return usecase.NewCreateOrderUseCase(orderRepository, orderCreated, eventDispatcher)
}

func setupListOrderUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	orderRepository := database.NewOrderRepository(db)
	return usecase.NewListOrderUseCase(orderRepository)
}

func startWebServer(
	port string,
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrderUseCase *usecase.ListOrderUseCase,
) {
	ws := webserver.NewServer(port)
	webOrderHandler := webserver.NewWebOrderHandler(createOrderUseCase, listOrderUseCase)

	ws.AddHandler(webserver.NewRoute("/orders", "POST", webOrderHandler.Create))
	ws.AddHandler(webserver.NewRoute("/orders", "GET", webOrderHandler.GetOrders))

	fmt.Println("Starting Web server on port", port)
	go ws.Start()
}

func startGRPCServer(port string, createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) {
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(createOrderUseCase, listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
}

func startGraphQLServer(port string, createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) {
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CreateOrderUseCase: createOrderUseCase,
			ListOrderUseCase:   listOrderUseCase,
		},
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start GraphQL server: %v", err)
	}
}
