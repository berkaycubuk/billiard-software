package models

type OrderDetail struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	OrderID		uint		`json:"order_id"`
	UserName	string		`gorm:"type:varchar(255)" json:"user_name"`
	UserSurname	string		`gorm:"type:varchar(255)" json:"user_surname"`
	UserPhone	string		`gorm:"type:varchar(255)" json:"user_phone"`
	UserEmail	string		`gorm:"type:varchar(255)" json:"user_email"`
}
