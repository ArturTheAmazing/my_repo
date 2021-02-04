package main

import (
	pb "Projects/SberAuto/service/validator"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func Validate (client pb.ValidatorClient, data string) (bool, error) {
	ctx := context.Background()

	req := pb.ValidationRequest{
		Input: data,
	}

	resp, err := client.Validate(ctx, &req)
	if err != nil {
		return false, err
	}
	return resp.Resp, nil
}

func Fix(client pb.ValidatorClient, data string) (string, error) {
	ctx := context.Background()

	req := pb.ValidationRequest{
		Input: data,
	}

	resp, err := client.Fix(ctx, &req)
	if err != nil {
		return "", err
	}
	return resp.Output, nil
}

func main(){
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewValidatorClient(conn)

	fixed, err := Fix(client, "[(])")
	if err != nil {
		log.Fatalf("failed to fix: %v", err)
	}
	fmt.Println(fixed)
}