package pkgs

import (
	"fmt"
)

//RoutePoint is exported
type RoutePoint struct {
	ID      string
	RouteID string
	Lat     float64
	Lon     float64
	Time    struct {
		Day   int
		Month int
		Year  int
		Hour  int
		Min   int
	}
	Position int
}

//Node is a struct
type Node struct {
	Item RoutePoint
	Next *Node
}

type LinkedList struct {
	Head *Node
	Size int
}

func (p *LinkedList) AddNode(rp RoutePoint) error {
	newNode := &Node{
		Item: rp,
		Next: nil,
	}
	if p.Head == nil {
		p.Head = newNode
	} else {
		currentNode := p.Head
		for currentNode.Next != nil {
			currentNode = currentNode.Next
		}
		currentNode.Next = newNode
	}
	p.Size++
	return nil
}

func (p *LinkedList) Remove(position int) error {
	if p.Head == nil {
		return fmt.Errorf("List is empty")
	}
	if position == 0 {
		p.Head = p.Head.Next
		p.Size--
		return nil
	}
	pointer := p.Head
	index := 0
	if position > 0 && position <= p.Size {
		for pointer.Next != nil {
			if index == position-1 {
				pointer.Next = pointer.Next.Next
				p.Size--
				return nil
			}
			pointer = pointer.Next
			index++
		}
	}
	return fmt.Errorf("Index is out of bounds")
}

func (p *LinkedList) AddAtPos(rp RoutePoint, position int) error {
	addNode := Node{rp, nil}
	if position == 0 {
		addNode.Next = p.Head
		p.Head = &addNode
		p.Size++
		return nil
	}
	pointer := p.Head
	index := 0
	if position > 0 && position <= p.Size {
		for pointer.Next != nil {
			if index == position-1 {
				addNode.Next = pointer.Next
				pointer.Next = &addNode
				p.Size++
				return nil
			}
			index++
			pointer = pointer.Next
		}
	}
	return fmt.Errorf("Index is out of bounds")
}

func (p *LinkedList) Get(position int) (RoutePoint, error) {
	if p.Head == nil {
		return RoutePoint{}, fmt.Errorf("List is empty")
	}
	if position == 0 {
		return p.Head.Item, nil
	}
	pointer := p.Head
	index := 0
	if position > 0 && position <= p.Size {
		for pointer.Next != nil {
			if index == position {
				return pointer.Item, nil
			}
			pointer = pointer.Next
			index++
		}
	}
	return RoutePoint{}, fmt.Errorf("Index is out of bounds")
}

func (p *LinkedList) Reverse() error {
	currNode := p.Head
	var nextNode *Node
	var prevNode *Node

	for currNode != nil {
		nextNode = currNode.Next
		currNode.Next = prevNode
		prevNode = currNode
		currNode = nextNode
	}
	p.Head = prevNode
	return nil
}

func (p *LinkedList) PrintAllNodes() error {
	currentNode := p.Head
	if currentNode == nil {
		fmt.Println("Linked list is empty.")
		return nil
	}
	fmt.Printf("%+v\n", currentNode.Item)
	for currentNode.Next != nil {
		currentNode = currentNode.Next
		fmt.Printf("%+v\n", currentNode.Item)
	}
	return nil
}
