package finder

import (
	"log"
	"context"
	"fmt"
	
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/types"
)

// find opcode in the contracts from latest block
func FindOpcode(client *ethclient.Client, opcode vm.OpCode) {
	block := latestBlock(client)
	inspectContractsInBlock(client, block, opcode)
}

func FindOpcodeInContract(client *ethclient.Client, address common.Address, opcode vm.OpCode) {
	_, bytecode := isContract(client, address)
	
	if compareOpcodes(bytecode, opcode) {
		success(opcode, address)
	}
}

//==========================Blocks==========================

// Find all contracts interactions in the latest block.
// looking here for `To` field in transactions to find
// transactions sent to contracts from EOA
func inspectContractsInBlock(client *ethclient.Client, block *types.Block, opcode vm.OpCode) {
	txs := block.Transactions()
	printStart(opcode, block)
	
	for _, tx := range txs {
		address := *tx.To()
		contract, bytecode := isContract(client, address)

		if contract {
			inspectBytecode(client, opcode, bytecode, address)
		}
	}
}

func latestBlock(client *ethclient.Client) *types.Block {
	// get latest block
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return block
}

//==========================Bytecode==========================

func inspectBytecode(client *ethclient.Client, opcode vm.OpCode, bytecode []byte, address common.Address) {
	if compareOpcodes(bytecode, opcode) {
		success(opcode, address)
	}
}

// actually not all bytes in bytecode is opcodes it can be arguments to opcodes also
// that's why all PUSH instructions are skipped
func compareOpcodes(bytecode []byte, opcode vm.OpCode) bool {
	for i := 0; i < len(bytecode); {
		op := vm.OpCode(bytecode[i])

		// skip PUSH opcodes
		if op.IsPush() {
			// increment i, skipping PUSH argument
			i += int(skipPush(op, i))
		} else if op == opcode {
			return true
		}
		i++
	}

	return false
}

//==========================Helpers==========================

func skipPush(pushXX vm.OpCode, counter int) uint8 {
	// PUSH1 in dec -> 96, PUSH32 in dec -> 127
	// to get amount of bytes to push, i.e. for us
	// amount of bytes to skip, we subtracting 95 from
	// OpCode, since byte is alias for uint8, we will
	// get the number of bytes to push => number to add
	// to counter
	var push uint8 = byte(pushXX)
	
	return (push - 95)
}

// check if address has runtime bytecode
func isContract(client *ethclient.Client, address common.Address) (bool, []byte) {
	bytecode, err := client.CodeAt(context.Background(), address, nil)

	if err != nil {
		log.Fatal(err)
	}

	if len(bytecode) > 0 {
		return true, bytecode
	} else {
		return false, bytecode
	}
}

// Separate logs in this two functions
func printStart(opcode vm.OpCode, block *types.Block) {
	fmt.Printf("========================\nLooking for %v in contracts\n in block %v\n========================\n", opcode.String(), block.Number())
}

func success(opcode vm.OpCode, address common.Address) {
	fmt.Printf("Catch opcode: %v\n in contract: %v\n", opcode, address)
}
