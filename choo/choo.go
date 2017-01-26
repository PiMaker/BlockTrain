package main

import (
    "github.com/PiMaker/BlockTrain/blocktrain"
    "github.com/PiMaker/BlockTrain/ramstore"
	"fmt"
    "github.com/chzyer/readline"
	"strconv"
	"strings"
    "io"
)

func main() {
    fmt.Println("\033[31mGe\033[32mne\033[33msi\033[34ms...\033[0m")
    chain := blocktrain.Genesis(ramstore.NewRAMStore())
    blocktrain.Log = true
    
    fmt.Println("Blocktrain ready, please enter a command:")

    completer := readline.NewPrefixCompleter(
        readline.PcItem("print",
            readline.PcItem("chain"),
            readline.PcItem("latest"),
            readline.PcItem("buffer"),
        ),
        readline.PcItem("commit"),
        readline.PcItem("retrieve"),
        readline.PcItem("verify"),
        readline.PcItem("exit"),
        readline.PcItem("help"),
    )


	l, _ := readline.NewEx(&readline.Config{
		Prompt:          "\033[34mBT>\033[0m ",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
    })

    for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "print "):
			switch line[6:] {
			case "chain":
                chain.PrintChain()
			case "latest":
                if chain.LatestBlock != nil {
                    chain.LatestBlock.PrintBlock()
                } else {
                    fmt.Println("No blocks yet")
                }
            case "buffer":
                fmt.Println("In Buffer: " + strconv.Itoa(len(chain.Buffer)) + "/" + strconv.Itoa(blocktrain.BufferSize))
                for i, b := range chain.Buffer {
                    fmt.Println(strconv.Itoa(i) + ": " + b.TxID)
                }
			default:
				fmt.Println("Cant't print: ", line[6:])
			}

        case line == "print":
            fmt.Println("Usage: print (chain|latest|buffer)")

        case strings.HasPrefix(line, "help"):
            fmt.Println(completer.Tree("  "))

        case strings.HasPrefix(line, "exit"):
            goto exit

        case strings.HasPrefix(line, "commit "):
            data := []byte(line[7:])
            chain.Commit(data)

        case strings.HasPrefix(line, "retrieve "):
            data := chain.Retrieve(line[9:])
            fmt.Println(string(data))
            status := chain.Verify(line[9:], data)
            fmt.Println("Status: " + blocktrain.StatusToString(status))

        case line == "verify":
            fmt.Println("Usage: verify <txID> <d a t a>")
            continue

        case strings.HasPrefix(line, "verify "):
            spaceIndex := strings.Index(line[7:], " ")
            if spaceIndex == -1 {
                fmt.Println("Usage: verify <txID> <d a t a>")
                continue
            }

            txID := strings.TrimSpace(line[7:7+spaceIndex])
            data := line[7+spaceIndex+1:]

            status := chain.Verify(txID, []byte(data))
            fmt.Println("Status: " + blocktrain.StatusToString(status))
        }
    }

exit:
    fmt.Println("The BlockTrain is leaving the station...")
}

