package jwtserver

import (
	"errors"
	"github.com/leyle/fabric-app-starter/context"
)

const (
	UserRoleAdmin  = "admin"
	UserRoleClient = "user"
)

const (
	UserStatusInvalid = 0
	UserStatusValid   = 1
)

type UserAccount struct {
	Id        string `json:"id" gorm:"column:id"`
	Username  string `json:"username" gorm:"column:username;unique"`
	Role      string `json:"role" gorm:"column:role"`
	Salt      string `json:"-" gorm:"column:salt"`
	PassHash  string `json:"-" gorm:"column:passHash"`
	Valid     int    `json:"-" gorm:"column:passHash"`
	CreatedAt int64  `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt int64  `json:"updatedAt" gorm:"column:updatedAt"`
}

type CaAccount struct {
	// userId is UserPasswd table's id value
	UserId string `json:"userId" gorm:"column:userId;unique"`

	// password get from ca server or generate by api server
	// it saves as unencrypted
	Password string `json:"-" gorm:"column:password"`

	// how many times this client ca credential exported
	EnrollTimes int   `json:"enrollTimes" gorm:"column:enrollTimes"`
	CreatedAt   int64 `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   int64 `json:"updatedAt" gorm:"column:updatedAt"`
}

func GetUserByUsername(ds *context.Ds, username string) (*UserAccount, error) {
	var ua UserAccount
	tx := ds.Db.Where("username = ?", username).First(&ua)
	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return &ua, nil
}

func SaveUserAccount(ds *context.Ds, ua *UserAccount) error {
	tx := ds.Db.Create(ua)
	if tx.RowsAffected == 0 {
		return errors.New("Save UserAccount failed")
	}
	return nil
}
