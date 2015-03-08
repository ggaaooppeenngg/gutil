package container

// 字符串排序算法的比较
// inplace 意思是说是否需要辅助函数转换
// 算法						stable  inplace	running	time	extra space	sweet spot
// insertion sort			yes		yes		N to N^2		1			small arrays,arrays in order
// quicksort				no		yes		N(logN)^2		logN		general-purpose when space is tight
// mergesort				yes		no		N(logN)^2		N			general-purpose stable srot
// 3-way quicksort			no		yes		N to NlogN		logN		large numbers of equal keys
// LSD string sort			yes		no		NW				N			short fixed-length strings
// MSD string sort			yes		no		N to Nw			N+WR		random strings
// 3-way string quicksort	no		yes		N to Nw			W+logN		general-purpose strings with long prefix matches

import ()

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

// 检查是否含有该字
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
func LSD(strings []string, w int) {
	// sort strings on leading len characters.
	n := len(strings)
	r := 256                 // radix
	aux := make([]string, n) // axuilary slice
	for d := w - 1; d >= 0; d-- {
		count := make([]int, r+1) // 后面要叠加,为了统一第一个空出来
		for i := 0; i < n; i++ {
			count[strings[i][d]+1]++
		}
		for i := 0; i < r; i++ {
			count[i+1] += count[i]
		}
		for i := 0; i < n; i++ {
			// 前面加1,这里没有加1是因为要靠第一个元素应该在前一个元素紧跟最后面
			// (比如第一个空元素是0,第一类元素是从0开始的,而不是count[1]对应的位置
			aux[count[strings[i][d]]] = strings[i]
			// 统计完一次加一次,分布完了以后就可以编程末尾的索引+1
			count[strings[i][d]]++
		}
		for i := 0; i < n; i++ {
			strings[i] = aux[i]
		}
	}
}

// Proof:
// LSD 算法依赖于key索引实现的稳定性,
// 对于i个尾部的keys,他们被排序要么是因为
// 第i个key不同,按照key排序了,
// 要么是相同,这个就是有序的了,这一定依靠
// key索引实现的稳定性,不能对应的索引不总是相同的.

// MSD(most significant first) 字符串排序.
// MSD 的特点是需要一个多余的字符判断结束,另外本身count数组变
// 索引数组又要多一个位置,所以整个的空间是R+2.
// 这个算法分到后面很浪费空间,有时候字符串完全是一样的,就在重复分配空间
func MSD(strs []string) {
	var (
		r   = 256                       //radix
		m   = 0                         // cutoff for small subarrays 小的子部分就不用基数排序不然开销很大.
		aux = make([]string, len(strs)) // auxiliary array for distribution.
	)
	sortMSD(strs, aux, 0, len(strs)-1, 0, m, r)
}

// sort strings from strs[lo] to strs[hi] at char d.
// m is cutoff for small subarrays.
// r is the alphabet radix.
func sortMSD(strs []string, aux []string, lo int, hi int, d int, m int, r int) {
	if hi <= lo+m {
		//小规模的划分不用基数排序
		return
	}
	// 这里比LSD多一个,是少于d的字符串统一在-1+2=1的位置,需要多出一个位置来记录小于的情况.
	var count = make([]int, r+2) // compute the frequency,one for computing the indices,one for end of string.
	// 统计频率的时候,
	// 基数组的用法:
	// count[1]是长度为d的字符串的数量
	// count[2]-count[R+1]是对应的字母表的数量.
	// count[0]没有用来计数,但是索引的时候作为开始起点(因为是0)
	for i := lo; i <= hi; i++ {
		count[charAt(strs[i], d)+2]++
	}
	// count[0]是起始索引,依次类推,要出现 R+1个
	// 1到R是原本的,0是小于d的字符串的.
	for i := 0; i < r+1; i++ {
		count[i+1] += count[i]
	}
	for i := lo; i <= hi; i++ {
		// 累加得出索引参考值之后 +2的地方是末尾,+1的地方是上一个的末尾,作为起点.
		k := count[charAt(strs[i], d)+1]
		aux[k] = strs[i] // -1+1=0表示
		// 统计完一次加一次,分布完了以后就可以变成末尾的索引+1,也可以是下个字母的开始.
		count[charAt(strs[i], d)+1]++
	}
	for i := lo; i <= hi; i++ {
		strs[i] = aux[i-lo]
	}
	for i := 0; i < r; i++ {
		sortMSD(strs, aux, lo+count[i], lo+count[i+1]-1, d+1, m, r)
	}
}

// returns -1 if it is end of string.
func charAt(s string, d int) int {
	if d < len(s) {
		return int(s[d])
	} else {
		return -1
	}
}

// 3-way string quicksort.
// 因为是递归调用,对于小的数组其实可以更换排序方式,
// 这样就可以提高效率,而且这也是常见的方法.
// 三路快排较于直接排序的好处就是,这些排序会把相同的前缀分类,
// 而不是直接比较每个字符.
// 随机化可以防止一些最坏现象,比如数组已经排好序或接近有序.
// 但是随机化也是需要产生随机数的代价的,这个需要权衡.
func Quick3string(strs []string) {
	sort3way(strs, 0, len(strs)-1, 0)
}

// internal implementation of 3-way string quicksort.
func sort3way(strs []string, lo int, hi int, d int) {
	if hi <= lo {
		return
	}
	lt := lo
	gt := hi
	v := charAt(strs[lo], d)
	// 每次划分三个区域,递归进行快排.
	// 这就是三路快排.
	for i := lo + 1; i <= gt; {
		t := charAt(strs[i], d)
		if t < v {
			strs[lt], strs[i] = strs[i], strs[lt]
			lt++
			i++
		}
		if t > v {
			strs[i], strs[gt] = strs[gt], strs[i]
			gt--
		}
		if t == v {
			i++
		}

	}

	// strs[lo:lt] < v = strs[lt,gt+1] < strs[gt+1:hi]
	sort3way(strs, lo, lt-1, d)
	if v >= 0 {
		sort3way(strs, lt, gt, d+1)
	}
	sort3way(strs, gt+1, hi, d)
}

// TODO: Tries tree.
