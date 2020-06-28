package util

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成随机数字
func GetRandomCode(length int) string {
	var container string
	for i := 0; i < length; i++ {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		container += fmt.Sprintf("%01v", rnd.Int31n(10))
	}
	return container
}

//生成UUID
func GetUUIDString(is bool) string {
	if is {
		return fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil))
	}
	return strings.Replace(fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil)), "-", "", -1)
}
