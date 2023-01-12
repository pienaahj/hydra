package hydradblayer

import (
	"errors"
	"log"
)

const (
	mongodb = "mongodb"
	mysql   = "mysql"
)

var ErrDBTypeNotFound = errors.New("database type not found ")

type DBLayer interface {
	AddMember(cm *CrewMember) error
	FindMember(id int) (CrewMember, error)
	AllMembers() (crew, error)
}

type CrewMember struct {
	ID           int    `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
	SecClearance int    `json:"clearance" bson:"security clearance"`
	Position     string `json:"position" bson:"position"`
}

type crew []CrewMember

// ConnectDatabase connects to a database type o using the provided connection string
func ConnectDatabase(o, cstring string) (DBLayer, error) {
	switch o {
	case mongodb:
		return NewMongoStore(cstring)
	case mysql:
		return NewMySQLDataStore(cstring)
	}
	log.Println("Could not find ", o)
	return nil, ErrDBTypeNotFound
}
