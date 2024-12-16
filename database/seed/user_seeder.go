package seed

import (
	"bwanews/internal/core/domain/model"
	"bwanews/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB){
	// cypher Password
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatal().Err(err).Msg("Error	creating password hash")
	}

	// Create account for User
	admin := model.User{
		Name: "Admin",
		Email: "admin@gmail.com",
		Password: string(bytes),
	}

	//? Jika ada Akun admin console "Error seeding admin role" klo gk dia buat akun baru
	//? Cek apakah akun sudah ada
	var existingUser model.User
	if err := db.Where("email= ?", "admin@gmail.com").First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//! Akun tidak di temukan, buat Baru
			if err := db.Create(&admin).Error; err != nil {
				log.Fatal().Err(err).Msg("Error Creating admin Role")
			} else {
				log.Info().Msg("Admin role seeded successfully")
			}
		} else {
			log.Fatal().Err(err).Msg("Error checking admin role")
		}	
	} else {
		log.Info().Msg("Admin account already exists")
	}
	 


	// if err := db.FirstOrCreate(&admin,model.User{Email: "admin@gmail.com"}).Error; err != nil {
	// 	log.Fatal().Err(err).Msg("Error seeding admin role")
	// } else {
	// 	log.Info().Msg("Admin role seeded Successfully")
	// } 
 

}