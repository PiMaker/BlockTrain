package blocktrain

type stack []*MerkleNode

func (s *stack) Push(v *MerkleNode) {
    *s = append(*s, v)
}

func (s *stack) Pop() *MerkleNode {
    if len(*s) == 0 {
        return nil
    }

    res:=(*s)[len(*s)-1]
    *s=(*s)[:len(*s)-1]
    return res
}

func (s *stack) Top() *MerkleNode {
    if len(*s) == 0 {
        return nil
    }

    res:=(*s)[len(*s)-1]
    return res
}