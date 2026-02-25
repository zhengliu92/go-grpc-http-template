package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"go-grpc-http-template/internal/util"
)

// JwtAuth 返回 JWT 验证中间件。
// publicPaths: 不需要鉴权的路径（如 /api/v1/health、/api/v1/login）
func JwtAuth(secret string, publicPaths ...string) func(http.HandlerFunc) http.HandlerFunc {
	public := make(map[string]struct{}, len(publicPaths))
	for _, p := range publicPaths {
		public[p] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 公开路径直接放行
			if _, ok := public[r.URL.Path]; ok {
				next(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := util.ParseToken(secret, parts[1])
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			// 将用户身份注入 gRPC metadata（gateway 会转发 Grpc-Metadata-* 前缀的 header）
			r.Header.Set("Grpc-Metadata-User-Id", strconv.FormatInt(claims.UserId, 10))
			r.Header.Set("Grpc-Metadata-Username", claims.Username)

			next(w, r)
		}
	}
}
