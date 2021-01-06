# Payman

Payman is a golang driven payout tool for delegation services that is built upon the [go-tezos](https://github.com/DefinitelyNotAGoat/go-tezos) library.

## Installation

### Source

```
go get github.com/DefinitelyNotAGoat/payman
cd $GOPATH/github.com/DefinitelyNotAGoat/payman
go build
```

### Linux

```
wget https://github.com/DefinitelyNotAGoat/payman/releases/download/v1.0.4/payman_linux_amd64
sudo mv payman_linux_amd64 /usr/local/bin/payman
sudo chmod a+x /usr/local/bin/payman
```

### MacOS

```
wget https://github.com/DefinitelyNotAGoat/payman/releases/download/v1.0.4/payman_darwin_amd64
sudo mv payman_linux_amd64 /usr/local/bin/payman
sudo chmod a+x /usr/local/bin/payman

```

## Payman Usage

=======

## Standard Stake Capital Payout

Simply insert the secret and password for the wallet:

```
./payman payout --delegate=tz1Z3KCf8CLGAYfvVWPEr562jDDyWkwNF7sT --cycle=143 --node=https://mainnet.tezrpc.me --fee=0.1 --secret=edesk22... --password=... --network-fee=1420
```

To first check the total cost of a cycle simply run:

```
./payman report --delegate=tz1Z3KCf8CLGAYfvVWPEr562jDDyWkwNF7sT --cycle=143 --node=https://mainnet.tezrpc.me --fee=0.1
```

## Payman Documentation

Github Pages: https://definitelynotagoat.github.io/payman/

### Report

#### Help

```
payman report --help
report simulates a payout and generates a table and csv report

Usage:
  payman report [flags]

Flags:
  -c, --cycle int         cycle to payout for (e.g. 95)
  -d, --delegate string   public key hash of the delegate that's paying out (e.g. --delegate=<phk>)
  -f, --fee float32       fee for the delegate (e.g. 0.05 = 5%) (default -1)
  -h, --help              help for report
  -l, --log-file string   file to log to (default stdout)(e.g. ./payman.log) (default "/dev/stdout")
  -u, --node string       address to the node to query (default http://127.0.0.1:8732)(e.g. https://mainnet-node.tzscan.io:443) (default "http://127.0.0.1:8732")
      --payout-min int    will only payout to addresses that meet the payout minimum (e.g. --payout-min=<mutez>)
```

#### Example

```
payman report --delegate=tz1SF9wBoBQbFUF13agZ8EgihLCKM54G1ccV --cycle=184 --fee=0.05 --payout-min=5000

+--------------------------------------+-----------+-------------+-----------+-------------+
|               ADDRESS                |   SHARE   |    GROSS    |    FEE    |     NET     |
+--------------------------------------+-----------+-------------+-----------+-------------+
| KT1S1aZU5ATcWRARcq3mVtR9Z5M9ajjjwtv5 | 14.532817 |  361.818692 | 18.090934 |  343.727758 |
| KT1VHMARaXX374QLg5rMg9Jrm83mWH2AguTa |  0.006282 |    0.156412 |  0.007820 |    0.148592 |
| KT1U4N84AGXgU8JPuiuNujjXqJjKNV4XXUNK |  0.025131 |    0.625674 |  0.031283 |    0.594391 |
| KT1W5soiJhwuLaG6eYjhjZPCZfikGMJjSzWE | 30.670178 |  763.585188 | 38.179260 |  725.405928 |
+--------------------------------------+-----------+-------------+-----------+-------------+
|                                          TOTAL   | 1126.185966 | 56.309297 | 1069.876669 |
+--------------------------------------+-----------+-------------+-----------+-------------+
```

### Payout

#### Help

```
payman payout --help
Payout pays out rewards to delegations for the delegate passed.

Usage:
  payman payout [flags]

Flags:
  -c, --cycle int                  cycle to payout for (e.g. 95)
  -d, --delegate string            public key hash of the delegate that's paying out (e.g. --delegate=<phk>)
  -f, --fee float32                fee for the delegate (e.g. 0.05 = 5%) (default -1)
      --gas-limit int              network gas limit for each transaction in mutez (default 10200)(e.g. 10300) (default 10200)
  -h, --help                       help for payout
  -l, --log-file string            file to log to (default stdout)(e.g. ./payman.log) (default "/dev/stdout")
      --network-fee int            network fee for each transaction in mutez (default 1270)(e.g. 2000) (default 1270)
  -u, --node string                address to the node to query (default http://127.0.0.1:8732)(e.g. https://mainnet-node.tzscan.io:443) (default "http://127.0.0.1:8732")
  -k, --password string            password to the secret key of the wallet paying (e.g. --password=<passwd>)
      --payments-override string   overrides the rewards calculation and allows you to pass in your own payments in a json file (e.g. path/to/my/file/payments.json)
      --payout-min int             will only payout to addresses that meet the payout minimum (e.g. --payout-min=<mutez>)
  -r, --reddit string              path to reddit agent file (initiates reddit bot)(e.g. https://turnage.gitbooks.io/graw/content/chapter1.html)
      --reddit-title string        pre title for the reddit bot to post (e.g. DefinitelyNotABot: -- will read DefinitelyNotABot: Payout for Cycle <cycle>)
  -s, --secret string              encrypted secret key of the wallet paying (e.g. --secret=<sk>)
      --serve                      run service to payout for all new cycles going foward (default false)(e.g. --serve)
  -t, --twitter                    turn on twitter bot, will look for api keys in twitter.yml in current dir or --twitter-path (e.g. --twitter)
      --twitter-path string        path to twitter.yml file containing API keys if not in current dir (e.g. path/to/my/file/)
      --twitter-title string       pre title for the twitter bot to post (e.g. DefinitelyNotABot: -- will read DefinitelyNotABot: Payout for Cycle <cycle>)
```

#### Generic Example

```
payman payout --delegate=tz1SF9wBoBQbFUF13agZ8EgihLCKM54G1ccV --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm --password=abcd1234 --cycle=184 --fee=0.05 --payout-min=5000

[payout][preflight] warning: no network fee passed for payout, using default 1270 mutez
[payout][preflight] warning: no gas limit passed for payout, using default 10200 mutez
2019/05/20 18:50:33 reporting.go:24: Successful operation: "oorNRgL3WoQ49Z57pWTN62jVjF94cPAra1YBResAu6HyDx3SKkH"

+--------------------------------------+-----------+-------------+-----------+-------------+
|               ADDRESS                |   SHARE   |    GROSS    |    FEE    |     NET     |
+--------------------------------------+-----------+-------------+-----------+-------------+
| KT1U4N84AGXgU8JPuiuNujjXqJjKNV4XXUNK |  0.025131 |    0.625674 |  0.031283 |    0.594391 |
| KT1VHMARaXX374QLg5rMg9Jrm83mWH2AguTa |  0.006282 |    0.156412 |  0.007820 |    0.148592 |
| KT1S1aZU5ATcWRARcq3mVtR9Z5M9ajjjwtv5 | 14.532817 |  361.818692 | 18.090934 |  343.727758 |
| KT1W5soiJhwuLaG6eYjhjZPCZfikGMJjSzWE | 30.670178 |  763.585188 | 38.179260 |  725.405928 |
+--------------------------------------+-----------+-------------+-----------+-------------+
|                                          TOTAL   | 1126.185966 | 56.309297 | 1069.876669 |
```

#### Override Payments Example

This will override payman's calculations with your own by creating a file (e.g. payments.json) in the following format:

```
[
  {
      "Address": "KT1W5soiJhwuLaG6eYjhjZPCZfikGMJjSzWE",
      "Amount": 562508162
    },
  {
      "Address": "KT1S1aZU5ATcWRARcq3mVtR9Z5M9ajjjwtv5",
      "Amount": 267494981
    }
]
```

```
payman payout --delegate=tz1SF9wBoBQbFUF13agZ8EgihLCKM54G1ccV --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm --password=abcd1234 --payments-override=./payments.json

[payout][preflight] warning: no network fee passed for payout, using default 1270 mutez
[payout][preflight] warning: no gas limit passed for payout, using default 10200 mutez
2019/05/20 18:55:56 reporting.go:24: Successful operation: "onyZi9q84fMZQ53VxqqmfMDXukxb59bxNmvnjuKUWtD2SzTfdht"
```

#### Reddit Bot Example

This feature is currently only functional with mainnet. If used with another network, the link in your reddit post will be broken (Future Fix)

```
-r, --reddit string     example https://turnage.gitbooks.io/graw/content/chapter1.html
--title string      example "MyService:"
```

To post the results of your payout to reddit, create a [reddit.agent](https://turnage.gitbooks.io/graw/content/chapter1.html) file containing your reddit client and secret.

```
user_agent: "<platform>:<app ID>:<version string> (by /u/<reddit username>)"
client_id: "client id (looks kind of like: sdkfbwi48rhijwsdn)"
client_secret: "client secret (looks kind of like: ldkvblwiu34y8hsldjivn)"
username: "reddit username"
password: "reddit password"
```

Pass that file to the `--reddit` / `-r` flag. You may also include an option pre title using the `--reddit-title` flag.

```
payman payout --delegate=tz1SF9wBoBQbFUF13agZ8EgihLCKM54G1ccV --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm --password=abcd1234 --cycle=184 --fee=0.05 --payout-min=5000 --reddit=./reddit.agent --reddit-title=DefinitelyNotABot:

```

Your reddit post will contain the title passed along with "Payout for Cycle (cycle)"

```
DefinitelyNotATestBot: Payout for Cycle 100
```

With a link to the tzscan operation related to the cycle.

#### Twitter Bot Example

This feature is currently only functional with mainnet. If used with another network, the link in your twitter post will be broken (Future Fix)

```
-t, --twitter                turn on twitter bot, will look for api keys in twitter.yml in current dir or --twitter-path (e.g. --twitter)
      --twitter-path string    path to twitter.yml file containing API keys if not in current dir (e.g. path/to/my/file/)
      --twitter-title string   pre title for the twitter bot to post (e.g. DefinitelyNotABot: -- will read DefinitelyNotABot: Payout for Cycle <cycle>)
```

Pass the `--twitter` flag to turn on the twitter bot, and `--twitter-title` to add a custom title for your tweets. Payman will look for a twitter.yml file containing your api key, access token, and secrets. If this file is not in the directory of where payman is executed, you need to specify the path to find the file with `--twitter-path`.

```
payman payout --delegate=tz1SF9wBoBQbFUF13agZ8EgihLCKM54G1ccV --secret=edesk1Qx5JbctVnFVHL4A7BXgyExihHfcAHRYXoxkbSBmKqP2Sp92Gg1xcU8mqqu4Qi9TXkXwomMxAfy19sWAgCm --password=abcd1234 --cycle=184 --fee=0.05 --payout-min=5000 --twitter --twitter-title="DefinitelyNotATestBot:" --twitter-path=./

```

Example twitter.yml

```
consumerKey: "<key>"
consumerKeySecret: "<key_secret>"
accessToken: "<access_token>"
accessTokenSecret: "<access_token_secret>"
```

Your tweet will contain the title passed along with "Payout for Cycle (cycle)"

```
DefinitelyNotATestBot: Payout for Cycle 100
```

With a link to the tzscan operation related to the cycle.

## Roadmap:

- blacklist addresses
- defer option that allows you to defer a percentage of your fee to another address
- inlucde fiat price of XTZ in reports
- tax reporting

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
