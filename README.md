# Payman

Payman is a golang driven payout tool for delegation services that is built upon the [go-tezos](https://github.com/DefinitelyNotAGoat/go-tezos) library. This project is in alpha and could change frequently. 

## Installation

Get Payman 
```
go get github.com/DefinitelyNotAGoat/payman
cd $GOPATH/github.com/DefinitelyNotAGoat/payman
go build
```

## Payman Documentation

### help:

```
Usage:
  payman payout [flags]

Flags:
  -c, --cycle int             cycle to payout for (e.g. 95)
      --cycles string         cycles to payout for (e.g. 95-100)
  -d, --delegate string       public key hash of the delegate that's paying out (e.g. --delegate=<phk>)
      --dry                   run payout in simulation with report (default false)(e.g. --dry)
  -f, --fee float32           fee for the delegate (e.g. 0.05 = 5%) (default -1)
      --gas-limit int         network gas limit for each transaction in mutez (default 10200)(e.g. 10300) (default 10200)
  -h, --help                  help for payout
  -l, --log-file string       file to log to (default stdout)(e.g. ./payman.log) (default "/dev/stdout")
      --network-fee int       network fee for each transaction in mutez (default 1270)(e.g. 2000) (default 1270)
  -n, --node string           address to the node to query (default http://127.0.0.1)(e.g. mainnet-node.tzscan.io) (default "http://127.0.0.1")
  -k, --password string       password to the secret key of the wallet paying (e.g. --password=<passwd>)
  -p, --port string           port to use for node (default 8732)(e.g. 443) (default "8732")
  -r, --reddit string         path to reddit agent file (initiates reddit bot)(e.g. https://turnage.gitbooks.io/graw/content/chapter1.html)
      --reddit-title string   pre title for the reddit bot to post (e.g. DefinitelyNotABot: -- will read DefinitelyNotABot: Payout for Cycle(s) <cycles>)
  -s, --secret string         encrypted secret key of the wallet paying (e.g. --secret=<sk>)
      --serve                 run service to payout for all new cycles going foward (default false)(e.g. --serve)
```

### example: 
The example below will payout delegations for delegate `tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB` with wallet sk `edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm` and password `abcd1234` for cycle `152` with a tezos node at `127.0.0.1:8732`:
```
./payman payout \     
        --delegate=tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB \
        --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm \
        --password=abcd1234 \
        --cycle=160 \
        --node=127.0.0.1 \
        --port=8732 \
        --fee=0.05

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
./payman payout \     
        --delegate=tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB \
        --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm \
        --password=abcd1234 \
        --cycle=160 \
        --node=127.0.0.1 \
        --port=8732 \
        --fee=0.05 \
        --serve
``` 

#### Using Reddit Bot Functionality 
This feature is currently only functional with mainnet. If used with another network, the link in your reddit post will be broken (Future Fix)
```
-r, --reddit string     example https://turnage.gitbooks.io/graw/content/chapter1.html
--title string      example "MyService:"
```
To post the results of your payout to reddit, create a [reddit.agent](https://turnage.gitbooks.io/graw/content/chapter1.html) file  containing your reddit client and secret. 
```
user_agent: "<platform>:<app ID>:<version string> (by /u/<reddit username>)"
client_id: "client id (looks kind of like: sdkfbwi48rhijwsdn)"
client_secret: "client secret (looks kind of like: ldkvblwiu34y8hsldjivn)"
username: "reddit username"
password: "reddit password"
```

Pass that file to the `--reddit` / `-r` flag. You may also include an option pre title using the `--reddit-title` flag. 
```
./payman payout \     
        --delegate=tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB \
        --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm \
        --password=abcd1234 \
        --cycle=160 \
        --node=127.0.0.1 \
        --port=8732 \
        --fee=0.05 \
        --reddit=./reddit.agent \
        --reddit-title=DefinitelyNotABot: 

```

Your reddit post will contain the title passed along with "Payout for Cycle(s) number"
```
DefinitelyNotATestBot: Payout for Cycle(s) 100
```
With a link to the tzscan operation related to the cycle.


## Roadmap:
* blacklist addresses
* defer option that allows you to defer a percentage of your fee to another address
* inlucde fiat price of XTZ in reports
* tax reporting

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details