package main

import (
	"context"
	"fmt"
	"shorturl/link"
	"shorturl/shint"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	sh := shint.New(128, "link:shint", client)
	link := link.New(sh)
	fmt.Println(link.NewID(context.Background()))
}
