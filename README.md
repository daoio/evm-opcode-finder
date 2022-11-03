# Abilities

1. Search for opcode in contracts in the latest block
```go
  func FindOpcode(client *ethclient.Client, opcode vm.OpCode) {
	  block := latestBlock(client)
	  inspectContractsInBlock(client, block, opcode)
  }
```
2. Search for opcode in certain contract
```go
  func FindOpcodeInContract(client *ethclient.Client, address common.Address, opcode vm.OpCode) {
    _, bytecode := isContract(client, address)

    if compareOpcodes(bytecode, opcode) {
      success(opcode, address)
    }
  }
```
Note: it doesn't search for PUSH opcodes, since they're exist in bytecode of all contracts and it skips all PUSH instructions with their arguments

# Example

```go
  const opcode = "SELFDESTRUCT"

  func main() {
    client, err := ethclient.Dial(URL)

    if err != nil {
      log.Fatal(err)
    }

    finder.FindOpcode(client, vm.StringToOp(opcode))
  }
```
Search for SELFDESTRUCT opcode in contracts in the latest block

OUTPUT
```shell
  ========================
  Looking for SELFDESTRUCT in contracts
   in block 15888893
  ========================
  Catch opcode: SELFDESTRUCT
   in contract: 0x46A82Ec528d89154EF3Dc66d9E03fEd617886d2c
  Catch opcode: SELFDESTRUCT
   in contract: 0x1111111254fb6c44bAC0beD2854e76F90643097d
  Catch opcode: SELFDESTRUCT
   in contract: 0xe069aE4B336Ca73142cDc5206ed4a4d3A3ff39f6
```
Inspecting contracts in etherscan:
1. https://etherscan.io/address/0x46A82Ec528d89154EF3Dc66d9E03fEd617886d2c#code 👍 it has SELFDESTRUCT
2. https://etherscan.io/address/0x1111111254fb6c44bAC0beD2854e76F90643097d#code 👍 it has SELFDESTRUCT
3. https://etherscan.io/address/0xe069aE4B336Ca73142cDc5206ed4a4d3A3ff39f6#code 👍 it has SELFDESTRUCT (close to the end it has `ff` opcode  which is SELFDESTRUCT)

# Limitations
1. At the moment it can return the same addresses in the block, i.e. it doesn't check for collisions
2. It inspects only one latest block, stopping then execution
