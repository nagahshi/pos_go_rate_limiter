package middleware

import (
	"net/http"
	"strings"

	"github.com/nagahshi/pos_go_rate_limiter/internal/usecase"
)

type middleware struct {
	limiter *usecase.Limiter
}

func NewMiddleware(limiter *usecase.Limiter) *middleware {
	return &middleware{
		limiter: limiter,
	}
}

func (l *middleware) Run(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("API_KEY"); token != "" && l.limiter.GetMaxTokenRequests() > 0 {
			blocked, err := l.limiter.AllowToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !blocked {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}

		addrSplited := strings.Split(r.RemoteAddr, ":")
		if IP := addrSplited[0]; IP != "" && l.limiter.GetMaxIPRequests() > 0 {
			blocked, err := l.limiter.AllowIP(r.Context(), IP)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !blocked {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}

		f(w, r)
	}
}
