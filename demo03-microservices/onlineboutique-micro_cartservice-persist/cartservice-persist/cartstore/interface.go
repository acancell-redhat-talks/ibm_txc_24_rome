package cartstore

import (
	"context"

	pb "github.com/go-micro/demo/cartservice/proto"

	"github.com/valkey-io/valkey-go"
)

type CartStore interface {
	AddItem(ctx context.Context, userID, productID string, quantity int32) error
	EmptyCart(ctx context.Context, userID string) error
	GetCart(ctx context.Context, userID string) (*pb.Cart, error)
}

func NewMemoryCartStore(cdb valkey.Client, crt map[string]map[string]int32) CartStore {
	return &memoryCartStore{
		client_db: cdb,
		carts: crt,
	}
}
