package seeders

import (
	"eth-transaction-api/models"
	"log"

	"gorm.io/gorm"
)

func SeedAccounts(db *gorm.DB) {
	addresses := []string{
		"0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
		"0x7C456a5eA65E03fbb5F3Bd7B7Ec5f9A04C1a9E3D",
		"0xBf241F9aA6c542dE3d2b6C89D9A43Fc9149A2E35",
	}

	accountUuids := []string{
		"9b3af3a7-51f1-49a7-aa3b-c700cf82a835",
		"81f5c001-45a5-4922-8fcb-b961ae312ec0",
		"bb7b48b4-4481-4a72-8079-74372cdeea92",
	}

	for i, address := range addresses {
		var account models.Account
		if err := db.Where("address = ?", address).First(&account).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Account does not exist, create a new one
				account = models.Account{
					AccountUuid: accountUuids[i],
					Address:     address,
				}
				if err := account.Save(db); err != nil {
					log.Printf("failed to seed account with address %s: %v", address, err)
				} else {
					log.Printf("successfully seeded account %s with AccountUuid %s", address, accountUuids[i])
				}
			} else {
				log.Printf("error checking account with address %s: %v", address, err)
			}
		} else {
			// Account exists, check if AccountUuid is correct
			if account.AccountUuid != accountUuids[i] {
				account.AccountUuid = accountUuids[i]
				if err := account.Save(db); err != nil {
					log.Printf("failed to update account with address %s: %v", address, err)
				} else {
					log.Printf("successfully updated account %s with AccountUuid %s", address, accountUuids[i])
				}
			} else {
				log.Printf("account %s already exists with AccountUuid %s", address, accountUuids[i])
			}
		}
	}
}
