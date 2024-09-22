package main

import (
	"context"
	"sync"

	"strconv"

	"github.com/valkey-io/valkey-go"
	"fmt"
	//"reflect"
)

type memoryCartStore struct {
	sync.RWMutex

	client_db *valkey.Client

	carts map[string]map[string]int32
}

func NewMemoryCartStore(cdb *valkey.Client, initmap map[string]map[string]int32) *memoryCartStore {
	return &memoryCartStore{
		client_db: cdb,
		carts: initmap, 
	}
}

func (s *memoryCartStore) AddItem(ctx context.Context, userID, productID string, quantity int32) error {
	s.Lock()
	defer s.Unlock()

	conn := *(s.client_db)

	if cart, ok := s.carts[userID]; ok {
		if currentQuantity, ok := cart[productID]; ok {
			cart[productID] = currentQuantity + quantity
		} else {
			cart[productID] = quantity
		}
		s.carts[userID] = cart
	} else {
		s.carts[userID] = map[string]int32{productID: quantity}
	}

	for user, cart := range s.carts {
		for item, quantity := range cart {
				err := conn.Do(ctx, conn.B().Arbitrary("HSET").Keys(user).Args(item, strconv.Itoa(int(quantity))).Build()).Error()
				if err != nil {
					//logger.Fatal(err.Error())
					return err
				}
		}
	}

	return nil
}

func (s *memoryCartStore) EmptyCart(ctx context.Context, userID string) error {
	s.Lock()
	defer s.Unlock()

	conn := *(s.client_db)

	delete(s.carts, userID)

	err := conn.Do(ctx, conn.B().Arbitrary("DEL").Keys(userID).Build()).Error()
	if err != nil {
		//logger.Fatal(err.Error())
		return err
	}

	return nil
}

// MAIN

// Connect to db
func connectDB(ctx context.Context) valkey.Client {
	client_db, err_db := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err_db != nil {
		panic(err_db)
	}
	//defer client_db.Close()
    
	return client_db
}

// Load from db to memorymap
func initMap(ctx context.Context, client_db *valkey.Client) map[string]map[string]int32 {
    conn := *(client_db)
	initmap := make(map[string]map[string]int32)
	
    scan, err := conn.Do(ctx, conn.B().Scan().Cursor(0).Count(9999).Build()).AsScanEntry()
	if valkey.IsParseErr(err) {
		//logger.Fatal(err)
	}
	initkeys := scan.Elements
	for _, key := range initkeys {
		initmap[key] = map[string]int32{}
	
		m, err := conn.Do(ctx, conn.B().Hgetall().Key(key).Build()).AsStrMap()
		if valkey.IsParseErr(err) {
			//logger.Fatal(err)
		}
		for k, v := range m {
			initvalue, _ := strconv.Atoi(v)
			initmap[key][k] = int32(initvalue)
		}
	}
	return initmap
}

func main(){

	//fmt.Printf("%+v\n", conn)
	fmt.Printf("START\n")

	ctx := context.Background()

	var conn valkey.Client = connectDB(ctx) 
	var initmap map[string]map[string]int32 = initMap(ctx, &conn)
	var cartStore = NewMemoryCartStore(&conn, initmap) 

	//TESTS
	cartStore.AddItem(ctx, "john", "thing1", 1) 
	cartStore.AddItem(ctx, "jane", "thing2", 2) 
	cartStore.EmptyCart(ctx, "john") 
}

