package hydradblayer

import (
	"context"
	"errors"
	"log"
)

const (
	mongodb = "mongodb"
	mysql   = "mysql"
)

var ErrDBTypeNotFound = errors.New("Database type not found...")

type DBLayer interface {
	AddMember(cm *CrewMember) error
	FindMember(id int) (CrewMember, error)
	AllMembers() (crew, error)
}

type CrewMember struct {
	ID           int    `bson:"id" json:"id"`
	Name         string `bson:"name" json:"name"`
	SecClearance int    `bson:"security clearance" json:"security clearance"`
	Position     string `bson:"position" json:"position"`
}

type crew []CrewMember

// ConnectDatabase connects to a database type o using provided connection string
func ConnectDatabase(o string, cstring string) (DBLayer, error) {
	switch o {
	case mongodb:
		return NewMongoStore(context.TODO(), cstring)
	case mysql:
		return NewMySQLDataStore(cstring)
	}
	log.Println("Could not find ", o)
	return nil, ErrDBTypeNotFound
}
