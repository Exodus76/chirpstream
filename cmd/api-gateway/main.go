package main

import (
	rateLimiter "chirpstream/internal/ratelimit"
	"chirpstream/pkg/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//TODO: Validate JWT token and extract user ID

		log.Println("Auth Middleware: Request Authorized")
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(rl *rateLimiter.RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAllowed := rateLimiter.Validate(rl, r.RemoteAddr)

			if !isAllowed {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	cfg, err := config.Init(".")
	if err != nil {
		log.Fatalf("Cant parse config file %v\n", err)
	}

	rlc := rateLimiter.NewRateLimterConfig(10, 5)

	userURL, _ := url.Parse(cfg.User_Service.Addr)
	userProxy := httputil.NewSingleHostReverseProxy(userURL)
	chirpURL, _ := url.Parse(cfg.Chirp_Service.Addr)
	chirpProxy := httputil.NewSingleHostReverseProxy(chirpURL)

	// Set timeouts for the reverse proxies to prevent hanging requests
	userProxyWithTimeout := http.TimeoutHandler(userProxy, 5*time.Second, "Gateway Timeout")
	chirpProxyWithTimeout := http.TimeoutHandler(chirpProxy, 5*time.Second, "Gateway Timeout")

	//not using httprouter here as there is no need for a complex router,
	// we just need to forward requests to the appropriate services based on the path
	mux := http.NewServeMux()
	mux.Handle("POST /api/user/login", http.StripPrefix("", userProxyWithTimeout))
	mux.Handle("POST /api/user/register", http.StripPrefix("", userProxyWithTimeout))

	mux.Handle("/api/users/", AuthMiddleware(userProxyWithTimeout))
	mux.Handle("/api/chirps/", AuthMiddleware(chirpProxyWithTimeout))

	http.Handle("/", rateLimitMiddleware(rlc)(mux))

	log.Println("API Gateway started on address", cfg.API_Gateway.Addr)
	err = http.ListenAndServe(cfg.API_Gateway.Addr, nil)
	if err != nil {
		log.Fatalf("Error starting API Gateway: %v\n", err)
	}
}
