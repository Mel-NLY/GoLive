package pkgs

import (
	"errors"
	"fmt"
)

//BinaryNode is a struct
type BinaryNode struct {
	item  Post        // to store the data item
	left  *BinaryNode // pointer to point to left node
	right *BinaryNode // pointer to point to right node
}

//BST is a struct
type BST struct {
	root *BinaryNode
}

//Post is exported
type Post struct {
	ID       string
	Username string
	Time     struct {
		Day   int
		Month int
		Year  int
		Hour  int
		Min   int
	}
	Title       string
	Image       string
	Description string
	Tag         string
}

func (bst *BST) insertNode(t **BinaryNode, item Post) error {

	if *t == nil {
		newNode := &BinaryNode{
			item:  item,
			left:  nil,
			right: nil,
		}
		*t = newNode
		return nil
	}

	if item.Username < (*t).item.Username {
		bst.insertNode(&((*t).left), item)
	} else {
		bst.insertNode(&((*t).right), item)
	}

	return nil
}

//Insert is exported
func (bst *BST) Insert(item Post) {
	bst.insertNode(&bst.root, item)
}

var tempPost []Post
func (bst *BST) inOrderTraverse(t *BinaryNode) {
	if t != nil {
		bst.inOrderTraverse(t.left)
		tempPost = append(tempPost, t.item)
		bst.inOrderTraverse(t.right)
	}
}

func (bst *BST) preOrderTraverse(t *BinaryNode, U string) {
	if t != nil {
		if t.item.Username == U{
			tempPost = append(tempPost, t.item)
		}
		bst.preOrderTraverse(t.left, U)
		bst.preOrderTraverse(t.right, U)
	}
}

func (bst *BST) postOrderTraverse(t *BinaryNode) {
	if t != nil {
		bst.postOrderTraverse(t.left)
		bst.postOrderTraverse(t.right)
		fmt.Println(t.item)
	}
}

func (bst *BST) InOrder() []Post {
	tempPost = []Post{}
	bst.inOrderTraverse(bst.root)
	return tempPost
}

func (bst *BST) PreOrder(U string) []Post {
	tempPost = []Post{}
	bst.preOrderTraverse(bst.root, U)
	return tempPost
}

func (bst *BST) PostOrder() {
	bst.postOrderTraverse(bst.root)
}

func (bst *BST) searchNode(t **BinaryNode, x string) (Post, error) {
	if *t == nil {
		fmt.Println("Binary tree is empty.")
		return Post{}, errors.New("BST is empty")
	}
	if x < (*t).item.Username {
		bst.searchNode(&((*t).left), x)
	} else if x > (*t).item.Username {
		bst.searchNode(&((*t).right), x)
	} else {
		fmt.Println("Found!")
		return (*t).item, nil
	}
	fmt.Println("Not Found.")
	return Post{}, errors.New("Item not found")
}

//Search is exported
func (bst *BST) Search(x string) (Post, error) {
	return bst.searchNode(&bst.root, x)
}

func (bst *BST) getItem(t *BinaryNode) Post {
	if t.right != nil {
		bst.getItem(t.right)
	}
	return t.item
}

func (bst *BST) removeNode(t *BinaryNode, item Post) (*BinaryNode, error) { // Relook at it again, understand the code
	if t == nil {
		return nil, errors.New("Empty Binary Search Tree")
	} else if item.Username < t.item.Username {
		bst.removeNode(t.left, item)
	} else if item.Username > t.item.Username {
		bst.removeNode(t.right, item)
	} else {
		if t.left == nil {
			return t.right, nil
		} else if t.right == nil {
			return t.left, nil
		} else {
			t.item = bst.getItem(t.left)
			t.left, _ = bst.removeNode(t.left, t.item)
		}
	}
	return nil, errors.New("Not found")
}

//Remove is exported
func (bst *BST) Remove(item Post) {
	bst.removeNode(bst.root, item)
}

func (bst *BST) countNodes(t *BinaryNode) int {
	count := 1
	if t.left != nil {
		count += bst.countNodes(t.left)
	}
	if t.right != nil {
		count += bst.countNodes(t.right)
	}
	return count
}

//Count is exported
func (bst *BST) Count() int {
	count := 0
	if bst.root != nil {
		count = bst.countNodes(bst.root)
	}
	return count
}
