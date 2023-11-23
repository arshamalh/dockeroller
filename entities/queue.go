package entities

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

// Removes a Node from the queue and returns it,
// Returns an error if there is no available node
func (q *Queue) Pop() (string, error) {
	if q.Length == 0 {
		return "", fmt.Errorf("no node available to pop")
	}

	node := q.head
	if q.Length == 1 {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.Next
	}

	q.Length--
	return node.Value, nil
}

func (q *Queue) String() string {
	printable_list := ""
	for node := q.head; node != nil; node = node.Next {
		printable_list += fmt.Sprintln(node.Value)
	}
	return printable_list
}
