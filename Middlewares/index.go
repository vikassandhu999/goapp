package Middlewares


import "net/http"
import "github.com/dgrijalva/jwt-go"



var jwtKey = []byte("my_secret_key")



type Middleware func(http.HandlerFunc) http.HandlerFunc


func Chain(f http.HandlerFunc, middlewares ... Middleware) http.HandlerFunc {
    for _, m := range middlewares {
        f = m(f)
    }
    return f
}


func AuthMiddleware() Middleware {
	return  func(f http.HandlerFunc) http.HandlerFunc {
	
		return	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				tokenString := r.Header.Get("x-access-token")
				
					claims := &jwt.StandardClaims{}
					
					tkn, err := jwt.ParseWithClaims(tokenString, claims , func(token *jwt.Token) (interface{}, error) {
						return jwtKey, nil
					})

					if err != nil {
						if err == jwt.ErrSignatureInvalid {
							w.WriteHeader(http.StatusUnauthorized)
							return
						}
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if !tkn.Valid {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

				r.Header.Set("user", claims.Subject)
				r.Header.Set("auth", "true")

				f(w,r)
			})
    }
}
