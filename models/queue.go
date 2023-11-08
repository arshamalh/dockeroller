package models

import "fmt"

type Node struct {
	Value string
	Next  *Node
}

func newNode(value string, next *Node) *Node {
	return &Node{
		Value: value,
		Next:  next,
	}
}

type Queue struct {
	Length int
	head   *Node
	tail   *Node
}

// Add to the end, remove from the beginning
func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) isZeroOrOneNode() (bool, *Node) {
	if q.Length == 0 {
		return true, nil
	} else if q.Length == 1 {
		value := q.head
		q.head = nil
		q.tail = nil
		q.Length--
		return true, value
	}
	return false, nil
}

func (q *Queue) Push(value string) *Queue {
	new_node := newNode(value, nil)
	if q.Length == 0 {
		q.head = new_node
	} else {
		q.tail.Next = new_node
	}
	q.tail = new_node
	q.Length++
	return q
}

func (q *Queue) Pop() *Node {
	if ok, node := q.isZeroOrOneNode(); ok {
		return node
	}
	node := q.head
	q.head = q.head.Next
	q.Length--
	return node
}

func (q *Queue) String() string {
	printable_list := ""
	for node := q.head; node != nil; node = node.Next {
		printable_list += fmt.Sprintln(node.Value)
	}
	return printable_list
}
