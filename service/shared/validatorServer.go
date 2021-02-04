package shared

import (
	pb "Projects/SberAuto/service/validator"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
	"sync"
)

var (
	Reg = prometheus.NewRegistry()

	//reqCountMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
	//	Name:"requests_count",
	//	Help: "Total number of requests",
	//}, []string{"name"})

	reqDuartionMetric = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "request_duraion_seconds",
		Help: "Histogram for the runtime of a method",
		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
	})
)

func init(){
	Reg.MustRegister(reqDuartionMetric)
	//reqCountMetric.WithLabelValues("requests")
}

type ValidatorServer struct {
	pb.UnimplementedValidatorServer
	mu sync.Mutex
}

func NewServer() *ValidatorServer {
	return &ValidatorServer{}
}

func (s *ValidatorServer) Validate(_ context.Context, req *pb.ValidationRequest) (*pb.ValidationResponse, error) {
	//count, _ := reqCountMetric.GetMetricWithLabelValues("requests")
	//count.Inc()

	timer:=prometheus.NewTimer(reqDuartionMetric)
	defer timer.ObserveDuration()

	log.Printf("Validate: got input %s", req.Input)

	err := checkOrder(req.Input)

	if err != nil {
		log.Print(err)
		return &pb.ValidationResponse{
			Resp: false,
		}, err
	}

	return &pb.ValidationResponse{
		Resp: true,
	}, nil
}

func checkOrder(input string) error {
	if len(input)%2 != 0 {
		return fmt.Errorf("Not all have pairs: %v", input)
	}
	if len(input) == 0 {
		return nil
	}

	var opposite string

	first := string(input[0])
	switch first {
	case "[":
		opposite = "]"
	case "(":
		opposite = ")"
	case "{":
		opposite = "}"
	default:
		return fmt.Errorf("First character is a closing bracket: %s", string(input[0]))
	}

	if !strings.Contains(input, opposite) ||
		len(input) == 2 && !strings.Contains(input, opposite) {
		return fmt.Errorf("Input does not conain the opposite: %s", input)
	}

	counter := 0
	in := 0
	found := false

	if string(input[1]) == opposite {
		in = 1
	} else {
		for index, ch := range input {
			if string(ch) == first {
				found = true
				counter++
			} else if string(ch) == opposite {
				counter--
			}

			if (counter == 0 && found ||
				!found && counter == -1) && string(ch) == opposite {
				in = index
				break
			}
		}
	}
	if len(input[1:in])%2 != 0 {
		return fmt.Errorf("Resulting string does not have enough matches: %s", input[1:in])
	}
	input = input[1:in] + input[in+1:]

	return checkOrder(input)
}

func (s *ValidatorServer) Fix(ctx context.Context, req *pb.ValidationRequest) (*pb.FixResponse, error) {
	log.Printf("Fix: got input %s", req.Input)

	if _, err := s.Validate(ctx, req); err == nil {
		log.Printf("Nothing to fix")
		return &pb.FixResponse{Output: req.Input}, nil
	}

	input := req.Input

	brackets := map[string]string{
		"[": "]",
		"{": "}",
		"(": ")",
	}

	var res string
	for _, ch := range input {
		for k, v := range brackets {
			if string(ch) == k || string(ch) == v {
				res += k + v
			}
		}
	}

	return &pb.FixResponse{Output: res}, nil
}
