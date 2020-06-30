package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 定义用户表
type User struct {
	UserID    string     `gorm:"primary_key;column:user_id;type:varchar(32);comment:'用户ID'"`
	UserName  string     `gorm:"unique;not null;column:user_name;type:varchar(50);comment:'用户名/登陆用户名'"`
	Email     string     `gorm:"unique;not null;column:email;type:varchar(50);comment:'邮箱/登陆邮箱'"`
	Password  string     `gorm:"not null;column:password;type:varchar(50);comment:'密码'"`
	Salt      string     `gorm:"not null;column:salt;type:varchar(32);comment:'混淆盐'"`
	CreatedAt time.Time  `gorm:"not null;comment:'注册时间'"`
	UpdatedAt time.Time  `gorm:"comment:'更新资料时间'"`
	DeletedAt *time.Time `gorm:"comment:'注销账户时间'"`
}

// 设置表名
func (u User) TableName() string {
	return "user"
}

// 创建初始化表
func initUser() {
	if !db.HasTable(&User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
			panic(err)
		}
	}
}

// 添加用户
func (u User) InsertUser() error {
	return db.Create(&u).Error
}

// 检查用户是否存在
func CheckUser(username, email string) (string, error) {
	var user User
	err := db.Where(User{UserName: username}).Or(User{Email: email}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	return user.UserID, nil
}

// 通过用户Id获取用户
func GetUserById(userId string) (*User, error) {
	var user User
	err := db.Where(User{UserID: userId}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if user.UserID != "" {
		return &user, nil
	}
	return nil, nil
}
