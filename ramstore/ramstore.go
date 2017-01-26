package ramstore

import (
	"github.com/PiMaker/BlockTrain/blocktrain"
)

type RAMStore struct {
    blockstore map[string]*blocktrain.Block
    txstore map[string]*blocktrain.MerkleNode
    datastore map[string][]byte
}

func NewRAMStore() RAMStore {
    store := RAMStore{}
    store.blockstore = make(map[string]*blocktrain.Block)
    store.txstore = make(map[string]*blocktrain.MerkleNode)
    store.datastore = make(map[string][]byte)
    return store
}

func (store RAMStore) InsertTree(root *blocktrain.MerkleNode) {
    store.txstore[string(root.Hash)] = root
}

func (store RAMStore) InsertBlock(block *blocktrain.Block) {
    store.blockstore[string(block.Hash())] = block
}

func (store RAMStore) InsertData(data []byte, txID string) {
    store.datastore[txID] = data
}

func (store RAMStore) GetTree(hash []byte) *blocktrain.MerkleNode {
    tree, ok := store.txstore[string(hash)]
    if ok {
        return tree
    }
    
    return nil
}

func (store RAMStore) GetBlock(hash []byte) *blocktrain.Block {
    block, ok := store.blockstore[string(hash)]
    if ok {
        return block
    }
    
    return nil
}

func (store RAMStore) GetData(txID string) []byte {
    data, ok := store.datastore[txID]
    if ok {
        return data
    }
    
    return nil
}