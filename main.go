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
	"google.golang.org/api/option"
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
	RoomId   string    `firestore:"RoomId,omitempty"`
	Sender   string    `firestore:"Sender,omitempty"`
	SendTime time.Time `firestore:"SendTime,omitempty"`
}

func main() {
	ctx := context.Background()
	projectID, serviceAccountID, err := loadProjectIDFromEnvFile()
	fmt.Println(projectID)
	if err != nil {
		log.Fatalln(err)
	}
	conf := &firebase.Config{
		ProjectID:        projectID,
		ServiceAccountID: serviceAccountID,
	}
	opt := option.WithCredentialsFile("./secret.json")
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	// addCityData(ctx, client)
	// addNotificationData(ctx, client)
	// getNotifications(ctx, client)
	// roomIDs := []string{"1", "2", "3", "4", "5", "19"}
	// listenDocuments(ctx, client, roomIDs)
	// listenRoomNotifications(ctx, client)
	// listenWorkspaceNotifications(ctx, client, "workspace_1")
	listenUseBaseNotifications(ctx, client)
}

func listenDocuments(ctx context.Context, client *firestore.Client, roomIDs []string) error {
	fmt.Println("call lintenDocuments()")
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	// lastFetchTime := time.Date(2023, 02, 22, 19, 00, 00, 00, time.UTC)
	lastFetchTime := time.Now().Add(time.Hour * -1)

	snapIter := client.Collection("notifications").Where("RoomId", "in", roomIDs).Where("SendTime", ">", lastFetchTime).Snapshots(ctx)
	defer snapIter.Stop()

	for {
		snap, err := snapIter.Next()
		if err != nil {
			log.Fatalln(err)
		}

		if snap != nil {
			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("data: %+v \n", doc.Data())
			}
		}
	}
}

func listenRoomNotifications(ctx context.Context, client *firestore.Client) {
	fmt.Println("All notifications: ")
	roomIDs := []string{"3", "6", "19"}
	snapIter := client.Collection("rooms").Where("room_id", "in", roomIDs).Snapshots(ctx)
	for {
		snap, err := snapIter.Next()
		if err != nil {
			log.Fatalln(err)
		}

		if snap != nil {
			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("data: %+v \n", doc.Data())
			}
		}
	}
}

func listenWorkspaceNotifications(ctx context.Context, client *firestore.Client, workSpace string) {
	fmt.Println("All notifications: ")
	roomIDs := []string{"1", "2"}
	snapIter := client.Collection("workspaces").Doc(workSpace).Collection("notifications").Where("RoomId", "in", roomIDs).Snapshots(ctx)
	for {
		snap, err := snapIter.Next()
		if err != nil {
			log.Fatalln(err)
		}

		if snap != nil {
			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("data: %+v \n", doc.Data())
			}
		}
	}
}

func listenUseBaseNotifications(ctx context.Context, client *firestore.Client) {
	fmt.Println("All notifications: ")
	snapIter := client.Collection("users").Doc("user1").Collection("rooms").Snapshots(ctx)
	for {
		snap, err := snapIter.Next()
		if err != nil {
			log.Fatalln(err)
		}

		if snap != nil {
			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("data: %+v \n", doc.Data())
			}
		}
	}
}

func getNotifications(ctx context.Context, client *firestore.Client) {
	fmt.Println("All notifications: ")
	roomIDs := []string{"2", "4", "19"}
	iter := client.Collection("notifications").Where("RoomId", "in", roomIDs).Documents(ctx)
	for {
		fmt.Println("1")
		notificationDoc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		fmt.Println("2")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(notificationDoc.Data())
	}
}

func addRoomNotificationData(ctx context.Context, client *firestore.Client) {
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
		{roomId: "3", notifications: []Notification{
			{Sender: "hoge5", SendTime: time.Now().Add(time.Hour * -1)},
			{Sender: "hoge9", SendTime: time.Now().Add(time.Minute * -30)},
			{Sender: "hoge8", SendTime: time.Now()},
			{Sender: "hoge2", SendTime: time.Now().Add(time.Second * -1)},
			{Sender: "hoge10", SendTime: time.Now().Add(time.Hour * -10)},
		}},
		{roomId: "6", notifications: []Notification{
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

func addNotificationData(ctx context.Context, client *firestore.Client) {
	roomNotifications := []struct {
		notifications []Notification
	}{
		{notifications: []Notification{
			{RoomId: "1", Sender: "hoge1", SendTime: time.Now().Add(time.Hour * -10)},
			{RoomId: "1", Sender: "hoge2", SendTime: time.Now().Add(time.Minute * -10)},
			{RoomId: "1", Sender: "hoge8", SendTime: time.Now()},
			{RoomId: "1", Sender: "hoge2", SendTime: time.Now().Add(time.Second * -10)},
			{RoomId: "1", Sender: "hoge1", SendTime: time.Now().Add(time.Hour * -1)},
		}},
		{notifications: []Notification{
			{RoomId: "3", Sender: "hoge5", SendTime: time.Now().Add(time.Hour * -1)},
			{RoomId: "3", Sender: "hoge9", SendTime: time.Now().Add(time.Minute * -30)},
			{RoomId: "3", Sender: "hoge8", SendTime: time.Now()},
			{RoomId: "3", Sender: "hoge2", SendTime: time.Now().Add(time.Second * -1)},
			{RoomId: "3", Sender: "hoge10", SendTime: time.Now().Add(time.Hour * -10)},
		}},
		{notifications: []Notification{
			{RoomId: "6", Sender: "hoge22", SendTime: time.Now().Add(time.Minute * -10)},
			{RoomId: "6", Sender: "hoge19", SendTime: time.Now().Add(time.Minute * -9)},
			{RoomId: "6", Sender: "hoge30", SendTime: time.Now()},
			{RoomId: "6", Sender: "hoge22", SendTime: time.Now().Add(time.Hour * -1)},
			{RoomId: "6", Sender: "hoge1", SendTime: time.Now().Add(time.Second * -10)},
		}},
	}

	batch := client.Batch()

	for _, roomNotification := range roomNotifications {
		for _, notf := range roomNotification.notifications {
			ref := client.Collection("notifications").NewDoc()
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

func loadProjectIDFromEnvFile() (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", "", err
	}
	return os.Getenv("ProjectID"), os.Getenv("ServiceAccountID"), nil
}
