package imdb_test

import (
	"context"
	"fmt"

	"github.com/kenshaw/imdb"
)

func ExampleNew_findTitle() {
	cl := imdb.New()
	res, err := cl.FindTitle(context.Background(), "luca")
	if err != nil {
		panic(err)
	}
	if len(res) == 0 {
		fmt.Println("expected at least one result")
		return
	}
	fmt.Printf("result: %s\n", res[0])
	// Output:
	// result: tt12801262: "Luca" (2021) https://www.imdb.com/title/tt12801262/
}
