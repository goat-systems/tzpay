# Payman

Payman is a golang driven payout tool for delegation services that is built upon the [go-tezos](https://github.com/DefinitelyNotAGoat/go-tezos) library. This project is in alpha and could change frequently. 

## Installation

Install pkg_config (debian example below):
```
sudo apt-get install pkg-config
```

Install [libsoidum](https://libsodium.gitbook.io/doc/installation)

Get Payman 
```
go get github.com/DefinitelyNotAGoat/payman
```

## Payman Documentation

### help:

```
Usage:
  payout [flags]

Flags:
  -c, --cycle int         cycle to payout for, example 20
      --cycles string     cycles to payout for, example 20-24
  -d, --delegate string   public key hash of the delegate that's paying out
      --dry               run payout in simulation with report
  -f, --fee float32       example 0.05 (default 0.05)
      --gas-limit int     network gas limit for each transaction in mutez (default 10200)
  -h, --help              help for payout
  -l, --log-file string   example ./payman.log (default "/dev/stdout")
      --network-fee int   network fee for each transaction in mutez (default 1270)
  -n, --node string       example mainnet-node.tzscan.io (default "http://127.0.0.1")
  -k, --password string   password to the secret key of the wallet paying
  -p, --port string       example 8732 (default "8732")
  -s, --secret string     encrypted secret key of the wallet paying
      --serve             run service to payout for all new cycles going foward
```

### example: 
The example below will payout delegations for delegate `tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB` with wallet sk `edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm` and password `abcd1234` for cycle `152` with a tezos node at `127.0.0.1:8732`:
```
./payman -d=tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB -s=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm -k=abcd1234 --cycle=152 -n=127.
0.0.1 -p=8732
```

The below will do the same as above but also start a server that will payout for every cycle going forward:
```
./payman -d=tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB -s=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm -k=abcd1234 --cycle=152 -n=127.
0.0.1 -p=8732 --serve
``` 

## Roadmap:
* blacklist addresses
* defer option that allows you to defer a percentage of your fee to another address
* inlucde fiat price of XTZ in reports
* tax reporting

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details