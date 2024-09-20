package migration

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

// AccountMigration performs the migration of the Account and AccountPassHashed tables. It drops the existing tables, creates
// new ones, and inserts a default account with a hashed password if it does not already exist. All operations are handled
// within a transaction to ensure atomicity.
func AccountMigration(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable(&model.Account{}, &model.AccountPassHashed{}); err != nil {
			return err
		}

		if err := tx.AutoMigrate(&model.Account{}, &model.AccountPassHashed{}); err != nil {
			return err
		}

		account := &model.Account{
			ID:       uuid.New(),
			FullName: "John Doe",
			Email:    "john.doe@mail.com",
			Avatar:   "https://api.dicebear.com/9.x/notionists/png?scale=130&size=260&backgroundColor=b6e3f4",
			Active:   true,
		}

		p := crypto.DefaultArgon2IDHash
		generatedHash, err := p.GenerateHash([]byte("12345678"), nil)
		if err != nil {
			return err
		}

		userPassHashed := &model.AccountPassHashed{AccountId: account.ID, PassHashed: generatedHash}

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
