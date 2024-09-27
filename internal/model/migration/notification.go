package migration

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

// NotificationMigration performs the internal of the PushNotification table by dropping the existing table and creating a new one.
func (m *Migration) NotificationMigration() {
	if err := createEnums(m.db); err != nil {
		log.Fatalf("error creating enums table: %v", err)
		return
	}

	if err := migrateNotification(m, &model.PushNotification{}, &model.Notification{}); err != nil {
		log.Fatalf("Error during notification internal: %v", err)
	}

	notifications := createDummiesNotifications(m.id)

	for _, notification := range notifications {
		if err := insertRecord(m.db, notification); err != nil {
			log.Fatalf("Error during insert notification internal: %v", err)
			return
		}
	}
}

// migrateNotification handles the internal for the given model by dropping the existing table and creating a new one.
func migrateNotification(m *Migration, models ...interface{}) error {
	for _, i := range models {
		if m.db.Migrator().HasTable(i) {
			if err := m.db.Migrator().DropTable(i); err != nil {
				return err
			}
		}

		if err := m.db.AutoMigrate(i); err != nil {
			return err
		}
	}
	return nil
}

// createEnums creates the necessary enums in PostgreSQL.
func createEnums(db *gorm.DB) error {
	if err := db.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'platform') THEN CREATE TYPE platform AS ENUM ('Android', 'iOS', 'Web'); END IF; END $$;").Error; err != nil {
		return err
	}

	if err := db.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN CREATE TYPE status AS ENUM ('Pending', 'Success', 'Failure'); END IF; END $$;").Error; err != nil {
		return err
	}

	return nil
}

func createDummiesNotifications(id uuid.UUID) [2]*model.Notification {
	readAt := time.Now()

	return [2]*model.Notification{
		{
			ID:               uuid.New(),
			AccountID:        id,
			Title:            "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			Image:            "https://placehold.co/600x300/grey/white/png",
			ShortDescription: "Vestibulum ullamcorper nunc vel massa auctor, quis fringilla turpis fringilla",
			Description:      "Nunc a tristique massa. Quisque ut tellus arcu. Donec id neque elementum, porta ante id, bibendum neque. Vivamus non turpis sem. Suspendisse potenti. Nulla nibh sem, porttitor quis vulputate at, interdum sed augue. Nam volutpat luctus suscipit. Etiam at risus nec sem sollicitudin fringilla at id odio. Pellentesque elementum lacinia tortor, ac lobortis tortor congue a.",
			Read:             true,
			ReadAt:           &readAt,
		},
		{
			ID:               uuid.New(),
			AccountID:        id,
			Title:            "Vestibulum quis efficitur turpis",
			Image:            "https://placehold.co/600x300/grey/white/png",
			ShortDescription: "Nulla dictum nibh vel dapibus pellentesque",
			Description:      "Praesent lorem dui, rhoncus at gravida sed, mollis ut lorem. Nulla at gravida erat. Sed id ante ut neque cursus rhoncus et ut diam. Nulla id malesuada justo. In orci enim, dictum at tellus bibendum, fermentum convallis tortor. Quisque tincidunt, tellus vitae luctus pretium, leo ex laoreet dui, ut porta lectus dui nec nunc. Etiam gravida purus ex, sed efficitur ligula pharetra sit amet. Maecenas dapibus quam in mauris pretium, a vestibulum mi fermentum.",
		},
	}
}

func insertRecord(db *gorm.DB, record interface{}) error {
	if err := db.Create(record).Error; err != nil {
		return err
	}

	return nil
}
