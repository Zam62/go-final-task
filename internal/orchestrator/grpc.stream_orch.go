package orchestrator

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"

	pb "go-final-task/api/gen/go"
	"go-final-task/pkg/models"

	"google.golang.org/grpc"
)

const (
	tcp         = "tcp"
	addr string = ":5000"
)

type Task struct {
	ID   int32
	Arg1 string
	Arg2 string
	Type string
}

var (
	resultsCh = make(chan *models.TaskResult)
	tasksCh   = make(chan *models.Task)
)

type Server struct {
	pb.UnimplementedOrchestratorServer
	mu sync.Mutex
}

func NewServer() *Server {
	return &Server{mu: sync.Mutex{}}
}

func (s *Server) Calculate(stream pb.Orchestrator_CalculateServer) error {
	log.Println("agent connected to gRPC server")
	defer log.Println("agent disconnected")
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	done := make(chan struct{})
	defer close(done)

	go func() {
		defer cancel()

		for {
			select {
			case task := <-tasksCh:
				s.mu.Lock()
				i64, _ := strconv.ParseInt(task.ID, 10, 32)

				err := stream.Send(&pb.TaskRequest{
					Id:       int32(i64),
					Arg1:     strconv.FormatFloat(task.Arg1, 'f', -1, 64),
					Arg2:     strconv.FormatFloat(task.Arg2, 'f', -1, 64),
					Operator: task.Operation,
				})
				s.mu.Unlock()

				if err != nil {
					log.Printf("Failed to send task: %v", err)
					return
				}
			case <-ctx.Done():
				return
			case <-done:
				return
			}
		}
	}()

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := stream.Recv()
				if err != nil {
					log.Printf("Receive error: %v", err)
					return
				}
				resultsCh <- &models.TaskResult{
					ID:     string(res.Id),
					Result: float64(res.Result),
					Error:  res.Error,
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func runGRPC() {
	log.Println("Starting tcp server...")
	lis, err := net.Listen(tcp, addr)
	if err != nil {
		log.Fatalf("error starting tcp server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcServer, NewServer())

	log.Printf("tcp server started at: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error serving grpc: %v", err)
	}
}
