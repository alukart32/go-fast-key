package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alukart32/go-fast-key/internal/fastkey"
	"github.com/alukart32/go-fast-key/internal/fastkey/compute"
	"github.com/alukart32/go-fast-key/internal/fastkey/storage"
)

func main() {
	storage := storage.NewEngine(256)
	requestParser := compute.NewParser()

	db := fastkey.NewDatabase(requestParser, storage)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("[fastkey] > ")
		request, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[error] fail to read request")
		}

		result, err := db.HandleRequest(request)
		if err != nil {
			fmt.Printf("[error] %v\n", err)
		} else {
			fmt.Println(result)
		}
	}
}
