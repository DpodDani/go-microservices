package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func getObjectID(id string) (*primitive.ObjectID, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Couldn't create ObjectID for id: %s\n", id)
		return nil, err
	}

	return &docID, nil
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	// collections are analogous with tables in relational databases
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("❗ - Failed to write log to MongoDB!")
	}

	return err
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cur, err := collection.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		log.Println("❗ - Failed to fetch all log entries")
		return nil, err
	}

	var logs []*LogEntry
	for cur.Next(ctx) {
		var item LogEntry

		err := cur.Decode(&item)
		if err != nil {
			log.Printf("Error decoding log entry: %s\n", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	docID, err := getObjectID(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		log.Printf("Couldn't decode document for id: %s\n", id)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Couldn't drop 'logs' collection")
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := getObjectID(l.ID)
	if err != nil {
		return nil, err
	}

	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "Name", Value: l.Name},
				{Key: "Data", Value: l.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		log.Printf("Failed to update for ID: %s\n", l.ID)
		return nil, err
	}

	return res, err
}
