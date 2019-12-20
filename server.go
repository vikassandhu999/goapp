package main

import "net/http"
import "log"
import 	"github.com/gorilla/mux"
import "cloud.google.com/go/firestore"
import "./Db"
import "./Middlewares"
import "./Controllers"


var client *firestore.Client


func main(){
	r := mux.NewRouter()

	Db.InitDb()
	client = Db.Client
	defer client.Close()

	apiRouter := r.PathPrefix("/api").Subrouter()
	
	apiRouter.HandleFunc("/user/signup", Controllers.SignupController).Methods("POST")
	apiRouter.HandleFunc("/user/signin", Controllers.SigninController).Methods("POST")
	apiRouter.HandleFunc("/submit/form", Middlewares.Chain( Controllers.SubmitController , Middlewares.AuthMiddleware() )).Methods("POST")
	apiRouter.HandleFunc("/user/verify", Controllers.VerifyController).Methods("POST")

    log.Printf("The Serve is running at port 4100")
	log.Fatal(http.ListenAndServe(":4100",r))
}
