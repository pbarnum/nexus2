package middleware

import (
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/msrevive/nexus2/internal/log"
	"github.com/msrevive/nexus2/internal/rate"
	"github.com/msrevive/nexus2/internal/system"
)

var (
	globalLimiter *rate.Limiter
)

func getIP(r *http.Request) string {
	ip := r.Header.Get("X_Real_IP")
	if ip == "" {
		ips := strings.Split(r.Header.Get("X_Forwarded_For"), ",")
		if ips[0] != "" {
			return strings.TrimSpace(ips[0])
		}

		// Ignoring error and defaulting to zero-value
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return ip
}

func setControlHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	// Maximum age allowable under Chromium v76 is 2 hours, so just use that since
	// anything higher will be ignored (even if other browsers do allow higher values).
	//
	// @see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age#Directives
	w.Header().Set("Access-Control-Max-Age", "7200")
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setControlHeaders(w) // best place to set control headers?
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Log.Printf("%s %s from %s (%v)", r.Method, r.RequestURI, getIP(r), time.Since(start))
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if panic := recover(); panic != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Log.Errorln("500: We have encountered an error with the last request.")
				log.Log.Errorf("500: Error: %s", panic.(error).Error())
				log.Log.Errorf(string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RateLimit(cfg *system.ApiConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if globalLimiter == nil {
				globalLimiter = rate.NewLimiter(1, cfg.RateLimit.MaxRequests, cfg.RateLimit.MaxAge, 0)
			}

			globalLimiter.CheckTime()
			if globalLimiter.IsAllowed() == false {
				log.Log.Println("Received too many requests.")
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Auth(cfg system.AuthConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			// IP Auth
			if cfg.IsEnforcingIP() {
				if !cfg.IsKnownIP(ip) {
					log.Log.Printf("%s Is not authorized.", ip)
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
			}

			// API Key Auth
			if cfg.IsEnforcingKey() {
				if r.Header.Get("Authorization") != cfg.GetKey() {
					log.Log.Printf("%s failed API key check.", ip)
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func NoAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		return
	}
}
