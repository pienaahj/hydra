package main

import (
	"flag"
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

func main() {
	// Open a boltdb connection
	db, err := bolt.Open("hydra.db", 0600, nil)
	if err != nil {
		log.Fatal("Cannot connect to bolt database", err)
	}
	defer db.Close()

	// add command line arguments
	op := flag.String("op", "", "Add or Get operation")
	user := flag.String("u", "", "Enter username")
	pass := flag.String("p", "", "Enter password")
	flag.Parse()

	switch strings.ToUpper(*op) {
	case "ADD":
		if len(*user) != 0 && len(*pass) != 0 {
			err = addToVault(db, *user, *pass)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "GET":
		if len(*user) != 0 {
			pass, err := getPassword(db, *user)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Found password: ", pass)
		}
	}

}

// addToVault adds a password to the vault
func addToVault(db *bolt.DB, username, password string) error {
	return db.Update(func(tx *bolt.Tx) error {
		// if the bucket does not exist - create it
		b, err := tx.CreateBucketIfNotExists([]byte("PasswordVault"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(username), []byte(password))
		return err
	})
}

// getPassword will retrieve the password from the vault
func getPassword(db *bolt.DB, username string) (string, error) {
	password := ""
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PasswordVault"))
		v := b.Get([]byte(username))
		password = string(v)
		return nil
	})
	return password, err
}

// add Nathan to the vault
// go run boltdbhandler.go -op Add -u Nathan -p Nat
// go run boltdbhandler.go -op Get -u Nathan
