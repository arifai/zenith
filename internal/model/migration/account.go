package migration

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/arifai/zenith/pkg/errormessage"
	logg "github.com/arifai/zenith/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AccountMigration performs the migration of the Account and AccountPassHashed tables. It drops the existing tables, creates
// new ones, and inserts a default account with a hashed password if it does not already exist. All operations are handled
// within a transaction to ensure atomicity.
const defaultPassword = "12345678"

// AccountMigration performs the migration of the Account and AccountPassHashed tables.
func (m *Migration) AccountMigration() {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := migrateAccount(tx); err != nil {
			return err
		}

		account := createDefaultAccount(m.id)
		hashedPassword, err := hashDefaultPassword()
		if err != nil {
			return err
		}
		accountPassHashed := &model.AccountPassHashed{AccountId: account.ID, PassHashed: hashedPassword}

		if err := updateOrInsertRecord(tx, &account, "email = ?", account.Email); err != nil {
			return err
		}
		if err := updateOrInsertRecord(tx, accountPassHashed, "account_id = ?", account.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logg.Logger.Error(errormessage.ErrMigrationText, zap.String("migration_name", "account"), zap.Error(err))
	}
}

// migrateAccount performs the migration of Account and AccountPassHashed tables within a transaction, ensuring atomicity.
func migrateAccount(tx *gorm.DB) error {
	if err := tx.Migrator().DropTable(&model.Account{}, &model.AccountPassHashed{}); err != nil {
		return err
	}
	if err := tx.AutoMigrate(&model.Account{}, &model.AccountPassHashed{}); err != nil {
		return err
	}
	return nil
}

// createDefaultAccount creates a new account with default values for ID, FullName, Email, Avatar, and sets Active to true.
func createDefaultAccount(id uuid.UUID) *model.Account {
	return &model.Account{
		ID:       id,
		FullName: "John Doe",
		Email:    "john.doe@mail.com",
		Avatar:   "https://api.dicebear.com/9.x/notionists/png?scale=130&size=260&backgroundColor=b6e3f4",
		Active:   true,
	}
}

// hashDefaultPassword generates a hashed representation of a predefined default password using the Argon2ID algorithm.
// Returns the encoded hash or an error if the hashing process fails.
func hashDefaultPassword() (string, error) {
	p := crypto.DefaultArgon2IDHash
	return p.GenerateHash([]byte(defaultPassword), nil)
}

// updateOrInsertRecord updates a record based on the provided query, or inserts it if no record matches the query.
func updateOrInsertRecord(tx *gorm.DB, record interface{}, query string, args ...interface{}) error {
	if tx.Model(record).Where(query, args...).Updates(record).RowsAffected == 0 {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
	}
	return nil
}
