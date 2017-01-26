package blocktrain

import (
	"time"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
    "strconv"
)

type Block struct {
    Timestamp time.Time
    PrevHash []byte
    TXHash []byte
}

type Chain struct {
    Store DataStore
    SeedHash []byte
    Buffer []MerkleData
    LatestBlock *Block
}

type DataStore interface {
    InsertTree(root *MerkleNode)
    InsertBlock(block *Block)
    InsertData(data []byte, txID string)
    GetTree(hash []byte) *MerkleNode
    GetBlock(hash []byte) *Block
    GetData(txID string) []byte
}

var Log = true

type VerificationStatus int
const (
    InBuffer VerificationStatus = iota
    InBufferInvalid
    Invalid
    UnknownTxID
    Verified
)

const (
    BufferSize = 4
)

func StatusToString(status VerificationStatus) string {
    switch status {
    case InBuffer:
        return "In buffer (not yet in a block)"
    case InBufferInvalid:
        return "In buffer (not yet in a block) / Invalid (data-txID mismatch)"
    case Invalid:
        return "Invalid"
    case UnknownTxID:
        return "Unknown transaction ID"
    case Verified:
        return "Verified"
    default:
        return "Unknown status"
    }
}

func Genesis(store DataStore) *Chain {
    genesisHash := make([]byte, 32)
    n, err := rand.Reader.Read(genesisHash)
    if err != nil || n != 32 {
        panic("Error on creating the genesis hash")
    }

    chain := &Chain{
        Store: store,
        SeedHash: genesisHash,
        Buffer: make([]MerkleData, 0),
    }

    return chain
}

func (block *Block) Hash() []byte {
    hash := block.TXHash
    hash = append(hash, []byte(strconv.FormatInt(block.Timestamp.UnixNano(), 10))...)
    hash = append(hash, block.PrevHash...)
    return Hash(hash)
}

func NewBlock(tx *MerkleNode, prevBlock *Block) *Block {
    hash := prevBlock.Hash()

    newBlock := &Block{
        PrevHash: hash,
        Timestamp: time.Now(),
        TXHash: tx.Hash,
    }

    return newBlock
}

func (chain *Chain) Commit(data []byte) (txID string) {
    txIDb := sha256.Sum256(data)
    txIDb = sha256.Sum256(append(txIDb[:], []byte(strconv.FormatInt(time.Now().UnixNano(), 10))...))
    txID = base32.StdEncoding.EncodeToString(txIDb[:])

    chain.Buffer = append(chain.Buffer, MerkleData{
        Data: data,
        TxID: txID,
    })
    chain.Store.InsertData(data, txID)

    if Log {
        fmt.Println("TXid committed: " + txID)
    }

    if len(chain.Buffer) == BufferSize {
        tx := NewMerkleTree(chain.Buffer)
        chain.Store.InsertTree(tx)
        
        if chain.LatestBlock == nil {
            chain.LatestBlock = &Block{
                PrevHash: chain.SeedHash,
                Timestamp: time.Now(),
                TXHash: tx.Hash,
            }
        } else {
            chain.LatestBlock = NewBlock(tx, chain.LatestBlock)
        }

        chain.Store.InsertBlock(chain.LatestBlock)

        chain.Buffer = make([]MerkleData, 0)

        if Log {
            fmt.Println("Block was appended: " + base32.StdEncoding.EncodeToString(chain.LatestBlock.Hash()))
        }
    }

    return txID
}

func (chain *Chain) Retrieve(txID string) []byte {
    return chain.Store.GetData(txID)
}

func (chain *Chain) Verify(txID string, data []byte) VerificationStatus {
    for _, bufferTX := range chain.Buffer {
        if txID == bufferTX.TxID {
            if byteSliceEqual(bufferTX.Data, data) {
                return InBuffer
            }

            return InBufferInvalid
        }
    }

    block := chain.LatestBlock
    for block != nil {
        tree := chain.Store.GetTree(block.TXHash)
        if tree.Contains(txID) {
            verified := tree.Verify(txID, data)
            if verified {
                return Verified
            }
            
            return Invalid
        }

        block = chain.Store.GetBlock(block.PrevHash)
    }

    return UnknownTxID
}

func (chain *Chain) PrintChain() {
    fmt.Println("~~ BlockTrain Chainlog ~~")
    fmt.Println("Seed: " + base32.StdEncoding.EncodeToString(chain.SeedHash))
    fmt.Println("TX-Buffer: " + strconv.Itoa(len(chain.Buffer)) + "/" + strconv.Itoa(BufferSize))
    fmt.Println()
    fmt.Println("Attached blocks:")
    fmt.Println()

    count := 0

    block := chain.LatestBlock
    for block != nil {
        block.PrintBlock()
        fmt.Println()
        block = chain.Store.GetBlock(block.PrevHash)
        count++
    }

    fmt.Println(strconv.Itoa(count) + " Blocks in chain\n")

    fmt.Println("~~ End of Chainlog ~~")
}

func (block *Block) PrintBlock() {
    fmt.Println("### " + base32.StdEncoding.EncodeToString(block.Hash()))
    fmt.Println("Timestamp: \t" + strconv.FormatInt(block.Timestamp.UnixNano(), 10))
    fmt.Println("TX-Hash: \t" + base32.StdEncoding.EncodeToString(block.TXHash))
    fmt.Println("Linked Block: \t" + base32.StdEncoding.EncodeToString(block.PrevHash))
}

func byteSliceEqual(s1, s2 []byte) bool {
    if len(s1) != len(s2) {
        return false
    }

    for i := 0; i < len(s1); i++ {
        if s1[i] != s2[i] {
            return false
        }
    }

    return true
}

func Hash(input []byte) []byte {
    hashed := sha256.Sum256(input)
    return hashed[:]
}