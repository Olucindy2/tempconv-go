package main

import (
    "context"
    "fmt"
    "log"
    "time"

    pb "tempconv-go/proto"

    "google.golang.org/grpc"
)

func main() {
    // Connect to the gRPC server
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewTempConvServiceClient(conn)

    // Call CelsiusToFahrenheit
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    celsius := 25.0
    resp1, err := client.CelsiusToFahrenheit(ctx, &pb.TempRequest{Value: celsius})
    if err != nil {
        log.Fatalf("error calling CelsiusToFahrenheit: %v", err)
    }
    fmt.Printf("%.2f°C = %.2f°F\n", celsius, resp1.Value)

    // Call FahrenheitToCelsius
    fahrenheit := 77.0
    resp2, err := client.FahrenheitToCelsius(ctx, &pb.TempRequest{Value: fahrenheit})
    if err != nil {
        log.Fatalf("error calling FahrenheitToCelsius: %v", err)
    }
    fmt.Printf("%.2f°F = %.2f°C\n", fahrenheit, resp2.Value)
}
