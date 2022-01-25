package gadget

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func UUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

// func MD5(s string) string {
// 	signByte := []byte(s)
// 	hash := md5.New()
// 	hash.Write(signByte)
// 	return hex.EncodeToString(hash.Sum(nil))
// }

// RandString 生成随机字符串
func RandString(len int) string {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}

	return string(bytes)
}
