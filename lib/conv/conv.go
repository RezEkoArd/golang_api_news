package conv

import (
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Fungsi Hashed Password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}


// CheckPassword Hash
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// String
func StringToInt64(s string)(int64, error){
	newData, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0 , err
	}

	return newData, nil
}


// Generate Slug
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug," ","-")

	return slug	
}

// Convert String to Integer
func StringToInt(s string) (int, error) {
	numb, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return numb, err
}

