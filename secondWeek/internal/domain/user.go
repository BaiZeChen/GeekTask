package domain

import (
	"time"
)

// User 领域对象，是 DDD 中的 entity
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Nickname string // 最多5个字符
	Birthday string
	Details  string
	Ctime    time.Time
}

//type Address struct {
//}
