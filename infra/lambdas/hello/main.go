package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) error {
	fmt.Println("Hello, World!")
	return nil
}

func main() {
	lambda.Start(handler)
}
