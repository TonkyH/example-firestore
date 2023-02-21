package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

// City represents a city.
type City struct {
	Name       string   `firestore:"name,omitempty"`
	State      string   `firestore:"state,omitempty"`
	Country    string   `firestore:"country,omitempty"`
	Capital    bool     `firestore:"capital,omitempty"`
	Population int64    `firestore:"population,omitempty"`
	Regions    []string `firestore:"regions,omitempty"`
}

type Room struct {
	Id int64
}

type Notification struct {
	Sender   string    `firestore:"Sender,omitempty"`
	SendTime time.Time `firestore:"SendTime,omitempty"`
}

func main() {
	ctx := context.Background()
	projectID, err := loadProjectIDFromEnvFile()
	if err != nil {
		log.Fatalln(err)
	}
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	// addCityData(ctx, client)
	addNotificationData(ctx, client)
	// getNotifications(ctx, client)

}

func getNotifications(ctx context.Context, client *firestore.Client) {
	fmt.Println("All notifications: ")
	roomIds := []int32{1, 2}
	iter := client.Collection("rooms").Where("room_id", "in", roomIds).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}

func addNotificationData(ctx context.Context, client *firestore.Client) {
	roomNotifications := []struct {
		roomId        string
		notifications []Notification
	}{
		{roomId: "1", notifications: []Notification{
			{Sender: "hoge1", SendTime: time.Now().Add(time.Hour * -10)},
			{Sender: "hoge2", SendTime: time.Now().Add(time.Minute * -10)},
			{Sender: "hoge8", SendTime: time.Now()},
			{Sender: "hoge2", SendTime: time.Now().Add(time.Second * -10)},
			{Sender: "hoge1", SendTime: time.Now().Add(time.Hour * -1)},
		}},
		{roomId: "6", notifications: []Notification{
			{Sender: "hoge5", SendTime: time.Now().Add(time.Hour * -1)},
			{Sender: "hoge9", SendTime: time.Now().Add(time.Minute * -30)},
			{Sender: "hoge8", SendTime: time.Now()},
			{Sender: "hoge2", SendTime: time.Now().Add(time.Second * -1)},
			{Sender: "hoge10", SendTime: time.Now().Add(time.Hour * -10)},
		}},
		{roomId: "19", notifications: []Notification{
			{Sender: "hoge22", SendTime: time.Now().Add(time.Minute * -10)},
			{Sender: "hoge19", SendTime: time.Now().Add(time.Minute * -9)},
			{Sender: "hoge30", SendTime: time.Now()},
			{Sender: "hoge22", SendTime: time.Now().Add(time.Hour * -1)},
			{Sender: "hoge1", SendTime: time.Now().Add(time.Second * -10)},
		}},
	}

	batch := client.Batch()

	for _, roomNotification := range roomNotifications {
		for _, notf := range roomNotification.notifications {
			ref := client.Collection("rooms").Doc(roomNotification.roomId).Collection("notifications").NewDoc()
			batch.Set(ref, notf)
		}
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		log.Fatalf("An has occurred: %s", err)
	}
}

func addCityData(ctx context.Context, client *firestore.Client) {
	cities := []struct {
		id string
		c  City
	}{
		{id: "SF", c: City{Name: "San Francisco", State: "CA", Country: "USA", Capital: false, Population: 860000}},
		{id: "LA", c: City{Name: "Los Angeles", State: "CA", Country: "USA", Capital: false, Population: 3900000}},
		{id: "DC", c: City{Name: "Washington D.C.", Country: "USA", Capital: true, Population: 680000}},
		{id: "TOK", c: City{Name: "Tokyo", Country: "Japan", Capital: true, Population: 9000000}},
		{id: "BJ", c: City{Name: "Beijing", Country: "China", Capital: true, Population: 21500000}},
	}
	for _, c := range cities {
		_, err := client.Collection("cities").Doc(c.id).Set(ctx, c.c)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}
	}
}

func loadProjectIDFromEnvFile() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	return os.Getenv("ProjectID"), nil
}
