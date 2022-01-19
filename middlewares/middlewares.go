package middlewares

import (
	"fmt"
	"net/http"
)

func Use(next http.HandlerFunc, mfs ...middlewareFunc) http.HandlerFunc {
	var hf = next
	for i := len(mfs) - 1; i >= 0; i-- {
		hf = mfs[i](hf)
	}
	return hf
}

type middlewareFunc func(handler http.HandlerFunc) http.HandlerFunc

func ValidateJWT(handler http.HandlerFunc) http.HandlerFunc {
	fmt.Println("1~~~")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware ValidateJWT req")
		token := r.Header.Get("Authorization")
		fmt.Println(token)
		handler.ServeHTTP(w, r)

		fmt.Println("middleware 1 writer")
	})
}

// func Middleware2(handler http.HandlerFunc) http.HandlerFunc {
// 	fmt.Println("2~~~")
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("middleware 2 req")

// 		r.Header.Set("BCD", "234")
// 		fmt.Println(r.Header.Get("BCD"))
// 		handler.ServeHTTP(w, r)

// 		fmt.Println("middleware 2 writer")

// 	})
// }
