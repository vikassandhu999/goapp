package Helpers


import "net/http"
import "encoding/json"
import "golang.org/x/crypto/bcrypt"


func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err
}


func Error ( w http.ResponseWriter , message string , code int ) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct{  Error string }{ Error: message })
	return
}


func Message (w http.ResponseWriter , message string ,  code int ) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct{  Message string }{ Message : message })
	return
}


