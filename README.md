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

2019/04/20 14:30:04 reporting.go:23: Successful operation: "oopG1PWdhM8YUkADTwWBHdgmtMxproSPurN8GBAHXLJTnCMb27T"

+--------------------------------------+-----------+
|               ADDRESS                |  PAYMENT  |
+--------------------------------------+-----------+
| KT1WvyJ1qUrWzShA2T6QeL7AW4DR6GspUimM |  0.024563 |
| KT1Veiiv3NvVF4YpZmSjUaJP7hhKqda4LX4A |  0.506465 |
| KT1VBXqcStE2sonpWSFTRurzVgdtAPT5HdEU |  2.532329 |
| KT1UPrvoYfAabgBT4YRFwmdxrHhURu45Yhtj |  0.075210 |
| KT1UPe6SuBx1tKg6fD2DgEzqhaJG1UpYTLVq |  0.001878 |
| KT1TJThjwjBuCJDDXp9CQkBgeNSRYbahQW8Y |  1.266164 |
| KT1SbZzZWGj7HTpZwSiU8BT9X9jVMvy4LX5t |  2.581583 |
| KT1NeVTJYNkrtGqKy6Z84rH1zF86VBnQYFTT |  0.101039 |
| KT1N82y6jxThaot6bQC71SjGLDAgrpbchhJm |  0.253708 |
| KT1N2iCTCpFJPZtncR7iTi5kTfnq7ahwvkaY |  5.163644 |
| KT1KhU8TK59KHwvoXnEPP8232tfBVQkDJN3M |  0.404740 |
| KT1JuKotYDPZd71cUtaR2AwCE2vfggMaGgKt |  0.253222 |
| KT1JexcFezMnUAaWmvUGY99jwTA4jcKiUgFp |  2.405681 |
| KT1HE1JLUZwsVYaCTBaGUnk652TfvWWeBsCu |  2.279089 |
| KT1FGkavmRTjFtJ3pzT96XHj1zxCFc91KwKQ |  0.005051 |
| KT1F5Gi1rLAuengNV9MdTsZbVSzJwEJPiQ1K |  0.430489 |
| KT1EmjARhUY7n44HTXCLJd4cCPfXKm7k6fZA | 52.761327 |
| KT1DKVQjcy2b9yDHQrsxNWZaKY7KfTWYvqzG |  0.164084 |
| KT1CBzyabE55icq7vCaGEPWWQwHiJ1vnDpeL |  0.253232 |
| KT1BmDCykjJr4zTHLNuojLkqKmcsGm7vxEMt |  1.519397 |
| KT1AvvU63C6GzrELbgEmsiyzTxQhNTQy3AZR |  2.532329 |
| KT1ACieGEX5WM1D8skTPGScw9jxv165t43Vk |  0.001871 |
| KT19LRjvtkaHNBHjyikNn7jWbkCVBbcPaB2K |  0.253232 |
| KT18j3UkJgNxqbjmy9J9rtFCmzKXynPWDqUy |  0.005058 |
+--------------------------------------+-----------+
|                TOTAL                 | 75.775384 |
+--------------------------------------+-----------+
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