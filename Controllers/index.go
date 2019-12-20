package Controllers

import "fmt"
import "encoding/json"
import "net/http"
import "time"
import	"google.golang.org/api/iterator"
import 	"context"
import "github.com/dgrijalva/jwt-go"
import "github.com/google/uuid"
import "../Db"
import "../Helpers"
import "../Modals"
import "log"


var jwtKey = []byte("my_secret_key")


//type MyCustomClaims jwt.StandardClaims


var jwtToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzYwODgwOTYsImlhdCI6MTU3NjA4NDQ5Niwic3ViIjoiNTVjODRiLTFjMmItMTFlYS1iY2YwLTg0YTkzZTM5MjY0NCJ9.0xBk-0d-OHPNWqsuGXLxCQMzEn-8BILnzzSZ_c80XkU"


func PaymentController(w http.ResponseWriter, r *http.ResponseWriter) {
	
}



func HomeController(w http.ResponseWriter , r *http.Request) {

	claims := &jwt.StandardClaims{}
	
	tkn, err := jwt.ParseWithClaims(jwtToken, claims , func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

 // 	claims := tkn.Claims.(jwt.MapClaims)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Fprintln(w , err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprintln(w , err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		fmt.Fprintf(w , "BADDD")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}


	log.Printf("%+v" , claims )
	fmt.Fprintln(w , claims.Subject )
}


func SignupController(w http.ResponseWriter , r *http.Request) {

	ctx := context.Background()

	User := &Modals.UserCrud{
		Email : r.FormValue("email") ,
		Name : r.FormValue("name") ,
		Password : r.FormValue("password") ,
	}

	err := User.Validator()

	if err != nil {
		log.Printf("%+v", err)
		Helpers.Error(w , "Please fill right info" , 401)
		return
	}

	query := Db.Client.Collection("users").Where("email", "==", User.Email).Limit(1).Documents(ctx)

	_ , exit := query.Next()

	//if user already exists

	if exit != iterator.Done {
		Helpers.Error(w , "User already exits" , 409)
		return
	}

	id, err := uuid.NewUUID()

	if err != nil {
		Helpers.Error(w , "Server Error" , 500)
		log.Fatalf("uuid.NewV4() failed with %s\n", err)
		return
	}

	log.Printf("%+v" , id.String())

	User.UserID  = id.String();

	log.Printf("%+v",User)

	User.Password , err = Helpers.HashPassword(User.Password)

	if err != nil {
		Helpers.Error(w , "Server Error" , 500)
		log.Fatalf("%+v", err)
		return
	}

	log.Printf("%+v",User)


	_, err = Db.Client.Collection("users").Doc(User.UserID).Set(ctx, User)

	if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
			Helpers.Error(w , "Serve Error" , 502)
			return
	}
	Helpers.Message(w , "You are in" , 200)
	return
}


func SigninController(w http.ResponseWriter , r *http.Request) {

	ctx := context.Background()

	User := &Modals.LoginCrud{
		Email : r.FormValue("email") ,
		Password : r.FormValue("password") ,
	}

	err := User.Validator()

	if err != nil {
		log.Printf("%+v", err)
		Helpers.Error(w , "Please fill right info" , 401)
		return
	}
	log.Printf("Sandhu");
	
	query := Db.Client.Collection("users").Where("email", "==", User.Email).Limit(1).Documents(ctx)
	log.Printf("Sandhu");
	Suser , exit := query.Next()

	//if user already exists

	if exit == iterator.Done {
		Helpers.Error(w , "User not found" , 404)
		return
	}

	err = Helpers.CheckPasswordHash( User.Password , Suser.Data()["password"].(string) )
	
	if err != nil {
		Helpers.Error(w , "Wrong Password" , 409)
		log.Printf("%+v", err)
		return
	}


	log.Printf("%+v",User)


	userId , _ := Suser.Data()["user_id"].( string );



	log.Printf("\n\n%+v\n\n" , userId)


	log.Printf("%+v" , Suser.Data())

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(time.Hour * 1).Unix()
    token.Claims.(jwt.MapClaims)["iat"] = time.Now().Unix()
    token.Claims.(jwt.MapClaims)["sub"] = userId

	// Create the JWT string
	tokenString, tokenError := token.SignedString(jwtKey)

	if tokenError != nil {
		// If there is an error in creating the JWT return an internal server error
		Helpers.Error(w , "Server Error" , 500)
		return
	}
	
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct{ Token string `json:"token"`}{ Token : tokenString })
	return
}


func SubmitController (w http.ResponseWriter , r *http.Request) {
	ctx := context.Background()

	Form := &Modals.FormCrud{
		Email : r.FormValue("email1"),
		Phone : r.FormValue("phone"),
		Deadline : r.FormValue("deadline"),
		Detail : r.FormValue("detail"),
		Lang : r.FormValue("lang"),
		Type : r.FormValue("type"),
		FileUrl : "https://www.google.com",
	}

	log.Printf( "%+v" , Form )

	err := Form.Validator()

	if err != nil {
		log.Printf("%+v", err)
		Helpers.Error(w , "Wrong information" , 409)
		return
	}

	user_id := r.Header.Get("user")

	fmt.Printf("\n%+v ::: %+v\n",user_id , r.Header)

	_, err = Db.Client.Collection("forms").Doc(user_id).Set(ctx,	Form)

	if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
			Helpers.Error(w , "Serve Error" , 502)
			return
	}

	Helpers.Message(w, "Form Submiited" , 200)
	return
}



func VerifyController (w http.ResponseWriter , r *http.Request) {

	tokenString := r.Header.Get("x-access-token")
		
	claims := &jwt.StandardClaims{}

	w.Header().Set("Content-Type","application/json")
	
	tkn, err := jwt.ParseWithClaims(tokenString, claims , func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	log.Printf("%+v" , claims )

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			// json.NewEncoder(w).Encode(struct{ valid bool }{ valid : false })
			// w.WriteHeader(http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(struct{ Valid string }{ Valid : "false" })
			//Error(w , "Not valid" ,http.StatusUnauthorized)
			return
		}
		// json.NewEncoder(w).Encode(struct{ valid bool }{ valid : false })
		// w.WriteHeader(http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Valid string }{ Valid : "false" })
		//Error(w , "Not valid" ,http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(struct{ Valid string }{ Valid : "false" })
		//Error(w , "Not valid" ,http.StatusUnauthorized)
		return
	}


	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct{ Valid string }{ Valid : "true" })
	return
}