package main

import (
	"log"
	"net"
	"os"

	c "github.com/forum-gamers/glowing-octo-robot/controllers"
	"github.com/forum-gamers/glowing-octo-robot/database"
	transactionProto "github.com/forum-gamers/glowing-octo-robot/generated/transaction"
	h "github.com/forum-gamers/glowing-octo-robot/helpers"
	"github.com/forum-gamers/glowing-octo-robot/interceptor"
	"github.com/forum-gamers/glowing-octo-robot/pkg/transaction"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	h.PanicIfError(godotenv.Load())
	database.Conn()

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "50059"
	}

	lis, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		log.Fatalf("Failed to listen : %s", err.Error())
	}

	transactionRepo := transaction.NewTransactionRepo()

	interceptor := interceptor.NewInterCeptor()
	getUser := interceptor.GetUserFromCtx
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.Logging, interceptor.UnaryAuthentication),
	)

	transactionProto.RegisterTransactionServiceServer(grpcServer, &c.TransactionService{
		GetUser:         getUser,
		TransactionRepo: transactionRepo,
	})

	log.Printf("Starting to serve in port : %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %s", err.Error())
	}
}
