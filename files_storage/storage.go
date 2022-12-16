package storage

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	UserID         int64
	ChatID         int
	DateTill       time.Time
	WarningMessage string
}

type Payments struct {
	Time              time.Time
	UserID            int64
	Username          string
	Amount            int
	Product           string
	TelegramPaymentID string
	ProviderPaymentID string
}

type DataBase struct {
	*gorm.DB
}

func NewDB() DataBase {
	db, err := gorm.Open(sqlite.Open("dataBase.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return DataBase{db}
}

func (db *DataBase) MigrateDBs() {
	log.Println("migrating database User")
	db.AutoMigrate(&User{})
	log.Println("migrating database Payments")
	db.AutoMigrate(&Payments{})
}

func (db *DataBase) NewPaymentRecord(amount int, userID int64, username, product, telegramPaymentID, providerPaymentID string) {
	log.Println("new payment record in DB")
	db.Table("payments").Create(Payments{
		Time:              time.Now().UTC(),
		UserID:            userID,
		Username:          username,
		Amount:            amount / 100,
		Product:           product,
		TelegramPaymentID: telegramPaymentID,
		ProviderPaymentID: providerPaymentID,
	})
}

func (db *DataBase) UpdateSubscriptionTime(userID int64, chatID int, product string) {
	usr := db.FindUser(userID)
	isVal, _ := db.CheckSubscription(userID)
	switch product {
	case "1monthSub":
		if isVal {
			db.Table("users").Where("user_id = ?", userID).Updates(User{WarningMessage: " ", DateTill: usr.DateTill.Add(time.Hour * 24 * 30)})
		}
		if !isVal {
			db.Table("users").Where("user_id = ?", userID).Updates(User{WarningMessage: " ", DateTill: time.Now().UTC().Add(time.Hour * 24 * 30)})
		}
	case "3monthsSub":
		if isVal {
			db.Table("users").Where("user_id = ?", userID).Updates(User{WarningMessage: " ", DateTill: usr.DateTill.Add(time.Hour * 24 * 90)})
		}
		if !isVal {
			db.Table("users").Where("user_id = ?", userID).Updates(User{WarningMessage: " ", DateTill: time.Now().UTC().Add(time.Hour * 24 * 90)})
		}
	case "6monthsSub":
		if isVal {
			db.Table("users").Where("user_id = ?", userID).Updates(User{WarningMessage: " ", DateTill: usr.DateTill.Add(time.Hour * 24 * 180)})
		}
		if !isVal {
			db.Table("users").Where("user_id = ?", userID).Updates(User{WarningMessage: " ", DateTill: time.Now().UTC().Add(time.Hour * 24 * 180)})
		}
	}
}

func (db *DataBase) CheckSubscription(userID int64) (isValid bool, validTill time.Time) {
	usr := db.FindUser(userID)
	validSubscription := usr.DateTill.After(time.Now().UTC())
	return validSubscription, usr.DateTill
}

func (db *DataBase) FindUser(userID int64) *User {
	var user User
	db.Table("users").First(&user, "user_id = ?", userID)
	return &user
}

func (db *DataBase) GetAllUsers() []*User {
	users := make([]*User, 0)
	db.Table("users").Find(&users)
	return users
}

func (db *DataBase) CreateNewUser(userID int64, chatID int) {
	db.Table("users").Create(&User{
		UserID:         userID,
		ChatID:         chatID,
		DateTill:       time.Now().UTC().AddDate(0, 0, -1),
		WarningMessage: "ended",
	})
}
