package secondWeek

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

// InsertStmt 作业里面我们这个只是生成 SQL，所以在处理 sql.NullString 之类的接口
// 只需要判断有没有实现 driver.Valuer 就可以了
func InsertStmt(entity interface{}) (string, []interface{}, error) {

	if entity == nil {
		return "", nil, errInvalidEntity
	}

	typ := reflect.TypeOf(entity)
	refVal := reflect.ValueOf(entity)
	// 检测 entity 是否符合我们的要求
	// 我们只支持有限的几种输入
	if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct {
		return "", nil, errInvalidEntity
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		refVal = refVal.Elem()
	}
	// 只支持一级指针
	if typ.Kind() != reflect.Struct {
		return "", nil, errInvalidEntity
	}

	// 使用 strings.Builder 来拼接 字符串
	bd := strings.Builder{}
	bd.WriteString("INSERT INTO ")
	// 构造 INSERT INTO XXX，XXX 是你的表名，这里我们直接用结构体名字
	tableName := snakeString(typ.Name())
	bd.WriteString("`" + tableName + "`(")
	// 遍历所有的字段，构造出来的是 INSERT INTO XXX(col1, col2, col3)
	// 在这个遍历的过程中，你就可以把参数构造出来
	// 如果你打算支持组合，那么这里你要深入解析每一个组合的结构体
	// 并且层层深入进去
	getSqlFile(typ, &bd, 0, make(map[string]struct{}))

	// 拼接 VALUES，达成 INSERT INTO XXX(col1, col2, col3) VALUES

	// 再一次遍历所有的字段，要拼接成 INSERT INTO XXX(col1, col2, col3) VALUES(?,?,?)
	// 注意，在第一次遍历的时候我们就已经拿到了参数的值，所以这里就是简单拼接 ?,?,?
	var args []interface{}
	record := make(map[interface{}]interface{})
	getSqlValue(refVal, &bd, 0, record)
	for _, value := range record {
		args = append(args, value)
	}

	return bd.String(), args, nil
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

func getSqlFile(typ reflect.Type, sqlStr *strings.Builder, lv int, record map[string]struct{}) {

	filedLen := typ.NumField()
	for i := 0; i < filedLen; i++ {
		v := typ.Field(i)
		if v.Type.Kind() == reflect.Struct {
			// 查看是否是SQL.NullXXX结构体
			_, ok := v.Type.MethodByName("Value")
			if !ok {
				temp := lv + 1
				getSqlFile(v.Type, sqlStr, temp, record)
				continue
			}
		}

		filedName := snakeString(v.Name)
		if _, ok := record[filedName]; !ok {
			if ((i + 1) == filedLen) && (lv == 0) {
				sqlStr.WriteString("`" + filedName + "`) VALUES(")
			} else {
				sqlStr.WriteString("`" + filedName + "`,")
			}
			record[filedName] = struct{}{}
		}
	}

}

func getSqlValue(refVal reflect.Value, sqlStr *strings.Builder, lv int, record map[interface{}]interface{}) {

	filedLen := refVal.NumField()
	for i := 0; i < filedLen; i++ {
		var file string
		v := refVal.Field(i)
		switch v.Kind() {
		case reflect.Int64:
			if v.IsZero() {
				file = "0"
			} else {
				file = strconv.FormatInt(v.Int(), 10)
			}
			record[v.Addr()] = v.Interface()
		case reflect.Uint64:
			if v.IsZero() {
				file = "0"
			} else {
				file = strconv.FormatUint(v.Uint(), 10)
			}
			record[v.Addr()] = v.Interface()
		case reflect.Ptr:
			if v.IsZero() {
				// 这里暴力处理一下，只要是nil，全部转换成NULL
				file = "NULL"
			} else {
				_, ok := v.Elem().Type().MethodByName("Value")
				if !ok {
					file = strconv.FormatInt(v.Elem().Int(), 10)
				} else {
					if v.Elem().FieldByName("Valid").Bool() {
						file = getNullValue(v.Elem())
					} else {
						file = "NULL"
					}
				}
			}
			record[v.Addr()] = v.Interface()
		case reflect.Struct:
			// 查看是否是SQL.NullXXX结构体
			_, ok := v.Type().MethodByName("Value")
			if !ok {
				temp := lv + 1
				getSqlValue(v, sqlStr, temp, record)
				continue
			} else {
				if v.FieldByName("Valid").Bool() {
					file = getNullValue(v)
				} else {
					file = "NULL"
				}
				record[v.Addr()] = v.Interface()
			}
		case reflect.String:
			if v.IsZero() {
				file = ""
			} else {
				file = v.String()
			}
			record[v.Addr()] = v.Interface()
		}

		if ((i + 1) == filedLen) && (lv == 0) {
			sqlStr.WriteString(file + ");")
		} else {
			sqlStr.WriteString(file + ",")
		}

	}

}

func getNullValue(refVal reflect.Value) string {

	switch refVal.Field(0).Kind() {
	case reflect.String:
		return refVal.Field(0).String()
	case reflect.Int32:
		return strconv.FormatInt(refVal.Field(0).Int(), 10)
	default:
		return ""
	}

}
