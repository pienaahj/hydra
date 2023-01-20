package main

import (
	"crypto/md5"

	"github.com/pienaahj/hydra/hydradblayer/passwordvault"
)

func main() {
	db, err := passwordvault.ConnectPasswordVault()
	if err != nil {
		return
	}
	minapass := md5.Sum([]byte("minapass"))
	jimpass := md5.Sum([]byte("jimpass"))
	caropass := md5.Sum([]byte("caropass"))
	passwordvault.AddBytesToVault(db, "Mina", minapass[:])
	passwordvault.AddBytesToVault(db, "Jim", jimpass[:])
	passwordvault.AddBytesToVault(db, "Carop", caropass[:])
	db.Close()
}
