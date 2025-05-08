package main

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

func parseExpire(expire string) (string, error) {
	s := strings.TrimSpace(expire)
	if s == "-1" || s == "infinite" || s == "forever" {
		return "-1", nil
	}

	// 年月日时分秒
	_, err := time.Parse("20060102150405", s)
	if err == nil {
		return s, nil
	}

	// 年月日时分
	t, err := time.Parse("200601021504", s)
	if err == nil {
		return t.Format("20060102150405"), nil
	}

	// 年月日时
	t, err = time.Parse("2006010215", s)
	if err == nil {
		return t.Format("20060102150405"), nil
	}

	// 年月日
	t, err = time.Parse("20060102", s)
	if err == nil {
		return t.Format("20060102150405"), nil
	}

	// 3d5h etc.
	d, err := time.ParseDuration(s)
	if err == nil {
		return time.Now().Add(d).Format("20060102150405"), nil
	}

	return "", err
}

func GenerateAPIKey(length int) (string, error) {
	// 每个字节转换为两个十六进制字符，因此需要 length/2 字节
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
