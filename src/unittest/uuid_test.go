package unittest

import (
	"INServer/src/common/uuid"
	"testing"
)

func TestUUID(t *testing.T) {
	uuids := make(map[string]bool)
	// 测试一千万次循环uuid是否会重复
	for i := 0; i < 10000000; i++ {
		uid := uuid.New()
		if _, ok := uuids[uid]; ok {
			t.Fail()
			return
		} else {
			uuids[uid] = true
		}
	}
}
