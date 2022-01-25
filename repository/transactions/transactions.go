package transactions

import (
	"fmt"
	"project-airbnb/entities"
	"time"

	"gorm.io/gorm"
)

type TransactionsRepository struct {
	db *gorm.DB
}

func NewTransactionsRepo(db *gorm.DB) *TransactionsRepository {
	return &TransactionsRepository{db: db}
}

func (tr *TransactionsRepository) Gets(userID uint) ([]entities.Transaction, error) {
	transaction := []entities.Transaction{}
	bookings := []entities.Book{}
	tr.db.Where("user_id=?", userID).Find(&bookings)
	tr.db.Joins("JOIN books ON books.transaction_id=transactions.id").Where("user_id=?", userID).Find(&transaction)
	return transaction, nil
}

func (tr *TransactionsRepository) Get(userID uint) ([]entities.Transaction, error) {
	transaction := []entities.Transaction{}
	bookings := []entities.Book{}
	tr.db.Where("user_id=?", userID).Find(&bookings)
	tr.db.Joins("JOIN books ON books.transaction_id=transactions.id").Where("user_id=? AND status='PENDING'", userID).Find(&transaction)
	return transaction, nil
}

func (tr *TransactionsRepository) Update(invoiceID string) (entities.Transaction, error) {
	transactionUpdate := entities.Transaction{}
	bookUpdate := entities.Book{}
	roomUpdate := entities.Room{}

	tr.db.Where("invoice=?", invoiceID).Find(&transactionUpdate)

	tr.db.Where("transaction_id=?", transactionUpdate.ID).Find(&bookUpdate)

	tr.db.Where("user_id=?", bookUpdate.User_id).Find(&roomUpdate)

	var now = time.Now()

	newBook := entities.Book{
		Checkin:  fmt.Sprint(now.Year(), "-", now.Month(), "-", now.Day()),
		Checkout: fmt.Sprint(now.Year(), "-", now.Month(), "-", now.Day()+roomUpdate.Duration),
	}
	tr.db.Where("id=?", bookUpdate.ID).Model(&bookUpdate).Updates(newBook)

	newRoom := entities.Room{
		Status: "CLOSED",
	}

	tr.db.Where("user_id=?", bookUpdate.User_id).Model(&roomUpdate).Updates(newRoom)

	transactionUpdate.Status = "LUNAS"
	tr.db.Save(&transactionUpdate)

	return transactionUpdate, nil
}
