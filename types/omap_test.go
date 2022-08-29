package types

import (
	"testing"
)

func TestOMap(t *testing.T) {
	o := New()
	key1, key2, key3, key4 := "yumontime", "yot", "jankin", "henry"
	value := "nb"
	if !o.Exists(key1) {
		t.Error("1:", key1, "不存在")
	}
	if o.RPush(key1, value) {
		t.Log("key1写入成功")
	}
	if !o.Exists(key1) {
		t.Error("2:", key1, "不存在")
	}
	t.Log("插入key2前的omap大小", o.Size())
	o.RPush(key2, value)
	t.Log("RPush插入key2后的omap大小", o.Size())
	t.Log("RPush插入key2后的omap", o.dataMap)

	o.LPush(key3, value)
	o.LPush(key4, value)
	o.LWalk(func(key, value any) {
		t.Log(key, value)
	})
	for i := 0; i < 10; i++ {
		o.RPush(i, i)
	}
	// t.Log("开始测试截取区间---初始化 ")
	// o.LWalk(func(key, value any) {
	// 	t.Log(key, value)
	// })
	newO := o.Between(-4, -1)
	t.Log("开始测试截取区间--- 打印结果")
	newO.LWalk(func(key, value any) {
		t.Log(key, value)
	})
	t.Error("----")
}

func TestBetween(t *testing.T) {
	o := New()
	for i := 0; i < 10; i++ {
		o.RPush(i, i)
	}
	newO := o.Between(-4, -2)
	t.Log("开始测试截取区间--- 打印结果")
	newO.LWalk(func(key, value any) {
		t.Log(key, value)
	})
}

var o = New()

func BenchmarkRPush(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				o.RPush(i, i)
				i++
			}
		})
	}
}
