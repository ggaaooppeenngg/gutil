package container

// 代表一个字符集合的字母表
// 使用字母表作映射可以使得代码更加紧凑,因为对应的数字可以直接作为下标
// 有时候连图也不要,直接一个数组就可以解决.
// 当然并不是所有的字符串算法都要使用这个字母表,
// 转换成索引的过程的复杂度和整个算法的结构也需要考量.
type Alphabet struct {
}

// 根据索引得到字符
func (apb Alphabet) ToByte(i int) byte {
	return '1'
}

// 根据字符得到索引
func (apb Alphabet) ToIndex(c byte) int {
	return 1
}

// 检查是否含有该字符
func (apb Alphabet) Contains(c byte) bool {
	return false
}

// radix (number of characters in alphabet)
// 字符集合的大小
func (apb Alphabet) R() int {
	return 0
}

// number of bits to represent an index
// 用来表示所有的字符需要的索引的位数
func (apb Alphabet) lgR() int {
	return 0
}

// 把字符集合转换成索引集合.
func (apb Alphabet) ToIndices(s string) []int {
	return []int{}
}

// 把索引数组转换成字符数组
func (apb Alphabet) ToString(indices []int) string {
	return ""
}

// 俗称的基数排序中的LSD.
// Key作为索引的排序,先用计数器统计每个key的个数,然后依次决定每一段key的位置,再把对应的元素分配过去
// 就能按照key作为索引的排序了. aux 可以用来表示一些辅助工具.空间规模是(8N+3R+1) R是字母表的基数.

// LSD(leadt-significant-digit first) 字符串排序,也可以用于整数的排序.
func LSD(strings []string, w int) []string {
	// sort strings on leading len characters.
	n := len(strings)
	r := 256                 // radix
	aux := make([]string, n) // axuilary slice
	for d := w - 1; d >= 0; d-- {
		count := make([]int, r+1) //后面要叠加,为了统一第一个空出来
		for i := 0; i < n; i++ {
			count[strings[i][d]+1]++
		}
		for i := 0; i < r; i++ {
			count[i+1] += count[i]
		}
		for i := 0; i < n; i++ {
			aux[count[strings[i][d]]] = strings[i]
			count[strings[i][d]]++
		}
		for i := 0; i < n; i++ {
			strings[i] = aux[i]
		}
	}
	return nil
}

// Proof:
// LSD 算法依赖于key索引实现的稳定性,
// 对于i个尾部的keys,他们被排序要么是因为
// 第i个key不同,按照key排序了,
// 要么是相同,这个就是有序的了,这一定依靠
// key索引实现的稳定性,不能对应的索引不总是相同的.
