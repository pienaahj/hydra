package hydradblayer

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type mySqlDataStore struct {
	*sql.DB
}

func NewMySQLDataStore(conn string) (*mySqlDataStore, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Printf("error opening connection, %s\n", conn)
		return nil, err
	}
	return &mySqlDataStore{DB: db}, nil
}

func (msql *mySqlDataStore) AddMember(cm *CrewMember) error {
	_, err := msql.Exec("INSERT INTO Personnel (Name,SecurityClearance,Position) VALUES (?,?,?)", cm.Name, cm.SecClearance, cm.Position)
	return err
}
func (msql *mySqlDataStore) FindMember(id int) (CrewMember, error) {
	row := msql.QueryRow("SELECT * FROM Personnel WHERE id = ?", id)
	cm := CrewMember{}
	err := row.Scan(&cm.ID, &cm.Name, &cm.SecClearance, &cm.Position)
	return cm, err
}
func (msql *mySqlDataStore) AllMembers() (crew, error) {
	rows, err := msql.Query("SELECT * FROM Personnel;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := crew{}
	for rows.Next() {
		member := CrewMember{}
		err := rows.Scan(&member.ID, &member.Name, &member.SecClearance, &member.Position)
		if err == nil {
			members = append(members, member)
		}
	}
	err = rows.Err()
	return members, err
}
