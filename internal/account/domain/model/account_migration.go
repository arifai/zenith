package model

import (
	"github.com/arifai/go-modular-monolithic/pkg/crypto"
	"gorm.io/gorm"
	"log"
)

// AccountMigration will create the account table and insert the initial data
func AccountMigration(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&Account{}, &UserPassHashed{})
	if err != nil {
		log.Fatalf("Error migrating account model: %v", err)
	}

	account := &Account{
		ID:       1,
		FullName: "John Doe",
		Email:    "john.doe@mail.com",
		Avatar:   "https://api.dicebear.com/9.x/notionists/png?scale=130&size=260&backgroundColor=b6e3f4",
		Active:   true,
	}

	p := crypto.Argon2IdHash{Time: 1, Memory: 64 * 1024, Threads: 4, KeyLen: 32, SaltLen: 16}
	hashSalt, err := crypto.GenerateHash(p, []byte("12345678"), nil)
	if err != nil {
		log.Fatalf("Error generating hash: %v", err)
	}

	passHashed := &UserPassHashed{
		AccountId:  account.ID,
		PassHashed: hashSalt,
	}

	if db.Model(&account).Where("email = ?", account.Email).Updates(&account).RowsAffected == 0 {
		db.Create(&account)
	}

	if db.Model(&passHashed).Where("account_id = ?", account.ID).Updates(&passHashed).RowsAffected == 0 {
		db.Create(&passHashed)
	}
}
