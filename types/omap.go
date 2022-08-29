package types

import (
	"container/list"
	"errors"
	"sync"
)

// 在这个基础上，可以实现sort方法，基于key做降序、升序的处理
// 初衷是能实现一个并发安全的高效的key-value格式的map。同时能结合一些方便的方法
// 我们希望一切的操作都是基于key-value格式的

type Elements struct {
	Value   any
	element *list.Element
}

type OMap struct {
	dataMap  map[any]*Elements
	dataList *list.List
	lock     *sync.RWMutex
}

func New() *OMap {
	return &OMap{
		dataMap:  make(map[any]*Elements),
		dataList: list.New(),
		lock:     &sync.RWMutex{},
	}
}

// Exists 判断一个key存不存在
func (o *OMap) Exists(key any) bool {
	_, ok := o.dataMap[key]
	return ok
}

// Push 往omap的右侧追加一个元素，即尾部
func (o *OMap) RPush(key, value any) bool {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.Exists(key) {
		return false
	}
	e := o.dataList.PushBack(key)
	o.dataMap[key] = &Elements{
		// Key:     key,
		Value:   value,
		element: e,
	}
	return true
}

// Push 往omap的左侧追加一个元素，即头部
func (o *OMap) LPush(key, value any) bool {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.Exists(key) {
		return false
	}
	e := o.dataList.PushFront(key)
	o.dataMap[key] = &Elements{
		// Key:     key,
		Value:   value,
		element: e,
	}
	return true
}

// Remove 移除特定的key
func (o *OMap) Remove(key any) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if !o.Exists(key) {
		return
	}
	o.dataList.Remove(o.dataMap[key].element)
	delete(o.dataMap, key)
}

// Size 判断大小
func (o *OMap) Size() int {
	o.lock.RLock()
	defer o.lock.RUnlock()

	return o.dataList.Len()
}

// Get 根据key 获取value
func (o *OMap) Get(key any) any {
	o.lock.RLock()
	defer o.lock.RUnlock()
	v, ok := o.dataMap[key]
	if ok {
		return v.Value
	}

	return nil
}

// Walk 左侧开始遍历迭代omap，即从左往右
func (o *OMap) LWalk(cb func(key, value any)) {
	for elem := o.dataList.Front(); elem != nil; elem = elem.Next() {
		cb(elem.Value, o.dataMap[elem.Value].Value)
	}
}

// Walk 左侧开始遍历迭代omap，即从左往右
func (o *OMap) RWalk(cb func(key, value any)) {
	for elem := o.dataList.Back(); elem != nil; elem = elem.Next() {
		cb(elem.Value, o.dataMap[elem.Value].Value)
	}
}

// MoveToFront 移动到最前面
func (o *OMap) MoveToFront(key any) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if !o.Exists(key) {
		return errors.New("the key not exist")
	}
	o.dataList.MoveToFront(o.dataMap[key].element)
	return nil
}

// MoveToBack 移动到最后面
func (o *OMap) MoveToBack(key any) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if !o.Exists(key) {
		return errors.New("the key not exist")
	}
	o.dataList.MoveToBack(o.dataMap[key].element)
	return nil
}

// MoveBefore 移动到指定key 之前
func (o *OMap) MoveBefore(key, mark any) error {
	o.lock.Lock()
	defer o.lock.Unlock()
	if !o.Exists(key) {
		return errors.New("the key not exist")
	}
	if !o.Exists(mark) {
		return errors.New("the mark not exist")
	}
	o.dataList.MoveBefore(o.dataMap[key].element, o.dataMap[mark].element)
	return nil
}

// MoveAfter 移动到指定key 之后
func (o *OMap) MoveAfter(key, mark any) error {
	if key == mark {
		return nil
	}
	o.lock.Lock()
	defer o.lock.Unlock()

	if !o.Exists(key) {
		return errors.New("the key not exist")
	}
	if !o.Exists(mark) {
		return errors.New("the mark not exist")
	}
	o.dataList.MoveAfter(o.dataMap[key].element, o.dataMap[mark].element)
	return nil
}

// Between 获取指定范围的value
// 非0的情况下  s <= e
// e != 0
func (o *OMap) Between(s, e int) *OMap {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if e == 0 || (s*e < 0 && s < e) || (s*e > 0 && s > e) {
		return nil
	}
	l := o.dataList.Len() // 总长度
	// 根据s,e的值，处理成非负数
	if s < 0 {
		s = s + l
	}
	if e < 0 {
		e = e + l
	}
	omap := New()
	i := 0
	for elem := o.dataList.Front(); elem != nil; elem = elem.Next() {
		if i >= s && i <= e {
			omap.RPush(elem.Value, o.dataMap[elem.Value].Value)
		}
		i++
	}
	return omap
}
