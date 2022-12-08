package hydradblayer

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkMySQLDBReads(b *testing.B) {
	dblayer, err := ConnectDatabase(b * testing.B)
	if err != nil {
		b.Fatal("Could not connect to hydra chat system", err)
	}

	findMembersBM(b, dbLayer) // find members benchmank
}

func BenchmarkMongoDBReads(b *testing.B) {
	var connectionStringAdmin string = "mongodb://admin:myadminpassword@192.168.0.148:27017"
	dblayer, err := ConnectDatabase("mongodb", connectionStringAdmin)
	if err != nil {
		b.Error("Could not connect to the hydra chat system", err)
	}

	findMembersBM(b, dbLayer)
}

func findMemberBM(b *testing.B, dblayer dbLayer) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		_, err := dblayer.FindMember(rand.Intn(16) + 1)

		if err != nil {
			b.Error("Query failed ", err)
			return
		}
	}
}
