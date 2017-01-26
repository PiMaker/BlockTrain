package blocktrain

import (
	"fmt"
	"encoding/base32"
)

type MerkleNode struct {
    Hash []byte
    TxID string
    left *MerkleNode
    right *MerkleNode

    visited bool `json:"-"`
}

type MerkleData struct {
    Data []byte
    TxID string
}

func NewMerkleTree(data []MerkleData) *MerkleNode {
    if len(data) == 0 {
        return nil
    }

    var list []*MerkleNode

    for _, newData := range data {
        chkSum := Hash(append(newData.Data, []byte(newData.TxID)...))
        newNode := &MerkleNode{Hash: chkSum[:], TxID: newData.TxID}
        list = append(list, newNode)
    }

    for len(list) > 1 {

        count := len(list)/2

        tmpList := make([]*MerkleNode, count)

        for i := 0; i < count; i++ {
            node := &MerkleNode{
                left: list[i*2],
                right: list[i*2+1],
            }

            var toHash []byte
            toHash = append(toHash, node.left.Hash...)
            toHash = append(toHash, node.right.Hash...)

            chkSum := Hash(toHash)
            node.Hash = chkSum[:]

            tmpList[i] = node
        }

        if len(list)%2==1 {
            tmpList = append(tmpList, list[len(list)-1])
        }

        list = tmpList
    }

    root := list[0]

    var toHash []byte
    toHash = append(toHash, root.left.Hash...)
    toHash = append(toHash, root.right.Hash...)

    chkSum := Hash(toHash)
    root.Hash = chkSum[:]

    return root
}

func (node *MerkleNode) Contains(txID string) bool {
    if node.TxID == txID {
        return true
    }

    left, right := false, false

    if node.left != nil {
        left = node.left.Contains(txID)
    }

    if node.right != nil {
        right = node.right.Contains(txID)
    }

    return left || right
}

func (node *MerkleNode) Verify(txID string, data []byte) bool {
    node.unvisit()

    stack := make(stack, 0)
    stack.Push(node)
    
    for len(stack) > 0 {
        n := stack.Top()
        n.visited = true

        if n.left == nil && n.right == nil {
            if n.TxID == txID {
                break
            }

            stack.Pop()
            continue
        }

        if n.left != nil && !n.left.visited {
            stack.Push(n.left)
            continue
        }
        if n.right != nil && !n.right.visited {
            stack.Push(n.right)
            continue
        }

        stack.Pop()
    }

    if len(stack) == 0 {
        return false
    }

    dataHash := Hash(append(data, []byte(txID)...))

    if !byteSliceEqual(stack.Top().Hash, dataHash) {
        return false
    }

    stack.Pop()

    for len(stack) > 0 {
        n := stack.Pop()

        var toHash []byte

        if n.left != nil {
            if n.left.TxID == txID {
                toHash = append(toHash, dataHash...)
            } else {
                toHash = append(toHash, n.left.Hash...)
            }
        }

        if n.right != nil {
            if n.right.TxID == txID {
                toHash = append(toHash, dataHash...)
            } else {
                toHash = append(toHash, n.right.Hash...)
            }
        }

        chkSum := Hash(toHash)

        if !byteSliceEqual(chkSum, n.Hash) {
            return false
        }
    }

    return true
}

func (node *MerkleNode) unvisit() {
    node.visited = false
    if node.left != nil {
        node.left.unvisit()
    }
    if node.right != nil {
        node.right.unvisit()
    }
}

func (node *MerkleNode) PrintRecursive(level int) {
    for i := 0; i < level; i++ {
        fmt.Print("  ")
    }

    if (len(node.Hash) > 0) {
        fmt.Println(base32.StdEncoding.EncodeToString(node.Hash))
    } else {
        fmt.Println("<no hash>")
    }

    if node.left != nil {
        node.left.PrintRecursive(level + 1)
    }

    if node.right != nil {
        node.right.PrintRecursive(level + 1)
    }
}