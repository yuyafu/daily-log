package trie

import "strings"

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}


type node struct {
	fullValue  []rune
	value      []rune
	sids       []int
	indicators []rune // children的所有第一个字,每一个字对应的下标就是孩子的下标。
	children   []*node
	priority   uint32 // 根据重复的前缀加权限，权限大的放在child的前面，提高性能更快的查到
}

func (n *node) AddValToSids(val []rune, sid int) {
	n.priority++

	if len(n.value) > 0 || len(n.children) > 0 {
	work:
		for {

			i := 0
			max := min(len(val), len(n.value))

			for i < max && val[i] == n.value[i] {
				i++
			}

			// 需要把重复的前缀提取出来，原有的作为孩子保存
			if i < len(n.value) {
				child := node{
					value:      n.value[i:],
					priority:   n.priority - 1,
					indicators: n.indicators,
					fullValue:  n.fullValue,
					sids:       n.sids,
					children:   n.children,
				}

				n.children = []*node{&child}
				n.indicators = []rune{n.value[i]}
				n.value = n.value[0:i]

			}

			n.sids = append(n.sids, sid)
			// 新输入的字符串
			if i < len(val) {

				val = val[i:]
				c := val[0]

				// beego.Error("val===",string(c),string(val), string(n.indicators))

				for i := 0; i < len(n.indicators); i++ {
					if c == n.indicators[i] {
						i = n.incrementChildPrio(i)
						n = n.children[i]
						continue work
					}
				}

				n.indicators = append(n.indicators, c)
				child := &node{}
				n.children = append(n.children, child)
				n.incrementChildPrio(len(n.indicators) - 1)
				n = child
				n.insertChild(val, sid)
				return

			}
		}

	} else {
		n.insertChild(val, sid)
	}
}

func (n *node) insertChild(val []rune, sid int) {

	n.value = val
	n.sids = append(n.sids, sid)
}

// increments priority of the given child and reorders if necessary
func (n *node) incrementChildPrio(pos int) int {
	n.children[pos].priority++
	prio := n.children[pos].priority

	// adjust position (move to front)
	newPos := pos
	for newPos > 0 && n.children[newPos-1].priority < prio {
		// swap node positions
		n.children[newPos-1], n.children[newPos] = n.children[newPos], n.children[newPos-1]

		newPos--
	}

	// build new index char string
	if newPos != pos {
		var tmp []rune
		tmp = append(tmp, n.indicators[:newPos]...)
		tmp = append(tmp, n.indicators[pos])
		tmp = append(tmp, n.indicators[newPos:pos]...)
		tmp = append(tmp, n.indicators[pos+1:]...)
		n.indicators = tmp
	}
	return newPos
}

func (n *node) GetNodeByVal(val []rune) (res *node) {
	if len(val) == 0 {
		return
	}

	fullName := []rune(n.value)
	res = &node{}

	work:
	for {

		if len(n.value) != 0 {
			// 与当前的val包含，当前的val，不能为空
			minLen := min(len(val),len(n.value))
			i := 0
			for i < minLen && val[i] == n.value[i] {
				i++
			}

			if strings.Contains(string(n.value), string(val)) {
				res = n
				res.fullValue = fullName
				return
			} else {
				val = val[i:]
			}
		}

		// 直接找孩子
		c := val[0]
		for i := 0; i < len(n.indicators); i++ {
			if c == n.indicators[i] {
				n = n.children[i]
				fullName = append(fullName, n.value...)
				continue work
			}
		}
		return
	}
	return
}
