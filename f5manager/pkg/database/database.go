package database

// import (
// 	"context"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func NewDB(uri, name string, timeout time.Duration, num uint64) (*mongo.Database, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	o := options.Client().ApplyURI(uri)
// 	o.SetMaxPoolSize(num)
// 	client, err := mongo.Connect(ctx, o)
// 	if err != nil {
// 		return nil, err
//   }
//   client.

// 	return client.Database(name), nil
// }
