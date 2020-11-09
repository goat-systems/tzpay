# Tzpay

Tzpay is a golang driven payout tool for delegation services on the tezos network. It is built with [go-tezos](https://github.com/goat-systems/go-tezos), and has support for [Dexter](http://camlcase.io/) liquidity providers. 

## Installation

### Docker
```
docker pull goatsystems/tzpay:latest

docker run --rm -ti goatsystems/tzpay:latest tzpay [command] \
-e TZPAY_BAKER=<TODO (e.g. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc)> \
-e TZPAY_BAKER_FEE=<TODO (e.g. 0.05 for 5%)> \
-e TZPAY_WALLET_ESK=<TODO (e.g. edesk...)> \
-e TZPAY_WALLET_PASSWORD=<TODO (e.g. password)>
```

## Configuration

| ENV                                  | Description                                          | Default                       | Required |
|--------------------------------------|------------------------------------------------------|:-----------------------------:|:--------:|
| TZPAY_BAKER                          | Pkh/Address of Baker                                 | N/A                           | True     |
| TZPAY_BAKER_FEE                      | Baker's Fee as a decimal (e.g. 5% would be 0.05)     | N/A                           | True     |
| TZPAY_WALLET_ESK                     | The tezos encrypted secret key (ed25519)             | N/A                           | True     |
| TZPAY_WALLET_PASSWORD                | The password to the encrypted secret key (ed25519)   | N/A                           | True     |
| TZPAY_BAKER_MINIMUM_PAYMENT          | Amounts below this amount will not be paid (MUTEZ)   | N/A                           | False    |
| TZPAY_BAKER_EARNINGS_ONLY            | Baker will not pay for missed endorsements or blocks | False                         | False    |
| TZPAY_BAKER_BLACK_LIST               | Baker will not pay addresses in blacklist            | N/A                           | False    |
| TZPAY_REWARDS_UNFROZEN_WAIT          | Baker pays out when rewards are unfrozen (tzpay serv)| False                         | False    |
| TZPAY_BAKER_LIQUIDITY_CONTRACTS_ONLY | Pays only liquidity providers                        | N/A                           | False    |
| TZPAY_BAKER_LIQUIDITY_CONTRACTS      | Pays liquidity providers in listed dexter contracts  | N/A                           | False    |
| TZPAY_API_TZKT                       | URL to a [tzkt api](api.tzkt.io)                     | https://api.tzkt.io           | False    |
| TZPAY_API_TEZOS                      | URL to a tezos RPC                                   | https://tezos.giganode.io/    | False    |
| TZPAY_OPERATIONS_NETWORK_FEE         | The network fee used in each transfer operation      | 2941                          | False    |
| TZPAY_OPERATIONS_GAS_LIMIT           | The gas limit used in each transfer operation        | 26283                         | False    |
| TZPAY_OPERATIONS_BATCH_SIZE          | The amount of transfers to include in an operation   | 125                           | False    |
| TZPAY_TWITTER_CONSUMER_KEY           | Twitter credentials for notifications                | N/A                           | False    |
| TZPAY_TWITTER_CONSUMER_SECRET        | Twitter credentials for notifications                | N/A                           | False    |
| TZPAY_TWITTER_ACCESS_TOKEN           | Twitter credentials for notifications                | N/A                           | False    |
| TZPAY_TWITTER_ACCESS_SECRET          | Twitter credentials for notifications                | N/A                           | False    |
| TZPAY_TWILIO_ACCOUNT_SID             | Twilio credentials for notifications                 | N/A                           | False    |
| TZPAY_TWILIO_AUTH_TOKEN              | Twilio credentials for notifications                 | N/A                           | False    |
| TZPAY_TWILIO_FROM                    | Twilio credentials for notifications                 | N/A                           | False    |
| TZPAY_TWILIO_TO                      | Twilio credentials for notifications                 | N/A                           | False    |

### Keys
As of now only ed25519 is supported.

### Notifications
If twilio or twitter credentials are provided, a notification will be sent after ever payout. 

### Help
```
➜  tzpay git:(dexter) ✗ ./tzpay help
A bulk payout tool for bakers in the Tezos Ecosystem

Usage:
  tzpay [command]

Available Commands:
  dryrun      dryrun simulates a payout
  help        Help about any command
  run         run executes a batch payout
  serv        serv runs a service that will continously payout cycle by cycle
  setup       setup prints a list of enviroment variables needed to get started.
  version     version prints tzpay's version

Flags:
  -h, --help   help for tzpay

Use "tzpay [command] --help" for more information about a command.
```

### Dryrun
```
➜  tzpay git:(dexter) ✗ ./tzpay dryrun 276 --table
+-------+--------------------------------------+----------+-----------+-----------+-----------+------------+
| CYLCE |                BAKER                 |  SHARE   |  REWARDS  |   FEES    |   TOTAL   | OPERATIONS |
+-------+--------------------------------------+----------+-----------+-----------+-----------+------------+
|   276 | tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc | 0.188334 | 73.550696 | 14.207660 | 87.758356 | N/A        |
+-------+--------------------------------------+----------+-----------+-----------+-----------+------------+
+--------------------------------------+----------+-----------+------------+-----------+
|              DELEGATION              |  SHARE   |   GROSS   |    NET     |    FEE    |
+--------------------------------------+----------+-----------+------------+-----------+
| tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd | 0.089355 | 34.895929 |  33.151133 |  1.744796 |
| KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy | 0.088661 | 34.624892 |  32.893648 |  1.731244 |
| KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv | 0.084802 | 33.118063 |  31.462160 |  1.655903 |
| KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj | 0.057932 | 22.624225 |  21.493014 |  1.131211 |
| KT1GBWviYFdRiNkhwM7LfrKDgHWnpdxpURtx | 0.043759 | 17.089149 |  16.234692 |  0.854457 |
| tz1Ykmc29JfQvWnjWRPYTPUZBLW4gwa9YKUD | 0.040391 | 15.773809 |  14.985119 |  0.788690 |
| tz1Zuav4ZBoiYhn4btW4HSr7G7J4txGZjvbu | 0.040379 | 15.769315 |  14.980850 |  0.788465 |
| KT1HccFB3cn4BR2za9XMuU7Wht64omed2UW8 | 0.025476 |  9.949126 |   9.451670 |  0.497456 |
| KT1WBDsJhoRvsvsRCmJirz9AFhSySSvzTWVd | 0.020242 |  7.905065 |   7.509812 |  0.395253 |
| KT1VzTs5piA7kYQkkfA9QNApVqGq1h6eMuV4 | 0.019578 |  7.646006 |   7.263706 |  0.382300 |
| tz1MXhttCeSJYpF3QRmPkMLCfNZaVufEuJmJ | 0.017759 |  6.935557 |   6.588780 |  0.346777 |
| KT1J4WFQRV3942phzRrh87WDFWKrNVcDJTP9 | 0.016668 |  6.509257 |   6.183795 |  0.325462 |
| tz1isExUQANnFb9YPmuwYmMmpeGHZ6T3CUT6 | 0.015796 |  6.168820 |   5.860379 |  0.308441 |
| KT1AThmRzcn51NwMf25NFYTqawjVo62hWiCv | 0.013499 |  5.271847 |   5.008255 |  0.263592 |
| KT19ABG9KxbEz2GrdN6uhGfxLmMY7REikBN8 | 0.011525 |  4.500767 |   4.275729 |  0.225038 |
| KT1Wp4tXL6GUtABkikB68fT7SaPQY2UuFkuE | 0.010131 |  3.956295 |   3.758481 |  0.197814 |
| tz1hbXhPVUX1fC8hN7fALyaUpdoC6EMgqM2h | 0.009624 |  3.758571 |   3.570643 |  0.187928 |
| KT1RuTPgQ6kdpnE3Adnw7Hr2KFN45uC3BdBy | 0.008928 |  3.486576 |   3.312248 |  0.174328 |
| KT1BjtEUxd25wwdwGH432LoP6PskvUc2bEYV | 0.008383 |  3.273781 |   3.110092 |  0.163689 |
| KT1JPeGNVarLsPZnSb3hG5xMVmJJmmBnrnpT | 0.007865 |  3.071596 |   2.918017 |  0.153579 |
| tz1gxmCTN8BSwuPLghDydtDKTqnAKyD8QTv7 | 0.006211 |  2.425409 |   2.304139 |  0.121270 |
| KT1A1sZmBQS9oZnPePRwP3Jyzv41xEppxfbF | 0.006142 |  2.398495 |   2.278571 |  0.119924 |
| tz1hSWBt6DD7SRH2Tq1kGbsKXZLrE7XGSMeF | 0.005812 |  2.269701 |   2.156216 |  0.113485 |
| tz1Nc2Zux98dEKqUW9Q9pL5rfUeLALBJTWGR | 0.005655 |  2.208601 |   2.098171 |  0.110430 |
| KT1C28u6DWsBfXk3UMyGrd8zTUVMpsyvjxmp | 0.004116 |  1.607420 |   1.527049 |  0.080371 |
| KT1TDrRrdz6SLYLBw8ZDxLWwJpx7FVpC52bt | 0.004110 |  1.604921 |   1.524675 |  0.080246 |
| KT1JcnHjWpkFxaLYMQD2URL8XEeAFqshz2uf | 0.003617 |  1.412444 |   1.341822 |  0.070622 |
| tz1VJa3ZkVwMzLFkGKhjvvrtzjRrnCJzMSKK | 0.003364 |  1.313906 |   1.248211 |  0.065695 |
| KT1TS49jiXxrnwhoJzAvCzGZCXLJs3XV1k6C | 0.003176 |  1.240476 |   1.178453 |  0.062023 |
| KT18kTf8UujihcF46Zn3rsFdEYFL1ZNFnGY4 | 0.003169 |  1.237687 |   1.175803 |  0.061884 |
| tz1VESLfEAEwDEKhyLZJYXVoervFk5ABPUUD | 0.003161 |  1.234597 |   1.172868 |  0.061729 |
| KT19Aro5JcjKH7J7RA6sCRihPiBQzQED3oQC | 0.003153 |  1.231416 |   1.169846 |  0.061570 |
| KT1CQiyDJ3mMVDoEqLY8Fz1onFXo5ycp5BDN | 0.003152 |  1.231047 |   1.169495 |  0.061552 |
| KT1QB9UAT1okYfcPQLi4jBmZkYg7LHcepERV | 0.003151 |  1.230739 |   1.169203 |  0.061536 |
| KT1QLo7DzPZnYK2EhmWpejVUnFjQUuWFKHnc | 0.003151 |  1.230554 |   1.169027 |  0.061527 |
| KT1UVUasDXH6mg8NCzRRgqvcjMoDUpETYEzH | 0.003151 |  1.230431 |   1.168910 |  0.061521 |
| KT1Na4maJ99GE6CGA1vEocWXrKRmxmsVUaTi | 0.003151 |  1.230431 |   1.168910 |  0.061521 |
| KT1MX2TwjSBzPaSsBUeW2k9DKehpiuMGfFcL | 0.003086 |  1.205179 |   1.144921 |  0.060258 |
| KT1BXmBgMSViAViNyhvkb441e2RBFMiKdnj7 | 0.003009 |  1.174983 |   1.116234 |  0.058749 |
| tz1iuFXyNN7nPHyHkfsj2tfZdnkK9MMJfFf1 | 0.002989 |  1.167275 |   1.108912 |  0.058363 |
| tz1Wq6LVwpofZ6zqjMBuLyEU53hRMepqkXEr | 0.002985 |  1.165852 |   1.107560 |  0.058292 |
| KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt | 0.002929 |  1.144055 |   1.086853 |  0.057202 |
| tz1hZZn4rsHLXdgQ9d8Rne9CLo6VFo29uQ3m | 0.002893 |  1.129820 |   1.073329 |  0.056491 |
| tz1Tjpy1ibFhioZ3Y1R6N9zoW4EL54AFYph1 | 0.002887 |  1.127494 |   1.071120 |  0.056374 |
| tz1Qadi21BxpHAjtfSrF6p4t3qMC5K8Ucjsw | 0.002872 |  1.121439 |   1.065368 |  0.056071 |
| tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH | 0.001326 |  0.517868 |   0.491975 |  0.025893 |
| KT1Jw925NVi4FzTVohZk5iLqagnhJGDEQoTS | 0.001134 |  0.442960 |   0.420812 |  0.022148 |
| KT1AUmLjJnmHmiieXnWWTPqHA98s65EeN7Mx | 0.000614 |  0.239825 |   0.227834 |  0.011991 |
| KT1PpVsfyVhWYTpyUaYigdmq1Aiv7zArTFYp | 0.000544 |  0.212280 |   0.201666 |  0.010614 |
| KT1CySPLDUSYyJ9vqNCF2dGgit4Rw2yUNEcj | 0.000377 |  0.147345 |   0.139978 |  0.007367 |
| tz1bHq6bUmTrvdepLVgYawcgEiLeeCMh2QJA | 0.000310 |  0.121075 |   0.115022 |  0.006053 |
| tz1Un6mfQ4Xie6U1nqmnedhnjNPAhfWx9jii | 0.000219 |  0.085511 |   0.081236 |  0.004275 |
| tz1VeiAS5wvYgNdri6vwDUrctQ5XhhaXY3K9 | 0.000204 |  0.079776 |   0.075788 |  0.003988 |
| KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC | 0.000189 |  0.073745 |   0.070058 |  0.003687 |
| tz1Vcu87ZuUK2e8BcoCBUWUhu2s2hPAabStm | 0.000156 |  0.060836 |   0.057795 |  0.003041 |
| KT1Aeg9D8kvkbAb6yikUdFcroReXvHtMBaZz | 0.000141 |  0.055183 |   0.052424 |  0.002759 |
| tz1aX2DF3ioDjqDcTVmrxVuqkxhZh1pLtfHU | 0.000095 |  0.037188 |   0.035329 |  0.001859 |
| KT18ni9Yar4UzwZozFbRF7SFUKg2EqyyUPPT | 0.000095 |  0.037070 |   0.035217 |  0.001853 |
| KT193c72q6eP1VpaY7hiheE7k1eDZiXeQUUw | 0.000086 |  0.033729 |   0.032043 |  0.001686 |
| tz1SnvfwMUYfD2uJrHBiaj4XPstW3eUE9RJU | 0.000079 |  0.031012 |   0.029462 |  0.001550 |
| KT1MSFeAGaWk8w7F1gmgUMaarU7mH385ueYC | 0.000041 |  0.015841 |   0.015049 |  0.000792 |
| KT1VUbpty8fER7npuvsfYDZXf2wVPhAHVqSx | 0.000036 |  0.014229 |   0.013518 |  0.000711 |
| KT1Lm4ZSyXSHod7U6znR7z9SGVmexntNQwAp | 0.000032 |  0.012472 |   0.011849 |  0.000623 |
| KT1NGd6RaRtmvwexYXGibtdvKBnNjjpBNknn | 0.000022 |  0.008518 |   0.008093 |  0.000425 |
| KT1WQWXvRcMjJB1y6mYZytoS5QsFJyFNDCk5 | 0.000021 |  0.008019 |   0.007619 |  0.000400 |
| tz1dfUssfLfTBoYqsWxMu86ycmLUvfF2abng | 0.000006 |  0.002377 |   0.002259 |  0.000118 |
| KT1JsHBFpoGRVXpcfC763YwvonKtNvaFotpG | 0.000006 |  0.002375 |   0.002257 |  0.000118 |
| tz1RomjUZ1j9F2vqE24h2Am8UeGUpcrf6vvJ | 0.000005 |  0.001982 |   0.001883 |  0.000099 |
| KT1Re5utTU2hrujXgZ3Ux5BgjN8rbru4sns2 | 0.000005 |  0.001949 |   0.001852 |  0.000097 |
| KT1AT7N9bGhViSorUrpivuYT6Wxs37hR2p9d | 0.000004 |  0.001569 |   0.001491 |  0.000078 |
| tz1a7ZrvfMm8reWSBQHcnAdjh9T5cXiu6EUT | 0.000004 |  0.001438 |   0.001367 |  0.000071 |
| KT1REp3D8dkiVVi37TCSMJNgGeX6UigBtfaL | 0.000004 |  0.001379 |   0.001311 |  0.000068 |
| KT18uqwoNyPRHpHCrg7xBFd7CiAZMbS1Ffne | 0.000002 |  0.000898 |   0.000854 |  0.000044 |
| KT1RbwPHzDwU9oPjnTWZrbCrMGjaFyj8dEtC | 0.000002 |  0.000836 |   0.000795 |  0.000041 |
| KT1JJcydTkinquNqh6kE5HYgFpD2124qHbZp | 0.000001 |  0.000315 |   0.000300 |  0.000015 |
| KT1JoAP7MfiigepR332u6xJqza9CG52ycYZ9 | 0.000000 |  0.000185 |   0.000176 |  0.000009 |
| KT1NfMCxyzwev243rKk3Y6SN8GfmdLKwASFQ | 0.000000 |  0.000184 |   0.000175 |  0.000009 |
| KT1EidADxWfYeBgK8L1ZTbf7a9zyjKwCFjfH | 0.000000 |  0.000178 |   0.000170 |  0.000008 |
| KT1XrBAocuiE3C2vvtgt7PFoazrC1KRi9ZF4 | 0.000000 |  0.000149 |   0.000142 |  0.000007 |
| KT1CeUNtCrXFNbLmvdGPNnxpcJw2sW5Hcpmc | 0.000000 |  0.000111 |   0.000106 |  0.000005 |
| tz1PB27kbPL64MWYoNZAfQAEmzCZFi9EvgBw | 0.000000 |  0.000103 |   0.000098 |  0.000005 |
| KT1T3dPMBm7D3kKqALKYnW2mViFqMMVCYtmo | 0.000000 |  0.000099 |   0.000095 |  0.000004 |
| KT1Dgma8bbDtAbtMbYYS5VmziyCANAZn8M7W | 0.000000 |  0.000096 |   0.000092 |  0.000004 |
| KT1NmVtU3CNqzhNWwLhE5BqAopjkcmHpWzT2 | 0.000000 |  0.000093 |   0.000089 |  0.000004 |
| KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3 | 0.000000 |  0.000077 |   0.000074 |  0.000003 |
| KT19Q8GiYqGpuuUjf9xfXXVu1WY889N8oxRe | 0.000000 |  0.000062 |   0.000059 |  0.000003 |
| KT1S9VbEnU8nj33ufxrGBYGxBCnqmeoAnKt4 | 0.000000 |  0.000055 |   0.000053 |  0.000002 |
| KT1Lnh39om2iqr4qb9AarF9T38ayNBLnfAVn | 0.000000 |  0.000039 |   0.000038 |  0.000001 |
| KT1Cz1jPLuaPR99XamKQDr9PKZY1PTXzTAHH | 0.000000 |  0.000038 |   0.000037 |  0.000001 |
| KT1PDBuQmFLVHfiWZjV248QdTrdcmAuSS7Tx | 0.000000 |  0.000032 |   0.000031 |  0.000001 |
| KT1EbMbqTUS8XnqGVRsdLZVKLhcT7Zc33jR1 | 0.000000 |  0.000020 |   0.000019 |  0.000001 |
| KT1KJ5Qt18yU9DrqN36tgyLtaSvFSZ5r6YL6 | 0.000000 |  0.000014 |   0.000014 |  0.000000 |
| KT1PY2MMiTUkZQv7CPekXy186N1qmu7GikcT | 0.000000 |  0.000012 |   0.000012 |  0.000000 |
| KT1E1MnvNgCDLqnGneStVY9CvmjnyZgcaPaD | 0.000000 |  0.000010 |   0.000010 |  0.000000 |
| tz1f7mbrPU2cMHhjqhYzw9SfmZYKUtZkG52A | 0.000000 |  0.000009 |   0.000009 |  0.000000 |
| KT1KeNNxEM4NyfrmF1CG6TLn3nRSmEGhP7Z2 | 0.000000 |  0.000008 |   0.000008 |  0.000000 |
| KT1W3oiS6s9NgSxhZY1nCsazW2QbwkmjkET1 | 0.000000 |  0.000005 |   0.000005 |  0.000000 |
| KT1NxnFWHW7bUxzks1oHVU2jn4heu48KC3eD | 0.000000 |  0.000004 |   0.000004 |  0.000000 |
| KT1MfT8XvQp9ZeGUx4cmCNF3wui55WLNYhq9 | 0.000000 |  0.000002 |   0.000002 |  0.000000 |
| KT1XPMJx2wuCbbzKZx5jJyKqLpPJMHv58wni | 0.000000 |  0.000000 |   0.000000 |  0.000000 |
+--------------------------------------+----------+-----------+------------+-----------+
|                                                     TOTAL   | 269.946543 | 14.207660 |
+--------------------------------------+----------+-----------+------------+-----------+
```

### Run
```
➜  tzpay git:(dexter) ✗ ./tzpay dryrun 276 --table
+-------+--------------------------------------+----------+-----------+-----------+-----------+------------+
| CYLCE |                BAKER                 |  SHARE   |  REWARDS  |   FEES    |   TOTAL   | OPERATIONS |
+-------+--------------------------------------+----------+-----------+-----------+-----------+------------+
|   276 | tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc | 0.188334 | 73.550696 | 14.207660 | 87.758356 | N/A        |
+-------+--------------------------------------+----------+-----------+-----------+-----------+------------+
+--------------------------------------+----------+-----------+------------+-----------+
|              DELEGATION              |  SHARE   |   GROSS   |    NET     |    FEE    |
+--------------------------------------+----------+-----------+------------+-----------+
| tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd | 0.089355 | 34.895929 |  33.151133 |  1.744796 |
| KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy | 0.088661 | 34.624892 |  32.893648 |  1.731244 |
| KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv | 0.084802 | 33.118063 |  31.462160 |  1.655903 |
| KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj | 0.057932 | 22.624225 |  21.493014 |  1.131211 |
| KT1GBWviYFdRiNkhwM7LfrKDgHWnpdxpURtx | 0.043759 | 17.089149 |  16.234692 |  0.854457 |
| tz1Ykmc29JfQvWnjWRPYTPUZBLW4gwa9YKUD | 0.040391 | 15.773809 |  14.985119 |  0.788690 |
| tz1Zuav4ZBoiYhn4btW4HSr7G7J4txGZjvbu | 0.040379 | 15.769315 |  14.980850 |  0.788465 |
| KT1HccFB3cn4BR2za9XMuU7Wht64omed2UW8 | 0.025476 |  9.949126 |   9.451670 |  0.497456 |
| KT1WBDsJhoRvsvsRCmJirz9AFhSySSvzTWVd | 0.020242 |  7.905065 |   7.509812 |  0.395253 |
| KT1VzTs5piA7kYQkkfA9QNApVqGq1h6eMuV4 | 0.019578 |  7.646006 |   7.263706 |  0.382300 |
| tz1MXhttCeSJYpF3QRmPkMLCfNZaVufEuJmJ | 0.017759 |  6.935557 |   6.588780 |  0.346777 |
| KT1J4WFQRV3942phzRrh87WDFWKrNVcDJTP9 | 0.016668 |  6.509257 |   6.183795 |  0.325462 |
| tz1isExUQANnFb9YPmuwYmMmpeGHZ6T3CUT6 | 0.015796 |  6.168820 |   5.860379 |  0.308441 |
| KT1AThmRzcn51NwMf25NFYTqawjVo62hWiCv | 0.013499 |  5.271847 |   5.008255 |  0.263592 |
| KT19ABG9KxbEz2GrdN6uhGfxLmMY7REikBN8 | 0.011525 |  4.500767 |   4.275729 |  0.225038 |
| KT1Wp4tXL6GUtABkikB68fT7SaPQY2UuFkuE | 0.010131 |  3.956295 |   3.758481 |  0.197814 |
| tz1hbXhPVUX1fC8hN7fALyaUpdoC6EMgqM2h | 0.009624 |  3.758571 |   3.570643 |  0.187928 |
| KT1RuTPgQ6kdpnE3Adnw7Hr2KFN45uC3BdBy | 0.008928 |  3.486576 |   3.312248 |  0.174328 |
| KT1BjtEUxd25wwdwGH432LoP6PskvUc2bEYV | 0.008383 |  3.273781 |   3.110092 |  0.163689 |
| KT1JPeGNVarLsPZnSb3hG5xMVmJJmmBnrnpT | 0.007865 |  3.071596 |   2.918017 |  0.153579 |
| tz1gxmCTN8BSwuPLghDydtDKTqnAKyD8QTv7 | 0.006211 |  2.425409 |   2.304139 |  0.121270 |
| KT1A1sZmBQS9oZnPePRwP3Jyzv41xEppxfbF | 0.006142 |  2.398495 |   2.278571 |  0.119924 |
| tz1hSWBt6DD7SRH2Tq1kGbsKXZLrE7XGSMeF | 0.005812 |  2.269701 |   2.156216 |  0.113485 |
| tz1Nc2Zux98dEKqUW9Q9pL5rfUeLALBJTWGR | 0.005655 |  2.208601 |   2.098171 |  0.110430 |
| KT1C28u6DWsBfXk3UMyGrd8zTUVMpsyvjxmp | 0.004116 |  1.607420 |   1.527049 |  0.080371 |
| KT1TDrRrdz6SLYLBw8ZDxLWwJpx7FVpC52bt | 0.004110 |  1.604921 |   1.524675 |  0.080246 |
| KT1JcnHjWpkFxaLYMQD2URL8XEeAFqshz2uf | 0.003617 |  1.412444 |   1.341822 |  0.070622 |
| tz1VJa3ZkVwMzLFkGKhjvvrtzjRrnCJzMSKK | 0.003364 |  1.313906 |   1.248211 |  0.065695 |
| KT1TS49jiXxrnwhoJzAvCzGZCXLJs3XV1k6C | 0.003176 |  1.240476 |   1.178453 |  0.062023 |
| KT18kTf8UujihcF46Zn3rsFdEYFL1ZNFnGY4 | 0.003169 |  1.237687 |   1.175803 |  0.061884 |
| tz1VESLfEAEwDEKhyLZJYXVoervFk5ABPUUD | 0.003161 |  1.234597 |   1.172868 |  0.061729 |
| KT19Aro5JcjKH7J7RA6sCRihPiBQzQED3oQC | 0.003153 |  1.231416 |   1.169846 |  0.061570 |
| KT1CQiyDJ3mMVDoEqLY8Fz1onFXo5ycp5BDN | 0.003152 |  1.231047 |   1.169495 |  0.061552 |
| KT1QB9UAT1okYfcPQLi4jBmZkYg7LHcepERV | 0.003151 |  1.230739 |   1.169203 |  0.061536 |
| KT1QLo7DzPZnYK2EhmWpejVUnFjQUuWFKHnc | 0.003151 |  1.230554 |   1.169027 |  0.061527 |
| KT1UVUasDXH6mg8NCzRRgqvcjMoDUpETYEzH | 0.003151 |  1.230431 |   1.168910 |  0.061521 |
| KT1Na4maJ99GE6CGA1vEocWXrKRmxmsVUaTi | 0.003151 |  1.230431 |   1.168910 |  0.061521 |
| KT1MX2TwjSBzPaSsBUeW2k9DKehpiuMGfFcL | 0.003086 |  1.205179 |   1.144921 |  0.060258 |
| KT1BXmBgMSViAViNyhvkb441e2RBFMiKdnj7 | 0.003009 |  1.174983 |   1.116234 |  0.058749 |
| tz1iuFXyNN7nPHyHkfsj2tfZdnkK9MMJfFf1 | 0.002989 |  1.167275 |   1.108912 |  0.058363 |
| tz1Wq6LVwpofZ6zqjMBuLyEU53hRMepqkXEr | 0.002985 |  1.165852 |   1.107560 |  0.058292 |
| KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt | 0.002929 |  1.144055 |   1.086853 |  0.057202 |
| tz1hZZn4rsHLXdgQ9d8Rne9CLo6VFo29uQ3m | 0.002893 |  1.129820 |   1.073329 |  0.056491 |
| tz1Tjpy1ibFhioZ3Y1R6N9zoW4EL54AFYph1 | 0.002887 |  1.127494 |   1.071120 |  0.056374 |
| tz1Qadi21BxpHAjtfSrF6p4t3qMC5K8Ucjsw | 0.002872 |  1.121439 |   1.065368 |  0.056071 |
| tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH | 0.001326 |  0.517868 |   0.491975 |  0.025893 |
| KT1Jw925NVi4FzTVohZk5iLqagnhJGDEQoTS | 0.001134 |  0.442960 |   0.420812 |  0.022148 |
| KT1AUmLjJnmHmiieXnWWTPqHA98s65EeN7Mx | 0.000614 |  0.239825 |   0.227834 |  0.011991 |
| KT1PpVsfyVhWYTpyUaYigdmq1Aiv7zArTFYp | 0.000544 |  0.212280 |   0.201666 |  0.010614 |
| KT1CySPLDUSYyJ9vqNCF2dGgit4Rw2yUNEcj | 0.000377 |  0.147345 |   0.139978 |  0.007367 |
| tz1bHq6bUmTrvdepLVgYawcgEiLeeCMh2QJA | 0.000310 |  0.121075 |   0.115022 |  0.006053 |
| tz1Un6mfQ4Xie6U1nqmnedhnjNPAhfWx9jii | 0.000219 |  0.085511 |   0.081236 |  0.004275 |
| tz1VeiAS5wvYgNdri6vwDUrctQ5XhhaXY3K9 | 0.000204 |  0.079776 |   0.075788 |  0.003988 |
| KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC | 0.000189 |  0.073745 |   0.070058 |  0.003687 |
| tz1Vcu87ZuUK2e8BcoCBUWUhu2s2hPAabStm | 0.000156 |  0.060836 |   0.057795 |  0.003041 |
| KT1Aeg9D8kvkbAb6yikUdFcroReXvHtMBaZz | 0.000141 |  0.055183 |   0.052424 |  0.002759 |
| tz1aX2DF3ioDjqDcTVmrxVuqkxhZh1pLtfHU | 0.000095 |  0.037188 |   0.035329 |  0.001859 |
| KT18ni9Yar4UzwZozFbRF7SFUKg2EqyyUPPT | 0.000095 |  0.037070 |   0.035217 |  0.001853 |
| KT193c72q6eP1VpaY7hiheE7k1eDZiXeQUUw | 0.000086 |  0.033729 |   0.032043 |  0.001686 |
| tz1SnvfwMUYfD2uJrHBiaj4XPstW3eUE9RJU | 0.000079 |  0.031012 |   0.029462 |  0.001550 |
| KT1MSFeAGaWk8w7F1gmgUMaarU7mH385ueYC | 0.000041 |  0.015841 |   0.015049 |  0.000792 |
| KT1VUbpty8fER7npuvsfYDZXf2wVPhAHVqSx | 0.000036 |  0.014229 |   0.013518 |  0.000711 |
| KT1Lm4ZSyXSHod7U6znR7z9SGVmexntNQwAp | 0.000032 |  0.012472 |   0.011849 |  0.000623 |
| KT1NGd6RaRtmvwexYXGibtdvKBnNjjpBNknn | 0.000022 |  0.008518 |   0.008093 |  0.000425 |
| KT1WQWXvRcMjJB1y6mYZytoS5QsFJyFNDCk5 | 0.000021 |  0.008019 |   0.007619 |  0.000400 |
| tz1dfUssfLfTBoYqsWxMu86ycmLUvfF2abng | 0.000006 |  0.002377 |   0.002259 |  0.000118 |
| KT1JsHBFpoGRVXpcfC763YwvonKtNvaFotpG | 0.000006 |  0.002375 |   0.002257 |  0.000118 |
| tz1RomjUZ1j9F2vqE24h2Am8UeGUpcrf6vvJ | 0.000005 |  0.001982 |   0.001883 |  0.000099 |
| KT1Re5utTU2hrujXgZ3Ux5BgjN8rbru4sns2 | 0.000005 |  0.001949 |   0.001852 |  0.000097 |
| KT1AT7N9bGhViSorUrpivuYT6Wxs37hR2p9d | 0.000004 |  0.001569 |   0.001491 |  0.000078 |
| tz1a7ZrvfMm8reWSBQHcnAdjh9T5cXiu6EUT | 0.000004 |  0.001438 |   0.001367 |  0.000071 |
| KT1REp3D8dkiVVi37TCSMJNgGeX6UigBtfaL | 0.000004 |  0.001379 |   0.001311 |  0.000068 |
| KT18uqwoNyPRHpHCrg7xBFd7CiAZMbS1Ffne | 0.000002 |  0.000898 |   0.000854 |  0.000044 |
| KT1RbwPHzDwU9oPjnTWZrbCrMGjaFyj8dEtC | 0.000002 |  0.000836 |   0.000795 |  0.000041 |
| KT1JJcydTkinquNqh6kE5HYgFpD2124qHbZp | 0.000001 |  0.000315 |   0.000300 |  0.000015 |
| KT1JoAP7MfiigepR332u6xJqza9CG52ycYZ9 | 0.000000 |  0.000185 |   0.000176 |  0.000009 |
| KT1NfMCxyzwev243rKk3Y6SN8GfmdLKwASFQ | 0.000000 |  0.000184 |   0.000175 |  0.000009 |
| KT1EidADxWfYeBgK8L1ZTbf7a9zyjKwCFjfH | 0.000000 |  0.000178 |   0.000170 |  0.000008 |
| KT1XrBAocuiE3C2vvtgt7PFoazrC1KRi9ZF4 | 0.000000 |  0.000149 |   0.000142 |  0.000007 |
| KT1CeUNtCrXFNbLmvdGPNnxpcJw2sW5Hcpmc | 0.000000 |  0.000111 |   0.000106 |  0.000005 |
| tz1PB27kbPL64MWYoNZAfQAEmzCZFi9EvgBw | 0.000000 |  0.000103 |   0.000098 |  0.000005 |
| KT1T3dPMBm7D3kKqALKYnW2mViFqMMVCYtmo | 0.000000 |  0.000099 |   0.000095 |  0.000004 |
| KT1Dgma8bbDtAbtMbYYS5VmziyCANAZn8M7W | 0.000000 |  0.000096 |   0.000092 |  0.000004 |
| KT1NmVtU3CNqzhNWwLhE5BqAopjkcmHpWzT2 | 0.000000 |  0.000093 |   0.000089 |  0.000004 |
| KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3 | 0.000000 |  0.000077 |   0.000074 |  0.000003 |
| KT19Q8GiYqGpuuUjf9xfXXVu1WY889N8oxRe | 0.000000 |  0.000062 |   0.000059 |  0.000003 |
| KT1S9VbEnU8nj33ufxrGBYGxBCnqmeoAnKt4 | 0.000000 |  0.000055 |   0.000053 |  0.000002 |
| KT1Lnh39om2iqr4qb9AarF9T38ayNBLnfAVn | 0.000000 |  0.000039 |   0.000038 |  0.000001 |
| KT1Cz1jPLuaPR99XamKQDr9PKZY1PTXzTAHH | 0.000000 |  0.000038 |   0.000037 |  0.000001 |
| KT1PDBuQmFLVHfiWZjV248QdTrdcmAuSS7Tx | 0.000000 |  0.000032 |   0.000031 |  0.000001 |
| KT1EbMbqTUS8XnqGVRsdLZVKLhcT7Zc33jR1 | 0.000000 |  0.000020 |   0.000019 |  0.000001 |
| KT1KJ5Qt18yU9DrqN36tgyLtaSvFSZ5r6YL6 | 0.000000 |  0.000014 |   0.000014 |  0.000000 |
| KT1PY2MMiTUkZQv7CPekXy186N1qmu7GikcT | 0.000000 |  0.000012 |   0.000012 |  0.000000 |
| KT1E1MnvNgCDLqnGneStVY9CvmjnyZgcaPaD | 0.000000 |  0.000010 |   0.000010 |  0.000000 |
| tz1f7mbrPU2cMHhjqhYzw9SfmZYKUtZkG52A | 0.000000 |  0.000009 |   0.000009 |  0.000000 |
| KT1KeNNxEM4NyfrmF1CG6TLn3nRSmEGhP7Z2 | 0.000000 |  0.000008 |   0.000008 |  0.000000 |
| KT1W3oiS6s9NgSxhZY1nCsazW2QbwkmjkET1 | 0.000000 |  0.000005 |   0.000005 |  0.000000 |
| KT1NxnFWHW7bUxzks1oHVU2jn4heu48KC3eD | 0.000000 |  0.000004 |   0.000004 |  0.000000 |
| KT1MfT8XvQp9ZeGUx4cmCNF3wui55WLNYhq9 | 0.000000 |  0.000002 |   0.000002 |  0.000000 |
| KT1XPMJx2wuCbbzKZx5jJyKqLpPJMHv58wni | 0.000000 |  0.000000 |   0.000000 |  0.000000 |
+--------------------------------------+----------+-----------+------------+-----------+
|                                                     TOTAL   | 269.946543 | 14.207660 |
+--------------------------------------+----------+-----------+------------+-----------+
```

## Roadmap:
* inlucde fiat price of XTZ in reports
* tax reporting

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
