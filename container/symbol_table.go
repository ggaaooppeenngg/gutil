package container

// 符号表接口.
type SymbolTable interface {
	Put(key interface{}, value interface{})
	Get(key interface{}) interface{}
	// delete有lazy模式和eager模式.
	// lazy模式并不马上删除这个结点,只是把值置为null.
	Delete(key interface{})
	Contains(key interface{}) bool
	IsEmpty() bool
	Size() int
	Keys() []interface{}
}

// 天花板是大于里面的最小值.
// 而地板是小于里面的最大值.
// rank key 得到key的排序.
// select 7 得到排序为7的value.
