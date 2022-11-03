package main

import (
	"log"
	
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/daoio/evm-opcode-finder/finder"
)

const URL = "YOUR_URL"
const opcode = "SELFDESTRUCT"

func main() {
	client, err := ethclient.Dial(URL)

	if err != nil {
		log.Fatal(err)
	}

	finder.FindOpcode(client, vm.StringToOp(opcode))
}
