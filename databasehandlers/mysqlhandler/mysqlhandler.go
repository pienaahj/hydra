package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type crewMember struct {
	id           int
	name         string
	secClearance int
	position     string
}

type Crew []crewMember

var connectionString string = "gouser:gouser@tcp(127.0.0.1:3306)/Hydra?parseTime=true"

func main() {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Could not connect, error ", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect, error ", err.Error())
	}

	//  populate the database from csv file
	// CSVToMySQL(db)

	cw := GetCrewByPositions(db, []string{"'Mechanic'", "'Biologist'"})
	fmt.Println(cw)

	fmt.Println(GetCrewMemberById(db, 11))
	// output {11 Noble Luczynski   4 Mechanic}

	fmt.Println(GetCrewMemberByPosition(db, "Chemist"))
	// output [{7 Marcus Durkee   2 Chemist} {14 Christel Sample   6 Chemist}]

	// AddCrewMember(db, crewMember{name: "Steve Bee", secClearance: 4, position: "Biologist"})
	// output 2022/12/14 19:35:36 Rows affected 1 Last inserted id 18

	// cr := Crew{
	// 	crewMember{name: "Adam stler", secClearance: 4, position: "Chemist"},
	// 	crewMember{name: "Zach Garph", secClearance: 5, position: "Biologist"},
	// }
	// CreateCrewMembersByTransaction(db, cr)
}

func GetCrewByPositions(db *sql.DB, positions []string) Crew {
	Qs := fmt.Sprintf("SELECT id,Name,SecurityClearance,Position FROM Personnel WHERE Position IN (%s)", strings.Join(positions, ","))
	// SELECT id,Name,SecurityClearance,Position form Personnel where Posision in (mechanic,biologist)
	rows, err := db.Query(Qs)
	if err != nil {
		log.Fatal("Could not get the data from the Personnel table ", err)
	}
	defer rows.Close()

	retVal := Crew{}
	cols, _ := rows.Columns()
	fmt.Println("Columns detected: ", cols)

	for rows.Next() {
		member := crewMember{}
		err := rows.Scan(&member.id, &member.name, &member.secClearance, &member.position)
		if err != nil {
			log.Fatal("Error scanning row", err)
		}
		retVal = append(retVal, member)

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

	}
	return retVal
}

func GetCrewMemberById(db *sql.DB, id int) (cm crewMember) {
	row := db.QueryRow("Select * from Personnel where id = ?", id)

	err := row.Scan(&cm.id, &cm.name, &cm.secClearance, &cm.position)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func GetCrewMemberByPosition(db *sql.DB, position string) (cr Crew) {

	stmt, err := db.Prepare("Select * from Personnel where Position = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(position)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		cm := crewMember{}
		err = rows.Scan(&cm.id, &cm.name, &cm.secClearance, &cm.position)
		if err != nil {
			log.Fatal(err)
		}
		cr = append(cr, cm)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func AddCrewMember(db *sql.DB, cm crewMember) int64 {
	res, err := db.Exec("INSERT INTO Personnel (Name, SecurityClearance, Position) VALUES (?, ?, ?)", cm.name, cm.secClearance, cm.position)
	if err != nil {
		log.Fatal(err)
	}

	ra, _ := res.RowsAffected()
	id, _ := res.LastInsertId()

	log.Println("Rows affected", ra, "Last inserted id", id)
	return id
}

func CreateCrewMembersByTransaction(db *sql.DB, cr Crew) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Could not begin tx", err)
	}

	stmt, err := tx.Prepare("INSERT INTO Personnel (Name, SecurityClearance, Position) VALUES (?,?,?);")
	if err != nil {
		tx.Rollback()
		log.Fatal("Could not do select statement ", err)
	}
	defer stmt.Close()
	for _, person := range cr {
		_, err := stmt.Exec(person.name, person.secClearance, person.position)
		if err != nil {
			tx.Rollback()
			log.Fatal("Could not query positions ", err)
		}
	}
	tx.Commit()
	return
}
func CSVToMySQL(db *sql.DB) {
	file, err := os.Open("Crews.csv")
	if err != nil {
		log.Fatal("Could not open CSV file", err)
	}
	defer file.Close()

	vargs := []interface{}{}
	sargs := []string{}

	r := csv.NewReader(file)
	r.Comment = '#'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		for _, rec := range record {
			vargs = append(vargs, rec)
		}
		//vargs = append(vargs,record...)
		sargs = append(sargs, "(?,?,?,?)")
	}

	insertStmt := fmt.Sprintf("INSERT INTO Personnel (id,Name,SecurityClearance,Position) VALUES %s ", strings.Join(sargs, ","))
	_, err = db.Exec(insertStmt, vargs...)
	if err != nil {
		log.Fatalf("Could not execute insert statement %s with args %s, error %s \n", insertStmt, vargs, err.Error())
	}

}
