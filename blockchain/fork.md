# update

## ore/block_validator.go

## consensus/clique/clique.go

```go
timing := time.Duration(100) * time.Millisecond
    if c.config.Period == 0 {
        if delay < time.Duration(0)*time.Millisecond {
            delay = timing
        } else {
            delay += timing
        }
    }
    if header.Difficulty.Cmp(diffNoTurn) == 0 {
        // It's not our turn explicitly to sign, delay it a bit
        wiggle := time.Duration(len(snap.Signers)/2+1)*wiggleTime - time.Duration(200)*time.Millisecond
        delay += time.Duration(rand.Int63n(int64(wiggle))) + time.Duration(200)*time.Millisecond

        log.Trace("Out-of-turn signing requested", "wiggle", common.PrettyDuration(wiggle))
    }
    log.Trace("Waiting for slot to sign and propagate", "delay", common.PrettyDuration(delay))
```

## params/protocol_params.go

GenesisGasLimit      uint64 = 471238800    (4712388)    --targetgaslimit 80000000

TxGas	uint64 = 11000	(21000)

TxGasContractCreation	uint64 = 23000	(53000)

TxDataNonZeroGas	uint64 = 8	(68)

LogGas	uint64 = 75	(375)

LogTopicGas	uint64 = 75	(375)

MaxCodeSize = 245760                       (24576)

## core/tx_pool.go

if tx.Size() > 320*1024 {                   (32*1024)   reject transactions over 32KB to prevent DOS attacks (extend to 320KB)
    return ErrOversizedData
}

## core/genesis.go

func DeveloperGenesisBlock() {
    GasLimit:   params.GenesisGasLimit		(6283185)
    ...
