package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	gorillaLambda *gorillamux.GorillaMuxAdapter
	log           *logrus.Logger
)

func init() {
	// Initialize logger with JSON format for structured logging.
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	
	// Enable debug level logging if required.
	// if os.Getenv("DEBUG") == "true" {
		log.SetLevel(logrus.DebugLevel)
	// }

	log.Info("Initializing Go service")
	
	// Setup router and create the API Gateway proxy adapter.
	gorillaLambda = gorillamux.New(setupRouter())
}

// healthCheck is an HTTP handler that returns service status and current time.
func healthCheck(w http.ResponseWriter, r *http.Request) {
	// Capture current time for response.
	now := time.Now()
	response := map[string]string{
		"status": "healthy",
		"time":   now.String(),
	}

	// Encode response as JSON and handle errors.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithError(err).Error("Failed to encode health check response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	// Log health check details.
	log.WithFields(logrus.Fields{
		"status": "healthy",
		"time":   now,
	}).Info("Health check completed")
}

// setupRouter configures the HTTP routes and middleware.
func setupRouter() *mux.Router {
	r := mux.NewRouter()
	
	// Middleware to log incoming request details.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.WithFields(logrus.Fields{
				"original_path": req.URL.Path,
				"method":       req.Method,
			}).Debug("Incoming request")
			next.ServeHTTP(w, req)
		})
	})

	// Health check endpoints mapped on both root and /health.
	healthHandler := http.HandlerFunc(healthCheck)
	r.Handle("/", healthHandler).Methods("GET")
	r.Handle("/health", healthHandler).Methods("GET")
	
	return r
}

// Handler processes API Gateway requests using Gorilla Mux as backend.
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.WithFields(logrus.Fields{
		"path":   req.Path,
		"method": req.HTTPMethod,
		"source": "lambda",
	}).Info("Handling Lambda request")

	// Handle health check paths directly.
	if req.Path == "/" || req.Path == "/health" {
		response := map[string]string{
			"status": "healthy",
			"time":   time.Now().String(),
		}
		jsonResponse, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonResponse),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Use Gorilla Mux adapter to proxy other requests.
	r, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))
	if err != nil {
		log.WithError(err).Error("Failed to proxy request")
	}
	return *r.Version1(), err
}

// main is the entry point of the service. It determines the runtime environment.
func main() {
	log.Info("Starting Go service...")

	// Run in Lambda mode if AWS Lambda environment variable is set.
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Info("Running in Lambda mode")
		lambda.Start(Handler)
	} else {
		// Run in local server mode.
		log.Info("Running in local server mode")
		r := setupRouter()
		
		srv := &http.Server{
			Handler:      r,
			Addr:         ":8080",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		
		log.WithField("port", 8080).Info("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil {
			log.WithError(err).Fatal("Server failed to start")
		}
	}
}
