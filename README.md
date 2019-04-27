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
Payout pays out rewards to delegations for the delegate passed.

Usage:
  payman payout [flags]

Flags:
  -c, --cycle int             cycle to payout for (e.g. 95)
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
The example below will payout delegations for delegate `tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB` with wallet sk `edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm` and password `abcd1234` for cycle `160` with a tezos node at `127.0.0.1:8732`:
```
./payman payout \
    --delegate=tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB \
    --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm \
    --password=abcd1234 \
    --cycle=160 \
    --node=127.0.0.1 \
    --port=8732 \
    --fee=0.05 

2019/04/27 11:47:22 reporting.go:24: Successful operation: "oomck2kQt4TXbyXseTLUmuTfVTzf8x3QxU4nFPSFyPfp9U6Vwsx"

+--------------------------------------+----------+-----------+----------+-----------+
|               ADDRESS                |  SHARE   |   GROSS   |   FEE    |    NET    |
+--------------------------------------+----------+-----------+----------+-----------+
| KT18j3UkJgNxqbjmy9J9rtFCmzKXynPWDqUy | 0.000015 |  0.005982 | 0.000299 |  0.005683 |
| KT1Veiiv3NvVF4YpZmSjUaJP7hhKqda4LX4A | 0.001540 |  0.599009 | 0.029950 |  0.569059 |
| KT1WvyJ1qUrWzShA2T6QeL7AW4DR6GspUimM | 0.000075 |  0.029051 | 0.001452 |  0.027599 |
| KT1VBXqcStE2sonpWSFTRurzVgdtAPT5HdEU | 0.007701 |  2.995048 | 0.149752 |  2.845296 |
| KT1UPrvoYfAabgBT4YRFwmdxrHhURu45Yhtj | 0.000229 |  0.088952 | 0.004447 |  0.084505 |
| KT1UPe6SuBx1tKg6fD2DgEzqhaJG1UpYTLVq | 0.000006 |  0.002221 | 0.000111 |  0.002110 |
| KT1TJThjwjBuCJDDXp9CQkBgeNSRYbahQW8Y | 0.003850 |  1.497524 | 0.074876 |  1.422648 |
| KT1T3ZwVCAsNPaFqmUj53qsWbkqEj4C8E14w | 0.000559 |  0.217445 | 0.010872 |  0.206573 |
| KT1NeVTJYNkrtGqKy6Z84rH1zF86VBnQYFTT | 0.000307 |  0.119502 | 0.005975 |  0.113527 |
| KT1RNyjoqV8AbeSGiqxKkRC4Wx4CsjK5vhUh | 0.000000 |  0.000000 | 0.000000 |  0.000000 |
| KT1N2iCTCpFJPZtncR7iTi5kTfnq7ahwvkaY | 0.015702 |  6.107170 | 0.305358 |  5.801812 |
| KT1N82y6jxThaot6bQC71SjGLDAgrpbchhJm | 0.000772 |  0.300067 | 0.015003 |  0.285064 |
| KT1L8HheomV6WzQrxTYqf9h3zFTcfdE7ELEg | 0.000698 |  0.271573 | 0.013578 |  0.257995 |
| KT1KhU8TK59KHwvoXnEPP8232tfBVQkDJN3M | 0.001231 |  0.478696 | 0.023934 |  0.454762 |
| KT1JuKotYDPZd71cUtaR2AwCE2vfggMaGgKt | 0.000770 |  0.299492 | 0.014974 |  0.284518 |
| KT1JexcFezMnUAaWmvUGY99jwTA4jcKiUgFp | 0.007316 |  2.845259 | 0.142262 |  2.702997 |
| KT1J3m8h86UiXKYwqXj6xBUMEjucbGeTJF52 | 0.000000 |  0.000000 | 0.000000 |  0.000000 |
| KT1FGkavmRTjFtJ3pzT96XHj1zxCFc91KwKQ | 0.000015 |  0.005974 | 0.000298 |  0.005676 |
| KT1HE1JLUZwsVYaCTBaGUnk652TfvWWeBsCu | 0.006931 |  2.695535 | 0.134776 |  2.560759 |
| KT1F5Gi1rLAuengNV9MdTsZbVSzJwEJPiQ1K | 0.001309 |  0.509150 | 0.025457 |  0.483693 |
| KT1EmjARhUY7n44HTXCLJd4cCPfXKm7k6fZA | 0.160444 | 62.402129 | 3.120106 | 59.282023 |
| KT1DKVQjcy2b9yDHQrsxNWZaKY7KfTWYvqzG | 0.000499 |  0.194066 | 0.009703 |  0.184363 |
| KT1CBzyabE55icq7vCaGEPWWQwHiJ1vnDpeL | 0.000770 |  0.299504 | 0.014975 |  0.284529 |
| KT1BmDCykjJr4zTHLNuojLkqKmcsGm7vxEMt | 0.004620 |  1.797028 | 0.089851 |  1.707177 |
| KT1AvvU63C6GzrELbgEmsiyzTxQhNTQy3AZR | 0.007701 |  2.995048 | 0.149752 |  2.845296 |
| KT1ACieGEX5WM1D8skTPGScw9jxv165t43Vk | 0.000006 |  0.002213 | 0.000110 |  0.002103 |
| KT19LRjvtkaHNBHjyikNn7jWbkCVBbcPaB2K | 0.000770 |  0.299504 | 0.014975 |  0.284529 |
+--------------------------------------+----------+-----------+----------+-----------+
|                                         TOTAL   | 87.057142 | 4.352846 | 82.704296 |
+--------------------------------------+----------+-----------+----------+-----------+
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