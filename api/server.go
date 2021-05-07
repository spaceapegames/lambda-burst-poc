package api

import (
	"fmt"
	"github.com/apex/gateway"
	"github.com/gorilla/mux"
	"github.com/spaceapegames/lambda-burst/primes"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	Router          *mux.Router
	lambdaMode      bool
	port            int64
	redirectAddress string
	disableBurst    bool
	maxPrime        int

	rateLimiter *rate.Limiter
}

func NewServer(lambdaMode bool, port int64, redirectAddress string, rateLimit int, disableBurst bool) Server {
	maxPrime, err := strconv.Atoi(os.Getenv("MAX_PRIME"))
	if err != nil {
		log.Printf("invalid value for MAX_PRIME, defaulting to 1000")
		maxPrime = 1000
	}

	s := Server{
		Router:          mux.NewRouter(),
		lambdaMode:      lambdaMode,
		port:            port,
		redirectAddress: redirectAddress,
		rateLimiter:     rate.NewLimiter(rate.Limit(rateLimit), 1),
		disableBurst:    disableBurst,
		maxPrime:        maxPrime,
	}
	s.Routes()
	return s
}

func (s *Server) Serve() {
	if s.lambdaMode {
		log.Println("running in Lambda mode")
		gateway.ListenAndServe(":8000", s.Router)
	} else {
		log.Printf("running in http mode on port %d", s.port)
		http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.Router)
	}
}

func (s *Server) Routes() {
	s.Router.Use(s.middleware)
	s.Router.HandleFunc("/doThing", s.DoThing()).Methods(http.MethodGet)
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldRedirect := false
		if !s.lambdaMode && !s.disableBurst && s.redirectAddress != "" {
			shouldRedirect = !s.rateLimiter.Allow() //Allow() is protected by a mutex
		}

		if shouldRedirect {
			log.Println("redirecting")
			w.Header().Add("Location", s.redirectAddress+"/doThing")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s Server) DoThing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// calculate all the prime numbers between 0 and 1000000
		primes.Calculate(s.maxPrime)
		msg := "Hello from Fargate!"
		if s.lambdaMode {
			msg = "Hello from Lambda!"
		}
		w.Write([]byte(msg))
	}
}
