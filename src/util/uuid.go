package util

import (
	"log"

	"github.com/google/uuid"
)

// 生成一个临时的uuid
func GenerageUUID() string {
	id := uuid.New()
	log.Println("generate uuid : ", id)
	return id.String()
}
