package main

import (
	"encoding/json"
	"fmt"
	proto "grpc-3/pb"
	"io"
	"log"
	"net/http"
)


func CallHttpRequest() ([]*proto.ProductResponse, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("res: %s\n", string(body)) // Print the body for debugging

	var products []*proto.ProductResponse
	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, err
	}

	// response := &proto.ProductList{Products: products}
	return products, nil
}

func testHTTP() {
	res, err := CallHttpRequest()
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range res {
		fmt.Printf("id : %d \t title :%s \n", result.Id, result.Title)
	}

	return

}