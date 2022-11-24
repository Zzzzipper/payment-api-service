package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	pbPayment "payment-api-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Backend implements the protobuf interface
type Backend struct {
	mu *sync.RWMutex
}

// New initializes a new Backend struct.
func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

// Examples
/*
// AddUser adds a user to the in-memory store.
func (b *Backend) AddUser(ctx context.Context, _ *pbPayment.AddUserRequest) (*pbPayment.User, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	user := &pbPayment.User{
		Id: uuid.Must(uuid.NewV4()).String(),
	}
	b.users = append(b.users, user)

	return user, nil
}

// ListUsers lists all users in the store.
func (b *Backend) ListUsers(_ *pbPayment.ListUsersRequest, srv pbPayment.PaymentService_ListUsersServer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, user := range b.users {
		err := srv.Send(user)
		if err != nil {
			return err
		}
	}

	return nil
}
*/

func (b *Backend) Block(ctx context.Context, in *pbPayment.BlockRequest) (*pbPayment.BlockHandler, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(md)
	}

	body, err := json.Marshal(in)
	if err == nil {
		fmt.Println(string(body))
	} else {
		fmt.Printf("can't marshall block request: %s", err.Error())
		fmt.Println("BlockRequest", in)
	}

	var block *pbPayment.BlockHandler

	requestURL := os.Getenv("MAPAPI_ADDR")
	if requestURL != "" {
		bodyReader := bytes.NewReader(body)

		req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{
			Timeout: 30 * time.Second,
		}

		var res *http.Response
		res, err = client.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			block = &pbPayment.BlockHandler{
				Success:    false,
				ErrCode:    "COMMUNICATE_ERROR",
				ErrMessage: "failed connection to MAP API",
			}
		} else {
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("client: could not read response body: %s\n", err)
				block = &pbPayment.BlockHandler{
					Success:    false,
					ErrCode:    "EMPTY_RESPONSE",
					ErrMessage: "processing error",
				}
			}
			fmt.Printf("client: response body: %s\n", resBody)

			block = &pbPayment.BlockHandler{}

			err = json.Unmarshal(resBody, block)
			if err != nil {
				fmt.Printf("client: could not parse response body: %s\n", err)
				block = &pbPayment.BlockHandler{
					Success:    false,
					ErrCode:    "ENCODING_ERROR",
					ErrMessage: "encoding response error",
				}
			}
		}
	} else {
		block = &pbPayment.BlockHandler{
			Success:    false,
			ErrCode:    "NOT_FOUND",
			ErrMessage: "address MAP API not present",
		}
	}

	return block, nil
}

func (b *Backend) Charge(ctx context.Context, in *pbPayment.ChargeRequest) (*pbPayment.ChargeHandler, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	charge := &pbPayment.ChargeHandler{
		Id: "charge_handler",
	}

	return charge, nil
}

func (b *Backend) Get(ctx context.Context, in *pbPayment.Order) (*pbPayment.OrderStatus, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	status := &pbPayment.OrderStatus{
		Status: "status ok",
	}

	return status, nil
}

func (b *Backend) Init(ctx context.Context, in *pbPayment.OrderRequest) (*pbPayment.Order, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	order := &pbPayment.Order{
		Id: "order ",
	}

	return order, nil
}

func (b *Backend) Pay(ctx context.Context, in *pbPayment.PayRequest) (*pbPayment.Payment, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	payment := &pbPayment.Payment{
		Id: "block_handler",
	}

	return payment, nil
}

func (b *Backend) Payout(ctx context.Context, in *pbPayment.PayoutRequest) (*pbPayment.PayoutHandler, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	payout := &pbPayment.PayoutHandler{
		Id: "block_handler",
	}

	return payout, nil
}

func (b *Backend) Refund(ctx context.Context, in *pbPayment.Order) (*pbPayment.RefundHandler, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	refund := &pbPayment.RefundHandler{
		Id: "block_handler",
	}

	return refund, nil
}

func (b *Backend) Void(ctx context.Context, in *pbPayment.Order) (*pbPayment.VoidHandler, error) {
	_, ok := GetUserMetadata(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user authentication found")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	void := &pbPayment.VoidHandler{
		Id: "block_handler",
	}

	return void, nil
}
