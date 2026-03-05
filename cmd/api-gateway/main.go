package main

import (
	rateLimiter "chirpstream/internal/ratelimit"
	"chirpstream/pkg/config"
	"chirpstream/pkg/response"
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
			response.Error(w, "Unauthorized", http.StatusUnauthorized)
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
				response.Error(w, "Too Many Requests", http.StatusTooManyRequests)
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

	userURL, err := url.Parse(cfg.User_Service.Addr)
	if err != nil || userURL.Host == "" {
		log.Fatalf("Invalid user service address: %v", cfg.User_Service.Addr)
	}
	userProxy := httputil.NewSingleHostReverseProxy(userURL)

	chirpURL, err := url.Parse(cfg.Chirp_Service.Addr)
	if err != nil || chirpURL.Host == "" {
		log.Fatalf("Invalid chirp service address: %v", cfg.Chirp_Service.Addr)
	}
	chirpProxy := httputil.NewSingleHostReverseProxy(chirpURL)

	relUrl, err := url.Parse(cfg.Rel_Service.Addr)
	if err != nil || chirpURL.Host == "" {
		log.Fatalf("Invalid chirp service address: %v", cfg.Chirp_Service.Addr)
	}

	relProxy := httputil.NewSingleHostReverseProxy(relUrl)

	// Set timeouts for the reverse proxies to prevent hanging requests
	userProxyWithTimeout := http.TimeoutHandler(userProxy, 5*time.Second, "Gateway Timeout")
	chirpProxyWithTimeout := http.TimeoutHandler(chirpProxy, 5*time.Second, "Gateway Timeout")
	relProxyWithTimeout := http.TimeoutHandler(relProxy, 5*time.Second, "Gateway Timeout")

	//not using httprouter here as there is no need for a complex router,
	// we just need to forward requests to the appropriate services based on the path
	mux := http.NewServeMux()
	mux.Handle("POST /api/user/login", http.StripPrefix("", userProxyWithTimeout))
	mux.Handle("POST /api/user/register", http.StripPrefix("", userProxyWithTimeout))

	mux.Handle("/api/users/", AuthMiddleware(userProxyWithTimeout))
	mux.Handle("/api/chirp/", AuthMiddleware(chirpProxyWithTimeout))
	mux.Handle("/api/user/", AuthMiddleware(relProxyWithTimeout))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, "Not Found", http.StatusNotFound)
	}))

	handler := rateLimitMiddleware(rlc)(mux)

	log.Println("API Gateway started on address", cfg.API_Gateway.Addr)
	err = http.ListenAndServe(cfg.API_Gateway.Addr, handler)
	if err != nil {
		log.Fatalf("Error starting API Gateway: %v\n", err)
	}
}
