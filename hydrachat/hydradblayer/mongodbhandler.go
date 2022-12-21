package hydradblayer

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDataStore struct {
	*mongo.Client
}

// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// defer cancel()
// defer func() {
// 	if err = client.Disconnect(ctx); err != nil {
// 		log.Fatal("Connection to mongodbclosed", err)
// 	}
// }()

func NewMongoStore(ctx context.Context, connectionString string) (*mongoDataStore, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	return &mongoDataStore{client}, nil
}

func (ms *mongoDataStore) AddMember(cm *CrewMember) error {
	personnel := ms.Client.Database("Hydra").Collection("Personnel")
	// Adding a new crew member
	insertResult, err := personnel.InsertOne(context.TODO(), cm)
	if err != nil {
		return err
	}
	log.Printf("Inserted a new crew member: %v\n", insertResult.InsertedID)
	return nil
}

func (ms *mongoDataStore) FindMember(id int) (CrewMember, error) {
	personnel := ms.Client.Database("Hydra").Collection("Personnel")
	filter := bson.D{{"id", id}}
	cm := crewMember{}
	err := personnel.FindOne(context.TODO(), filter).Decode(&cm)
	if err != nil {
		return CrewMember{}, err
	}
	return cm, nil
}

func (ms *mongoDataStore) AllMembers() (crew, error) {
	personnel := ms.Client.Database("Hydra").Collection("Personnel")
	members := crew{}
	cursor, err := personnel.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return crew{}, err
	}
	if err = cursor.All(context.TODO(), &members); err != nil {
		return crew{}, err
	}
	return members, nil
}
