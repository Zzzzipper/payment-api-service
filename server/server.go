package server

import (
	"context"
	"sync"

	"github.com/gofrs/uuid"

	pbPayment "payment-api-service/proto"
)

// Backend implements the protobuf interface
type Backend struct {
	mu    *sync.RWMutex
	users []*pbPayment.User
}

// New initializes a new Backend struct.
func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

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

func (b *Backend) Block(in *pbPayment.BlockRequest, out pbPayment.PaymentService_BlockServer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	block := pbPayment.BlockHandler{Id: "test"}

	err := out.Send(&block)
	if err != nil {
		return err
	}

	return nil
}
