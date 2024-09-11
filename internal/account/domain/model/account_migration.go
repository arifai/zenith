package model

import (
	"github.com/arifai/go-modular-monolithic/pkg/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

// AccountMigration will create the account table and insert the initial data
func AccountMigration(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().AutoMigrate(&Account{}, &AccountPassHashed{})
		if err != nil {
			return err
		}

		account := &Account{
			ID:       uuid.New(),
			FullName: "John Doe",
			Email:    "john.doe@mail.com",
			Avatar:   "https://api.dicebear.com/9.x/notionists/png?scale=130&size=260&backgroundColor=b6e3f4",
			Active:   true,
		}

		p := crypto.Argon2IdHash{Time: 1, Memory: 64 * 1024, Threads: 4, KeyLen: 32, SaltLen: 16}
		hashSalt, err := p.GenerateHash([]byte("12345678"), nil)
		if err != nil {
			return err
		}

		userPassHashed := &AccountPassHashed{
			AccountId:  account.ID,
			PassHashed: hashSalt,
		}

		if tx.Model(&account).Where("email = ?", account.Email).Updates(&account).RowsAffected == 0 {
			if err := tx.Create(&account).Error; err != nil {
				return err
			}
		}

		if tx.Model(&userPassHashed).Where("account_id = ?", account.ID).Updates(&userPassHashed).RowsAffected == 0 {
			if err := tx.Create(&userPassHashed).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error during account migration: %v", err)
	}
}
