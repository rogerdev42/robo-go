package bst

import (
	"fmt"
	"strings"
)

type BinarySearchTree struct {
	root *node
}

type node struct {
	key    string
	left   *node
	right  *node
	parent *node
}
 
func NewBST() *BinarySearchTree {
	return &BinarySearchTree{}
}

func (t *BinarySearchTree) Insert(key string) error {
	if key == "" {
		return fmt.Errorf("empty key is not allowed")
	}

	if t.root == nil {
		t.root = &node{key: key}
		return nil
	}

	return t.root.insert(key)
}

func (n *node) insert(key string) error {
	if key == "" {
		return fmt.Errorf("empty key is not allowed")
	}

	if n.key == "" {
		n.key = key
		return nil
	}
	if strings.Compare(key, n.key) < 0 {
		if n.left != nil {
			err := n.left.insert(key)
			if err != nil {
				return err
			}
		} else {
			n.left = &node{
				key:    key,
				parent: n,
			}
		}
	} else {
		if n.right != nil {
			err := n.right.insert(key)
			if err != nil {
				return err
			}
		} else {
			n.right = &node{
				key:    key,
				parent: n,
			}
		}
	}
	return nil
}

func (t *BinarySearchTree) Find(key string) (string, error) {
	if t.root == nil {
		return "", fmt.Errorf("tree is empty")
	}
	node, err := t.root.find(key)
	if err != nil {
		return "", err
	}
	return node.key, nil
}

func (n *node) find(key string) (*node, error) {
	if n == nil {
		return nil, fmt.Errorf("key %s not found", key)
	}
	comp := strings.Compare(key, n.key)
	if comp == 0 {
		return n, nil
	}

	if comp < 0 {
		if n.left == nil {
			return nil, fmt.Errorf("key %s not found", key)
		}

		return n.left.find(key)
	}
	if n.right == nil {
		return nil, fmt.Errorf("key %s not found", key)
	}

	return n.right.find(key)
}

func (t *BinarySearchTree) Delete(key string) error {
	if t.root == nil {
		return fmt.Errorf("tree is empty")
	}
	return t.root.delete(key)
}

func (n *node) delete(key string) error {
	if n == nil {
		return fmt.Errorf("node not found")
	}

	comp := strings.Compare(key, n.key)
	if comp < 0 {
		return n.left.delete(key)
	} else if comp > 0 {
		return n.right.delete(key)
	}
	return n.removeNode()
}

func (n *node) removeNode() error {
	if n.left == nil && n.right == nil {
		if n.parent == nil {
			n.key = ""
			return nil
		}
		if n.parent.left == n {
			n.parent.left = nil
		} else {
			n.parent.right = nil
		}
		return nil
	}

	if n.right == nil {
		if n.parent.left == n {
			n.parent.left = n.left
		} else {
			n.parent.right = n.left
		}
		n.left.parent = n.parent
	} else if n.left == nil {
		if n.parent.left == n {
			n.parent.left = n.right
		} else {
			n.parent.right = n.right
		}
		n.right.parent = n.parent
	} else {
		minNode, err := n.right.min()
		if err != nil {
			return err
		}
		n.key = minNode.key
		return n.right.delete(minNode.key)
	}
	return nil
}

func (t *BinarySearchTree) Min() (string, error) {
	if t.root == nil {
		return "", fmt.Errorf("tree is empty")
	}
	node, err := t.root.min()
	if err != nil {
		return "", err
	}
	return node.key, nil
}

func (n *node) min() (*node, error) {
	if n == nil {
		return nil, fmt.Errorf("cannot find minimum in nil node")
	}

	if n.left == nil {
		return n, nil
	}
	return n.left.min()
}

func (t *BinarySearchTree) Max() (string, error) {
	if t.root == nil {
		return "", fmt.Errorf("tree is empty")
	}
	node, err := t.root.max()
	if err != nil {
		return "", err
	}
	return node.key, nil
}

func (n *node) max() (*node, error) {
	if n == nil {
		return nil, fmt.Errorf("cannot find maximum in nil node")
	}

	if n.right == nil {
		return n, nil
	}
	return n.right.max()
}

func (t *BinarySearchTree) String() string {
	if t.root == nil {
		return "empty tree"
	}
	return t.root.string()
}

func (n *node) string() string {
	if n == nil {
		return "nil"
	}
	return fmt.Sprintf("{key: %s, left: %v, right: %v}", n.key, n.left.string(), n.right.string())
}

func (t *BinarySearchTree) InorderTraversal() []string {
	var result []string
	t.root.inorder(&result)
	return result
}

func (n *node) inorder(result *[]string) {
	if n == nil {
		return
	}
	n.left.inorder(result)
	*result = append(*result, n.key)
	n.right.inorder(result)
}

func (t *BinarySearchTree) ReverseInorderTraversal() []string {
	var result []string
	t.root.reverseInorder(&result)
	return result
}

func (n *node) reverseInorder(result *[]string) {
	if n == nil {
		return
	}
	n.right.reverseInorder(result)
	*result = append(*result, n.key)
	n.left.reverseInorder(result)
}

func (t *BinarySearchTree) RangeTraversal(minValue, maxValue string, desc bool) []string {
	var result []string
	t.root.rangeTraversal(&result, minValue, maxValue, desc)
	return result
}

func (n *node) rangeTraversal(result *[]string, minValue, maxValue string, desc bool) {
	if n == nil {
		return
	}
	compareMin := strings.Compare(n.key, minValue)
	compareMax := strings.Compare(n.key, maxValue)

	if desc {
		if compareMax <= 0 {
			n.right.rangeTraversal(result, minValue, maxValue, desc)
		}
		if compareMin >= 0 && compareMax <= 0 {
			*result = append(*result, n.key)
		}
		if compareMin >= 0 {
			n.left.rangeTraversal(result, minValue, maxValue, desc)
		}
	} else {
		if compareMin >= 0 {
			n.left.rangeTraversal(result, minValue, maxValue, desc)
		}
		if compareMin >= 0 && compareMax <= 0 {
			*result = append(*result, n.key)
		}
		if compareMax <= 0 {
			n.right.rangeTraversal(result, minValue, maxValue, desc)
		}
	}
}
