# MyNode

The humble basics of a bitcoin node that can handshake with another node, written as a coding exercise.

To run the test suite, simply

```
go test ./...
```

To check that the handshake with another node succeeds, first run a local bitcoin node. For example:

```shell
go run github.com/btcsuite/btcd@latest
```


> [!NOTE]  
> If you delete `~/.btcd/data/mainnet/peers.json` and pass `-ddebug --nodnsseed` as arguments to btcd, you will be able to see `mynode` connecting in the log.

`

Then, in another terminal, run the example program
```shell
go run ./cmd
```

Example output:

```
➜  fluffy:mynode git:(main) ✗ go run ./cmd
2024/06/20 16:00:54 Connected to 127.0.0.1:8333
2024/06/20 16:00:54 sending our version {Version:70016 Services:1 Timestamp:1718895654 AddrRecv:{Time:1718895654 Services:1 IP:127.0.0.1:8333} AddrFrom:{Time:0 Services:0 IP:invalid AddrPort} Nonce:7595319339721391208 UserAgent:/pscott31-mynode:0.0.1/ StartHeight:0 Relay:false}
2024/06/20 16:00:54 received their version: {Version:70016 Services:1101 Timestamp:1718895654 AddrRecv:{Time:0 Services:1101 IP:127.0.0.1:44784} AddrFrom:{Time:0 Services:1101 IP:invalid AddrPort} Nonce:7241563614374867889 UserAgent:/btcwire:0.5.0/btcd:0.24.2/ StartHeight:385463 Relay:true}
2024/06/20 16:00:54 sending version acknowledgement
2024/06/20 16:00:54 received message: {Magic:3652501241 Command:sendaddrv2 Length:0 Checksum:3806393949 Payload:[]}
2024/06/20 16:00:54 received message: {Magic:3652501241 Command:verack Length:0 Checksum:3806393949 Payload:[]}
```

