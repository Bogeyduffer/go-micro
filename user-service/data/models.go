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
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string    `bson:"email" json:"email"`
	FirstName string    `bson:"firstName" json:"firstName"`
	LastName  string    `bson:"lastName" json:"lastName"`
	Password  string    `bson:"password" json:"password"`
	Active    string    `bson:"active" json:"active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *User) Insert(entry User) error {
	collection := client.Database("users").Collection("users")

	_, err := collection.InsertOne(context.TODO(), User{
		Email:     entry.Email,
		FirstName: entry.FirstName,
		LastName:  entry.LastName,
		Password:  entry.Password,
		Active:    entry.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into users:", err)
		return err
	}

	return nil
}

func (l *User) All() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("users").Collection("users")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*User

	for cursor.Next(ctx) {
		var item User

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding user into slice:", err)
			return nil, err
		} else {
			users = append(users, &item)
		}
	}

	return users, nil
}

func (l *User) GetOne(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("users").Collection("users")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry User
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *User) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("users").Collection("users")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}
func (l *User) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("users").Collection("users")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"email", l.Email},
				{"first_name", l.FirstName},
				{"last_name", l.LastName},
				{"password", l.Password},
				{"active", l.Active},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

