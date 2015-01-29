package utils

// Nothing here is exported because these are only used by fileUtils.go

type list struct {
	root, last *node
	length     int
}

type node struct {
	value string
	next  *node
}

func (l *list) add(s string) {
	l.last.next = &node{value: s, next: nil}
	l.last = l.last.next
	l.length++
}
func newList() *list {
	a := list{}
	a.root = &node{value: "", next: nil}
	a.last = a.root
	return &a
}
func (l *list) len() int {
	return l.length
}
func (l *list) toStringArray() []string {
	s := make([]string, l.len())
	index := 0
	currentNode := l.root.next
	for currentNode != nil {
		s[index] = currentNode.value
		currentNode = currentNode.next
		index++
	}
	return s
}
