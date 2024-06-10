package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(c context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()

	// SELECT * FROM users where email=123@qq.com FOR UPDATE
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(c).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}

	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateByUserId(ctx context.Context, u User) error {
	var us User
	err := dao.db.WithContext(ctx).Where("id = ?", u.Id).First(&us).Error
	if err != nil {
		return err
	}

	us.NickName = u.NickName
	us.AboutMe = u.AboutMe
	us.Birthday = u.Birthday
	err = dao.db.WithContext(ctx).Save(&us).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

// User 直接对应数据库表结构
// 有些人称为 entity，有些称为 model，有些称为PO（persistence object）
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement" `
	Email    string `gorm:"unique"`
	Password string

	// 补充的 profile 信息
	NickName string
	Birthday int64
	// 限定自我简介的长度
	AboutMe string `gorm:"type:varchar(1024)"`

	// 创建时间, 毫秒数
	Ctime int64
	// 更新时间
	Utime int64
}

type Address struct {
	Id     int64
	UserId int64
}
