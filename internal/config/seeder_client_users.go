package config

import (
	"dimensy-bridge/internal/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserActivationData struct {
	PsreID                 string
	NIK                    *string
	Name                   *string
	BirthDate              *time.Time
	Email                  *string
	PhoneNumber            *string
	IsWNI                  *bool
	IsBackendPhoneVerified bool
	IsKYCVerified          bool
}

func SeedClientUsersFromUserActivations(db *gorm.DB) error {
	// Connect to client database
	clientDB := ConnectClientDatabase()

	// Raw SQL query to fetch user activation data
	rawQuery := `
	select ua.psre_id,
	       udi.nik, 
	       udi."name",
	       udi.birth_date, 
	       u.email, 
	       udi.phone_number, 
	       true as is_wni, 
	       ua.is_backend_phone_verified as is_active,
	       ua.is_backend_phone_verified, 
	       ua.is_kyc_verified
	from user_activations ua 
	join users u on ua.user_id = u.id
	join (select max(id) id, user_id from user_details ud group by user_id) ud on ud.user_id = u.id
	join user_details udi on udi.id = ud.id
	`

	var userActivationDataList []map[string]interface{}
	if result := clientDB.Raw(rawQuery).Scan(&userActivationDataList); result.Error != nil {
		return fmt.Errorf("failed to fetch user activation data: %w", result.Error)
	}

	if len(userActivationDataList) == 0 {
		fmt.Println("⚠️ No user activation data found to seed")
		return nil
	}

	fmt.Printf("🔄 Processing %d user activation records...\n", len(userActivationDataList))

	// Convert fetched data to ClientUser models
	clientUsers := make([]model.ClientUser, 0, len(userActivationDataList))

	for _, row := range userActivationDataList {
		isWNI := true
		clientID := int64(1)
		isActive := false
		isVerifyPhone := false
		isVerifyKYC := false

		// Extract data with type assertions
		if psreID, ok := row["psre_id"].(string); ok {
			externalID := uuid.Nil // Initialize with nil UUID
			// Try to parse psreID as UUID, or generate one
			if parsedID, err := uuid.Parse(psreID); err == nil {
				externalID = parsedID
			} else {
				// Generate UUID from psreID if it's not a valid UUID
				externalID = uuid.NewSHA1(uuid.NameSpaceDNS, []byte(psreID))
			}

			nik := extractString(row["nik"])
			name := extractString(row["name"])
			birthDate := extractTime(row["birth_date"])
			email := extractString(row["email"])
			phone := extractString(row["phone_number"])

			// Extract boolean fields
			if isActiveVal, ok := row["is_active"].(bool); ok {
				isActive = isActiveVal
			}
			if isPhoneVerified, ok := row["is_backend_phone_verified"].(bool); ok {
				isVerifyPhone = isPhoneVerified
			}
			if isKYC, ok := row["is_kyc_verified"].(bool); ok {
				isVerifyKYC = isKYC
			}

			clientUser := model.ClientUser{
				ClientID:        clientID,
				ClientCompanyID: nil,
				ExternalID:      &externalID,
				NIK:             nik,
				Name:            name,
				Birthdate:       birthDate,
				Email:           email,
				Phone:           phone,
				IsWNI:           &isWNI,
				IsActive:        isActive,
				IsVerifyPhone:   isVerifyPhone,
				IsVerifyKYC:     isVerifyKYC,
				ParentID:        nil,
			}

			clientUsers = append(clientUsers, clientUser)
		}
	}

	// Insert or skip on conflict (using external_id as unique constraint)
	if len(clientUsers) > 0 {
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},
			DoNothing: true, // Skip insert if external_id already exists
		}).CreateInBatches(clientUsers, 100)

		if result.Error != nil {
			return fmt.Errorf("failed to seed client users: %w", result.Error)
		}

		fmt.Printf("✅ Seeded %d client users (some may have been skipped if already exist)\n", result.RowsAffected)
	}

	return nil
}

// Helper function to extract string value from interface{}
func extractString(val interface{}) *string {
	if val == nil {
		return nil
	}
	if str, ok := val.(string); ok {
		return &str
	}
	return nil
}

// Helper function to extract time.Time value from interface{}
func extractTime(val interface{}) *time.Time {
	if val == nil {
		return nil
	}
	if t, ok := val.(time.Time); ok {
		return &t
	}
	return nil
}
