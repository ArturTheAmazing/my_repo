package main

import (
	"Projects/SberAuto/service/shared"
	pb "Projects/SberAuto/service/validator"
	"bufio"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"strings"
	"testing"
)

var validData = `[()]{}{[()()]()}`

var invalidData = map[string]string{
	"enclosedBracket": "[(])",
	"enclosedBracketFixed": "[]()[]()",
	"openedBrackets": "[({",
	"openedBracketsFixed": "[](){}",
	"closedBrackets": ")}]",
	"closedBracketsFixed": "(){}[]",
	"repeatingBrackets": "[()]{}{[()()]()}}}}}}}}",
	"repeatingBracketsFixed": "[]()()[]{}{}{}[]()()()()[]()(){}{}{}{}{}{}{}{}",
}

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterValidatorServer(s, shared.NewServer())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Error while serving: %v", err)
		}
	}()
}

func TestValid(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error while dialing: %v", err)
	}
	defer conn.Close()

	client := pb.NewValidatorClient(conn)

	scanner := bufio.NewScanner(strings.NewReader(validData))
	for scanner.Scan() {
		d := scanner.Text()
		log.Printf("Validating %v", d)

		resp, err := client.Validate(ctx, &pb.ValidationRequest{Input: d})
		if err != nil {
			t.Fatalf("Error while testing valid: %v", err)
		} else if !resp.Resp {
			t.Fatalf("Testing valid. expected: true, got: %v", resp.Resp)
		}
	}
}

func TestInvalid(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error while dialing: %v", err)
	}
	defer conn.Close()

	client := pb.NewValidatorClient(conn)

		d := invalidData["enclosedBracket"]
		log.Printf("Validating %v", d)

		resp, err := client.Validate(ctx, &pb.ValidationRequest{Input: d})
		if err != nil {
			t.Skip()
		} else if resp.Resp {
			t.Fatalf("Testing invalid. Expected: false, got: %v", resp.Resp)
		}
}

func TestFixNoChange(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error while dialing: %v", err)
	}
	defer conn.Close()

	client := pb.NewValidatorClient(conn)
		log.Printf("Fixing %s", validData)

		resp, err := client.Fix(ctx, &pb.ValidationRequest{Input: validData})
		if err!=nil{
			t.Fatalf("Error while fixing: %s", err)
		} else if resp.Output != validData{
			t.Fatalf("Not equal. Expected: %s, got: %s", validData, resp.Output)
		}
		log.Printf("Fixed: %s", resp.Output)
}

func TestFixEnclosedBracket(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error while dialing: %v", err)
	}
	defer conn.Close()

	client := pb.NewValidatorClient(conn)

		d := invalidData["enclosedBracket"]
		log.Printf("Fixing %s", d)

		resp, err := client.Fix(ctx, &pb.ValidationRequest{Input: d})
		if err!=nil{
			t.Fatalf("Error while fixing: %s", err)
		} else if resp.Output != invalidData["enclosedBracketFixed"] {
			t.Fatalf("Not equal. Expected: %s, got: %s", d, resp.Output)
		}
		log.Printf("Fixed: %s", resp.Output)
}

func TestFixOpenedBracket(t *testing.T){
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error while dialing: %v", err)
	}
	defer conn.Close()

	client := pb.NewValidatorClient(conn)

	d := invalidData["openedBrackets"]
	log.Printf("Fixing %s", d)

	resp, err := client.Fix(ctx, &pb.ValidationRequest{Input: d})
	if err!=nil{
		t.Fatalf("Error while fixing: %s", err)
	} else if resp.Output != invalidData["openedBracketsFixed"] {
		t.Fatalf("Not equal. Expected: %s, got: %s", d, resp.Output)
	}
	log.Printf("Fixed: %s", resp.Output)
}

func TestFixRepeatingBracket(t *testing.T){
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error while dialing: %v", err)
	}
	defer conn.Close()

	client := pb.NewValidatorClient(conn)

	d := invalidData["repeatingBrackets"]
	log.Printf("Fixing %s", d)

	resp, err := client.Fix(ctx, &pb.ValidationRequest{Input: d})
	if err!=nil{
		t.Fatalf("Error while fixing: %s", err)
	} else if resp.Output != invalidData["repeatingBracketsFixed"] {
		t.Fatalf("Not equal. Expected: %s, got: %s", d, resp.Output)
	}
	log.Printf("Fixed: %s", resp.Output)
}
