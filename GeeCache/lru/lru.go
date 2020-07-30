package lru

import "container/list"   // 直接使用 Go 语言标准库实现的双向链表list.List


type Cache struct {
	maxBytes   int64                              // 允许使用的最大内存
	curBytes   int64                              // 当前已经使用的内存
	linkedList *list.List
	cache      map[string]*list.Element           // 键是字符串，值是双向链表中对应节点的指针
	OnEvicted  func(key string, value Value)      // 某条记录被移除时的回调函数，可以为 nil
}

// 双向链表节点的数据类型, 在节点中仍保存每个值对应的 key 的好处在于, 淘汰队首节点时, 需要用 key 从字典中删除对应的映射.
type entry struct {
	key   string
	value Value
}

// 为了通用性, 我们允许值是实现了 Value 接口的任意类型, 该接口只包含了一个方法 Len() int, 用于返回值所占用的内存大小.
type Value interface {
	Len() int
}

/*
	@func 实例化 Cache
 */
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  	maxBytes,
		linkedList: list.New(),
		cache:      make(map[string]*list.Element),
		OnEvicted:  onEvicted,
	}
}


/*
	@func 新增/修改功能
*/
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(ele)        // 如果键存在, 则更新对应节点的值, 并将该节点移到队首.（双向链表作为队列, 队首队尾是相对的, 在这里约定 front 为队首)
		kv := ele.Value.(*entry)
		c.curBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.linkedList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.curBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.curBytes {  // 超过了设定的最大值 maxBytes，则移除最近最少访问的节点
		c.RemoveOldest()
	}
}

/*
	@func 查找功能 从字典中找到对应的双向链表的节点, 将该节点移动到队首
*/
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(ele)  // 即将链表中的节点 ele 移动到队首
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

/*
	@func 删除功能 缓存淘汰, 即移除最近最少访问的节点（队尾）
*/
func (c *Cache) RemoveOldest() {
	ele := c.linkedList.Back()  // 返回列表l的最后一个元素, 如果列表为空，则返回nil
	if ele != nil {
		c.linkedList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.curBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.linkedList.Len()
}