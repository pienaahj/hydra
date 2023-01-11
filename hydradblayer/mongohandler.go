package hydradblayer

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDataStore struct {
	*mongo.Client
}

func NewMongoStore(conn string) *mongoDataStore {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		log.Fatal("Failed to connect mongodb client: ", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Connection to mongodb closed", err)
		}
	}()
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Cannot ping mongodb", err)
	}
	fmt.Println("Connected to MongoDB!")
	return &mongoDataStore{Client: client}
}

func (s *mongoDataStore) AddMember(cm *CrewMember) error {
	personnel := s.Database("Hydra").Collection("Personnel")
	insertResult, err := personnel.InsertOne(context.TODO(), cm)
	if err != nil {
		return fmt.Errorf("Error while inserting new crew member: %v", err)
	}
	fmt.Printf("Inserted a new crew member: %v\n", insertResult.InsertedID)
	return nil
}

func (s *mongoDataStore) FindMember(id int) (CrewMember, error) {
	personnel := s.Database("Hydra").Collection("Personnel")
	// Perform a single query
	filter := bson.D{{"id", id}}
	cm := CrewMember{}
	err := personnel.FindOne(context.TODO(), filter).Decode(&cm)
	if err != nil {
		return CrewMember{}, fmt.Errorf("Error while retrieving crew member with id %d: %v\n", id, err)
	}
	fmt.Printf("Crew member with id, %d is: %s\n", id, cm.Name)
	return cm, nil
}

func (s *mongoDataStore) AllMembers() (crew, error) {
	personnel := s.Database("Hydra").Collection("Personnel")
	filter := bson.D{{}}
	crew := crew{}
	findOptions := options.Find()
	cur, err := personnel.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return crew, fmt.Errorf("Error while retrieving the crew %v\n", err)
	}
	//  loop through the cursor
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var doc CrewMember
		err := cur.Decode(&doc)
		if err != nil {
			log.Printf("Error while decoding doc %v\n", err)
		}

		crew = append(crew, doc)

		// check for cursor errors
		if err := cur.Err(); err != nil {
			log.Printf("Cursor error: %v", err)
		}
	}
	cur.Close(context.TODO())
	return crew, nil
}

/*
func main() {
	var connectionString string
	viper.AutomaticEnv()
	if connectionStringT, ok := viper.Get("ME_CONFIG_MONGODB_URI").(string); !ok {
		log.Fatal("Cannot get env variable")
	} else {
		connectionString = connectionStringT
		// fmt.Printf("viper : %s = %s \n", "Connection string", connectionString)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal("Failed to connect mongodb client: ", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Connection to mongodbclosed", err)
		}
	}()
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Cannot ping mongodb", err)
	}

	fmt.Println("Connected to MongoDB!")

	personnel := client.Database("Hydra").Collection("Personnel")

	// Populate the Hydra database
	// CSVToMongo(personnel)
	// Get the number of documents in the collection
	filter := bson.D{{}}

	count, err := personnel.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal("Cannot count records in Personnel collection", err)
	}
	fmt.Printf(" Number of personnel: %d\n", count)

	// Perform a single query
	filter = bson.D{{"id", 3}}
	cm := crewMember{}
	err = personnel.FindOne(context.TODO(), filter).Decode(&cm)
	if err != nil {
		log.Printf("Error while retrieving ")
	}
	fmt.Printf("Crew member with is id, 3 is: %s\n", cm.Name)

	// Query with expression - Note! you need to usebson.M whenever new {} start
	filterM := bson.M{
		"security clearance": bson.M{"$gt": 2},
		"position":           bson.M{"$in": []string{"Mechanic", "Biologist"}}}

	findOptions := options.Find()
	//  Set the limit to the number of records to return
	findOptions.SetLimit(4)
	// Set the fields to return note:  cannot use in conjuction with filter
	// findOptions.SetProjection(bson.D{{"name", 1}, {"security clearance", 1}, {"position", 1}, {"id", 1}})
	// Make a holder to return to
	var cmX Crew

	cur, err := personnel.Find(context.TODO(), filterM, findOptions)
	if err != nil {
		log.Printf("Error while retrieving filter: %v with error: %v\n", filter, err)
	}
	//  loop through the cursor
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var doc crewMember
		err := cur.Decode(&doc)
		if err != nil {
			log.Printf("Error while decoding doc %v\n", err)
		}

		cmX = append(cmX, doc)

		// check for cursor errors
		if err := cur.Err(); err != nil {
			log.Printf("Cursor error: %v", err)
		}
	}
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array with lenght: %d): %v\n", len(cmX), cmX)

	// Adding a new crew member
	// newcr := crewMember{ID: 18, Name: "Kaya Gal", SecClearance: 4, Position: "Biologist"}
	// insertResult, err := personnel.InsertOne(context.TODO(), newcr)
	// if err != nil {
	// 	log.Printf("Error while inserting new crew member: %v", err)
	// }
	// fmt.Printf("Inserted a new crew member: %v\n", insertResult.InsertedID)

	//  Updating a crew member
	// filter = bson.D{{"id", 16}}
	// update := bson.D{{"$set", bson.D{{"position", "Engineer IV"}}}}
	// updateResult, err := personnel.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	// 	log.Printf("Could not update record: %v", err)
	// }
	// fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	//  Deleting a crew member
	// filter = bson.D{{"id", 16}}
	// deleteResult, err := personnel.DeleteOne(context.TODO(), filter)
	// if err != nil {
	// 	log.Printf("Could not delete record: %v", err)
	// }
	// fmt.Printf("Deleted document %v documents.\n", deleteResult.DeletedCount)

}
*/
