package model

import "time"

type Content struct {
	ID 				int64 		`gorm:"id"`
	Title 			string 		`gorm:"title"`
	Excerpt 		string 		`gorm:"excerpt"`
	Description 	string 		`gorm:"description"`
	Image 			string 		`gorm:"image"`
	Tags 			string 		`gorm:"tags"`
	Status			string		`gorm:"status"`
	CategoryID		int64		`gorm:"category_id"`
	CreatedByID 	int64 		`gorm:"created_by_id"`
	User 			User 		`gorm:"foreignKey:CreatedByID"` //Untuk Foreign key to User
	Category 			Category 		`gorm:"foreignKey:CategoryId"` //Untuk Foreign key to User
	CreatedAt 		time.Time 	`gorm:"created_at"`
	UpdatedAt 		*time.Time 	`gorm:"updated_at"`
}