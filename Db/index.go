package Db


import "log"
import 	"context"
import "google.golang.org/api/option"
import firebase "firebase.google.com/go"
import "cloud.google.com/go/firestore"

var Client *firestore.Client


func InitDb(){
	ctx:= context.Background()

	//databae init
	opt := option.WithCredentialsFile("./Config/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	
	if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
	}

	Client, _ = app.Firestore(ctx)

	if err != nil {
	log.Fatalln(err)
	}
	log.Printf("%+v",*Client);
}

