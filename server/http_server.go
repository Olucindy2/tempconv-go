package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "tempconv-go/proto"

	"google.golang.org/grpc"
)

type TempResponse struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
}

func startHTTPServer(grpcAddr string, httpPort string) {
	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {

		// ⭐ CORS HEADERS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// ⭐ Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		// ⭐ Read query parameter
		cStr := r.URL.Query().Get("celsius")
		if cStr == "" {
			http.Error(w, "missing celsius parameter", http.StatusBadRequest)
			return
		}

		cVal, err := strconv.ParseFloat(cStr, 64)
		if err != nil {
			http.Error(w, "invalid celsius value", http.StatusBadRequest)
			return
		}

		// ⭐ Connect to gRPC server using provided address
		conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
		if err != nil {
			http.Error(w, "could not connect to gRPC server", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := pb.NewTempConvServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// ⭐ Call gRPC method
		resp, err := client.CelsiusToFahrenheit(ctx, &pb.TempRequest{Value: cVal})
		if err != nil {
			http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// ⭐ Prepare JSON response
		out := TempResponse{
			Celsius:    cVal,
			Fahrenheit: resp.Value,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	})

	log.Printf("HTTP server listening on %s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
