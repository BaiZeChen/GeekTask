package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
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

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, id).Error
	return u, err
}

func (dao *UserDAO) SaveDetails(ctx context.Context, user User) error {
	user.Utime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&user).Updates(user).Error
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

// User 直接对应数据库表结构
// 有些人叫做 entity，有些人叫做 model，有些人叫做 PO(persistent object)
type User struct {
	Id int64 `gorm:"column:id;primaryKey;autoIncrement"`
	// 全部用户唯一
	Email    string `gorm:"column:email;unique"`
	Password string `gorm:"column:password"`

	// 往这面加
	Nickname string `gorm:"column:nickname"`
	Birthday string `gorm:"column:birthday"`
	Details  string `gorm:"column:details"`

	// 创建时间，毫秒数
	Ctime int64 `gorm:"column:c_time"`
	// 更新时间，毫秒数
	Utime int64 `gorm:"column:u_time"`
}

func (u *User) TableName() string {
	return "user"
}
