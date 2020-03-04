# Tzpay

Tzpay is a golang driven payout tool for delegation services on the tezos network. Tzpay is built with [go-tezos](https://github.com/goat-systems/go-tezos).

## Installation

### Source 
```
go get -u github.com/goat-systems/tzpay
```

### Linux
```
wget https://github.com/DefinitelyNotAGoat/payman/releases/download/v2.0.0-alpha/tzpay_linux_amd64
sudo mv tzpay_linux_amd64 /usr/local/bin/tzpay
sudo chmod a+x /usr/local/bin/tzpay
```

### Docker
```
docker pull goatsystems/tzpay:latest

docker run --rm -ti goatsystems/tzpay:latest tzpay [command] \
-e TZPAY_HOST_NODE=<TODO (e.g. http://127.0.0.1:8732)> \
-e TZPAY_BAKERS_FEE=<TODO (e.g. 0.05 for 5%)> \
-e TZPAY_DELEGATE=<TODO (e.g. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc)> \
-e TZPAY_WALLET_SECRET=<TODO (e.g. edesk...)> \
-e TZPAY_WALLET_PASSWORD=<TODO (e.g. password)>
```

## Usage

### Required Enviroment Variables 
```
TZPAY_HOST_NODE=<TODO (e.g. http://127.0.0.1:8732)>
TZPAY_BAKERS_FEE=<TODO (e.g. 0.05 for 5%)>
TZPAY_DELEGATE=<TODO (e.g. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc)>
TZPAY_WALLET_SECRET=<TODO (e.g. edesk...)>
TZPAY_WALLET_PASSWORD=<TODO (e.g. password)>
```

### Optional Enviroment Variables
```
TZPAY_BLACKLIST=<TODO (e.g. tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH, tz1W3HW533csCBLor4NPtU79R2TT2sbKfjh7)>
TZPAY_NETWORK_GAS_LIMIT=<TODO (e.g. 30000)>
TZPAY_NETWORK_FEE=<TODO (e.g. 3000)>
TZPAY_MINIMUM_PAYMENT=<TODO (e.g. 3000)>
```

### Help
```
➜  tzpay git:(v2) ✗ ./tzpay --help
A bulk payout tool for bakers in the Tezos Ecosystem

Usage:
  tzpay [command]

Available Commands:
  dryrun      dryrun simulates a payout
  help        Help about any command
  run         run executes a batch payout
  setup       setup prints a list of enviroment variables needed to get started.
  version     version prints tzpay's version

Flags:
  -h, --help   help for tzpay

Use "tzpay [command] --help" for more information about a command.
```

### Dryrun
```
➜  tzpay git:(v2) ✗ ./tzpay dryrun 206 --table
+--------------------------------------+--------------------------------------+------------+-----------+
|                BAKER                 |                WALLET                |  REWARDS   | OPERATION |
+--------------------------------------+--------------------------------------+------------+-----------+
| tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc | tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH | 787.800000 | N/A       |
+--------------------------------------+--------------------------------------+------------+-----------+
+--------------------------------------+----------+------------+------------+-----------+
|              DELEGATION              |  SHARE   |   GROSS    |    NET     |    FEE    |
+--------------------------------------+----------+------------+------------+-----------+
| KT18kTf8UujihcF46Zn3rsFdEYFL1ZNFnGY4 | 0.001337 |   1.053357 |   1.000690 |  0.052667 |
| KT18ni9Yar4UzwZozFbRF7SFUKg2EqyyUPPT | 0.000040 |   0.031549 |   0.029972 |  0.001577 |
| KT18uqwoNyPRHpHCrg7xBFd7CiAZMbS1Ffne | 0.000460 |   0.362624 |   0.344493 |  0.018131 |
| KT193c72q6eP1VpaY7hiheE7k1eDZiXeQUUw | 0.000036 |   0.028706 |   0.027271 |  0.001435 |
| KT19ABG9KxbEz2GrdN6uhGfxLmMY7REikBN8 | 0.004862 |   3.830464 |   3.638941 |  0.191523 |
| KT19Aro5JcjKH7J7RA6sCRihPiBQzQED3oQC | 0.001330 |   1.048020 |   0.995619 |  0.052401 |
| KT19Q8GiYqGpuuUjf9xfXXVu1WY889N8oxRe | 0.000000 |   0.000053 |   0.000051 |  0.000002 |
| KT1A1sZmBQS9oZnPePRwP3Jyzv41xEppxfbF | 0.002591 |   2.041285 |   1.939221 |  0.102064 |
| KT1A5seo53aLSSyHgJKZFYnh7jTZBtFnNnjz | 0.000000 |   0.000002 |   0.000002 |  0.000000 |
| KT1Aeg9D8kvkbAb6yikUdFcroReXvHtMBaZz | 0.000060 |   0.046964 |   0.044616 |  0.002348 |
| KT1AT7N9bGhViSorUrpivuYT6Wxs37hR2p9d | 0.000002 |   0.001336 |   0.001270 |  0.000066 |
| KT1AThmRzcn51NwMf25NFYTqawjVo62hWiCv | 0.005695 |   4.486706 |   4.262371 |  0.224335 |
| KT1BjtEUxd25wwdwGH432LoP6PskvUc2bEYV | 0.003537 |   2.786213 |   2.646903 |  0.139310 |
| KT1BXmBgMSViAViNyhvkb441e2RBFMiKdnj7 | 0.001270 |   1.000411 |   0.950391 |  0.050020 |
| KT1C28u6DWsBfXk3UMyGrd8zTUVMpsyvjxmp | 0.002488 |   1.959845 |   1.861853 |  0.097992 |
| KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC | 0.029346 |  23.118649 |  21.962717 |  1.155932 |
| KT1CeUNtCrXFNbLmvdGPNnxpcJw2sW5Hcpmc | 0.000000 |   0.000095 |   0.000091 |  0.000004 |
| KT1CQiyDJ3mMVDoEqLY8Fz1onFXo5ycp5BDN | 0.001330 |   1.047706 |   0.995321 |  0.052385 |
| KT1CySPLDUSYyJ9vqNCF2dGgit4Rw2yUNEcj | 0.000159 |   0.125470 |   0.119197 |  0.006273 |
| KT1Cz1jPLuaPR99XamKQDr9PKZY1PTXzTAHH | 0.000000 |   0.000033 |   0.000032 |  0.000001 |
| KT1Dgma8bbDtAbtMbYYS5VmziyCANAZn8M7W | 0.000000 |   0.000083 |   0.000079 |  0.000004 |
| KT1E1MnvNgCDLqnGneStVY9CvmjnyZgcaPaD | 0.000000 |   0.000009 |   0.000009 |  0.000000 |
| KT1EbMbqTUS8XnqGVRsdLZVKLhcT7Zc33jR1 | 0.010737 |   8.458494 |   8.035570 |  0.422924 |
| KT1EidADxWfYeBgK8L1ZTbf7a9zyjKwCFjfH | 0.000000 |   0.000151 |   0.000144 |  0.000007 |
| KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy | 0.037406 |  29.468174 |  27.994766 |  1.473408 |
| KT1GBWviYFdRiNkhwM7LfrKDgHWnpdxpURtx | 0.018462 |  14.544046 |  13.816844 |  0.727202 |
| KT1GcSsQaTtMB2HvUKU9b6WRFUnGpGx9JwGk | 0.000000 |   0.000000 |   0.000000 |  0.000000 |
| KT1HccFB3cn4BR2za9XMuU7Wht64omed2UW8 | 0.010748 |   8.467393 |   8.044024 |  0.423369 |
| KT1J2uk1fYSnZjxkJcUhFDkaRDhjCTRBspqv | 0.000000 |   0.000000 |   0.000000 |  0.000000 |
| KT1J4WFQRV3942phzRrh87WDFWKrNVcDJTP9 | 0.007032 |   5.539827 |   5.262836 |  0.276991 |
| KT1JJcydTkinquNqh6kE5HYgFpD2124qHbZp | 0.042888 |  33.787274 |  32.097911 |  1.689363 |
| KT1JoAP7MfiigepR332u6xJqza9CG52ycYZ9 | 0.000000 |   0.000157 |   0.000150 |  0.000007 |
| KT1JPeGNVarLsPZnSb3hG5xMVmJJmmBnrnpT | 0.003366 |   2.652015 |   2.519415 |  0.132600 |
| KT1JsHBFpoGRVXpcfC763YwvonKtNvaFotpG | 0.000003 |   0.002021 |   0.001920 |  0.000101 |
| KT1Jw925NVi4FzTVohZk5iLqagnhJGDEQoTS | 0.000479 |   0.376989 |   0.358140 |  0.018849 |
| KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt | 0.001236 |   0.974078 |   0.925375 |  0.048703 |
| KT1KeNNxEM4NyfrmF1CG6TLn3nRSmEGhP7Z2 | 0.001643 |   1.294659 |   1.229927 |  0.064732 |
| KT1KJ5Qt18yU9DrqN36tgyLtaSvFSZ5r6YL6 | 0.002900 |   2.284941 |   2.170694 |  0.114247 |
| KT1LfoE9EbpczdzUzowRckGUfikGcd5PyVKg | 0.000000 |   0.000031 |   0.000030 |  0.000001 |
| KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv | 0.035778 |  28.185759 |  26.776472 |  1.409287 |
| KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3 | 0.000023 |   0.017797 |   0.016908 |  0.000889 |
| KT1Lm4ZSyXSHod7U6znR7z9SGVmexntNQwAp | 0.000013 |   0.010614 |   0.010084 |  0.000530 |
| KT1Lnh39om2iqr4qb9AarF9T38ayNBLnfAVn | 0.000000 |   0.000034 |   0.000033 |  0.000001 |
| KT1MfT8XvQp9ZeGUx4cmCNF3wui55WLNYhq9 | 0.000000 |   0.000002 |   0.000002 |  0.000000 |
| KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj | 0.024467 |  19.274821 |  18.311080 |  0.963741 |
| KT1MSFeAGaWk8w7F1gmgUMaarU7mH385ueYC | 0.005919 |   4.662630 |   4.429499 |  0.233131 |
| KT1MX2TwjSBzPaSsBUeW2k9DKehpiuMGfFcL | 0.001302 |   1.025690 |   0.974406 |  0.051284 |
| KT1Na4maJ99GE6CGA1vEocWXrKRmxmsVUaTi | 0.001329 |   1.047182 |   0.994823 |  0.052359 |
| KT1NfMCxyzwev243rKk3Y6SN8GfmdLKwASFQ | 0.000039 |   0.030682 |   0.029148 |  0.001534 |
| KT1NGd6RaRtmvwexYXGibtdvKBnNjjpBNknn | 0.000009 |   0.007249 |   0.006887 |  0.000362 |
| KT1NmVtU3CNqzhNWwLhE5BqAopjkcmHpWzT2 | 0.000028 |   0.022144 |   0.021037 |  0.001107 |
| KT1NxnFWHW7bUxzks1oHVU2jn4heu48KC3eD | 0.000000 |   0.000004 |   0.000004 |  0.000000 |
| KT1PDBuQmFLVHfiWZjV248QdTrdcmAuSS7Tx | 0.000000 |   0.000028 |   0.000027 |  0.000001 |
| KT1PY2MMiTUkZQv7CPekXy186N1qmu7GikcT | 0.000000 |   0.000010 |   0.000010 |  0.000000 |
| KT1QB9UAT1okYfcPQLi4jBmZkYg7LHcepERV | 0.001330 |   1.047444 |   0.995072 |  0.052372 |
| KT1QLo7DzPZnYK2EhmWpejVUnFjQUuWFKHnc | 0.001329 |   1.047287 |   0.994923 |  0.052364 |
| KT1RbwPHzDwU9oPjnTWZrbCrMGjaFyj8dEtC | 0.000001 |   0.000711 |   0.000676 |  0.000035 |
| KT1REp3D8dkiVVi37TCSMJNgGeX6UigBtfaL | 0.000001 |   0.001174 |   0.001116 |  0.000058 |
| KT1Re5utTU2hrujXgZ3Ux5BgjN8rbru4sns2 | 0.000574 |   0.452440 |   0.429818 |  0.022622 |
| KT1RuTPgQ6kdpnE3Adnw7Hr2KFN45uC3BdBy | 0.003767 |   2.967317 |   2.818952 |  0.148365 |
| KT1S9VbEnU8nj33ufxrGBYGxBCnqmeoAnKt4 | 0.000000 |   0.000048 |   0.000046 |  0.000002 |
| KT1T3dPMBm7D3kKqALKYnW2mViFqMMVCYtmo | 0.000028 |   0.022193 |   0.021084 |  0.001109 |
| KT1TDrRrdz6SLYLBw8ZDxLWwJpx7FVpC52bt | 0.001734 |   1.365962 |   1.297664 |  0.068298 |
| KT1TS49jiXxrnwhoJzAvCzGZCXLJs3XV1k6C | 0.001340 |   1.055731 |   1.002945 |  0.052786 |
| KT1Uh1G9tdq45N63ZBrreDKy7eZF8QVoydm1 | 0.000000 |   0.000024 |   0.000023 |  0.000001 |
| KT1UVUasDXH6mg8NCzRRgqvcjMoDUpETYEzH | 0.001329 |   1.047182 |   0.994823 |  0.052359 |
| KT1VUbpty8fER7npuvsfYDZXf2wVPhAHVqSx | 0.000004 |   0.002879 |   0.002736 |  0.000143 |
| KT1VzTs5piA7kYQkkfA9QNApVqGq1h6eMuV4 | 0.008260 |   6.507279 |   6.181916 |  0.325363 |
| KT1W3oiS6s9NgSxhZY1nCsazW2QbwkmjkET1 | 0.000000 |   0.000004 |   0.000004 |  0.000000 |
| KT1WBDsJhoRvsvsRCmJirz9AFhSySSvzTWVd | 0.008540 |   6.727756 |   6.391369 |  0.336387 |
| KT1Wp4tXL6GUtABkikB68fT7SaPQY2UuFkuE | 0.004274 |   3.367080 |   3.198726 |  0.168354 |
| KT1WQWXvRcMjJB1y6mYZytoS5QsFJyFNDCk5 | 0.000009 |   0.006824 |   0.006483 |  0.000341 |
| KT1XiGwpmguFEnZDtBDDGisGxXw6qKJHPjdB | 0.000000 |   0.000006 |   0.000006 |  0.000000 |
| KT1XPMJx2wuCbbzKZx5jJyKqLpPJMHv58wni | 0.000000 |   0.000000 |   0.000000 |  0.000000 |
| KT1XrBAocuiE3C2vvtgt7PFoazrC1KRi9ZF4 | 0.132739 | 104.571581 |  99.343002 |  5.228579 |
| tz1aX2DF3ioDjqDcTVmrxVuqkxhZh1pLtfHU | 0.000017 |   0.013635 |   0.012954 |  0.000681 |
| tz1dfUssfLfTBoYqsWxMu86ycmLUvfF2abng | 0.000003 |   0.002023 |   0.001922 |  0.000101 |
| tz1gxmCTN8BSwuPLghDydtDKTqnAKyD8QTv7 | 0.002620 |   2.064190 |   1.960981 |  0.103209 |
| tz1hbXhPVUX1fC8hN7fALyaUpdoC6EMgqM2h | 0.004060 |   3.198804 |   3.038864 |  0.159940 |
| tz1Lfs9xYtCvj1xe5UCPG8Gv78d3mFAJn4Dx | 0.321793 | 253.508708 | 240.833273 | 12.675435 |
| tz1LRir5SfRcC4LNfagetqzKRMRjGNBTiHNH | 0.088429 |  69.664269 |  66.181056 |  3.483213 |
| tz1Lv6nFvAWMvNRbQF7UcX4jobGLrAhKQLNN | 0.017043 |  13.426673 |  12.755340 |  0.671333 |
| tz1MWAyijbHHqwxA2zD8bjr75wpJJhwViqzW | 0.003123 |   2.460268 |   2.337255 |  0.123013 |
| tz1MXhttCeSJYpF3QRmPkMLCfNZaVufEuJmJ | 0.007493 |   5.902638 |   5.607507 |  0.295131 |
| tz1PB27kbPL64MWYoNZAfQAEmzCZFi9EvgBw | 0.000022 |   0.017191 |   0.016332 |  0.000859 |
| tz1R9vogbJQ4QpEnhFjut6SfyoopP17KkdMc | 0.005425 |   4.273754 |   4.060067 |  0.213687 |
| tz1RoDhaKjJjqcVy9MCN85bVCvbXHEnAFC7j | 0.000748 |   0.589201 |   0.559741 |  0.029460 |
| tz1RomjUZ1j9F2vqE24h2Am8UeGUpcrf6vvJ | 0.000043 |   0.033532 |   0.031856 |  0.001676 |
| tz1SnvfwMUYfD2uJrHBiaj4XPstW3eUE9RJU | 0.007168 |   5.646608 |   5.364278 |  0.282330 |
| tz1Vcu87ZuUK2e8BcoCBUWUhu2s2hPAabStm | 0.000066 |   0.051776 |   0.049188 |  0.002588 |
| tz1VESLfEAEwDEKhyLZJYXVoervFk5ABPUUD | 0.003192 |   2.514661 |   2.388928 |  0.125733 |
| tz1VeiAS5wvYgNdri6vwDUrctQ5XhhaXY3K9 | 0.000086 |   0.067751 |   0.064364 |  0.003387 |
| tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH | 0.000366 |   0.287945 |   0.273548 |  0.014397 |
| tz1Z48RMPT1vjqNyUASCexnCEvEEE93J1pwL | 0.000515 |   0.405388 |   0.385119 |  0.020269 |
+--------------------------------------+----------+------------+------------+-----------+
|                                                     TOTAL    | 664.453233 | 34.971180 |
+--------------------------------------+----------+------------+------------+-----------+
```

### Run
```
➜  tzpay git:(v2) ✗ ./tzpay dryrun 206 --table
+--------------------------------------+--------------------------------------+------------+-------------------------------------------------------+
|                BAKER                 |                WALLET                |  REWARDS   |                     OPERATION                         |
+--------------------------------------+--------------------------------------+------------+-------------------------------------------------------+
| tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc | tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH | 787.800000 | "oofFEAqiRLEXZzGyWj3hW14WemnKqLgsbAT3YH7GnoBUUfQNgs4" |
+--------------------------------------+--------------------------------------+------------+-------------------------------------------------------+
+--------------------------------------+----------+------------+------------+-----------+
|              DELEGATION              |  SHARE   |   GROSS    |    NET     |    FEE    |
+--------------------------------------+----------+------------+------------+-----------+
| KT18kTf8UujihcF46Zn3rsFdEYFL1ZNFnGY4 | 0.001337 |   1.053357 |   1.000690 |  0.052667 |
| KT18ni9Yar4UzwZozFbRF7SFUKg2EqyyUPPT | 0.000040 |   0.031549 |   0.029972 |  0.001577 |
| KT18uqwoNyPRHpHCrg7xBFd7CiAZMbS1Ffne | 0.000460 |   0.362624 |   0.344493 |  0.018131 |
| KT193c72q6eP1VpaY7hiheE7k1eDZiXeQUUw | 0.000036 |   0.028706 |   0.027271 |  0.001435 |
| KT19ABG9KxbEz2GrdN6uhGfxLmMY7REikBN8 | 0.004862 |   3.830464 |   3.638941 |  0.191523 |
| KT19Aro5JcjKH7J7RA6sCRihPiBQzQED3oQC | 0.001330 |   1.048020 |   0.995619 |  0.052401 |
| KT19Q8GiYqGpuuUjf9xfXXVu1WY889N8oxRe | 0.000000 |   0.000053 |   0.000051 |  0.000002 |
| KT1A1sZmBQS9oZnPePRwP3Jyzv41xEppxfbF | 0.002591 |   2.041285 |   1.939221 |  0.102064 |
| KT1A5seo53aLSSyHgJKZFYnh7jTZBtFnNnjz | 0.000000 |   0.000002 |   0.000002 |  0.000000 |
| KT1Aeg9D8kvkbAb6yikUdFcroReXvHtMBaZz | 0.000060 |   0.046964 |   0.044616 |  0.002348 |
| KT1AT7N9bGhViSorUrpivuYT6Wxs37hR2p9d | 0.000002 |   0.001336 |   0.001270 |  0.000066 |
| KT1AThmRzcn51NwMf25NFYTqawjVo62hWiCv | 0.005695 |   4.486706 |   4.262371 |  0.224335 |
| KT1BjtEUxd25wwdwGH432LoP6PskvUc2bEYV | 0.003537 |   2.786213 |   2.646903 |  0.139310 |
| KT1BXmBgMSViAViNyhvkb441e2RBFMiKdnj7 | 0.001270 |   1.000411 |   0.950391 |  0.050020 |
| KT1C28u6DWsBfXk3UMyGrd8zTUVMpsyvjxmp | 0.002488 |   1.959845 |   1.861853 |  0.097992 |
| KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC | 0.029346 |  23.118649 |  21.962717 |  1.155932 |
| KT1CeUNtCrXFNbLmvdGPNnxpcJw2sW5Hcpmc | 0.000000 |   0.000095 |   0.000091 |  0.000004 |
| KT1CQiyDJ3mMVDoEqLY8Fz1onFXo5ycp5BDN | 0.001330 |   1.047706 |   0.995321 |  0.052385 |
| KT1CySPLDUSYyJ9vqNCF2dGgit4Rw2yUNEcj | 0.000159 |   0.125470 |   0.119197 |  0.006273 |
| KT1Cz1jPLuaPR99XamKQDr9PKZY1PTXzTAHH | 0.000000 |   0.000033 |   0.000032 |  0.000001 |
| KT1Dgma8bbDtAbtMbYYS5VmziyCANAZn8M7W | 0.000000 |   0.000083 |   0.000079 |  0.000004 |
| KT1E1MnvNgCDLqnGneStVY9CvmjnyZgcaPaD | 0.000000 |   0.000009 |   0.000009 |  0.000000 |
| KT1EbMbqTUS8XnqGVRsdLZVKLhcT7Zc33jR1 | 0.010737 |   8.458494 |   8.035570 |  0.422924 |
| KT1EidADxWfYeBgK8L1ZTbf7a9zyjKwCFjfH | 0.000000 |   0.000151 |   0.000144 |  0.000007 |
| KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy | 0.037406 |  29.468174 |  27.994766 |  1.473408 |
| KT1GBWviYFdRiNkhwM7LfrKDgHWnpdxpURtx | 0.018462 |  14.544046 |  13.816844 |  0.727202 |
| KT1GcSsQaTtMB2HvUKU9b6WRFUnGpGx9JwGk | 0.000000 |   0.000000 |   0.000000 |  0.000000 |
| KT1HccFB3cn4BR2za9XMuU7Wht64omed2UW8 | 0.010748 |   8.467393 |   8.044024 |  0.423369 |
| KT1J2uk1fYSnZjxkJcUhFDkaRDhjCTRBspqv | 0.000000 |   0.000000 |   0.000000 |  0.000000 |
| KT1J4WFQRV3942phzRrh87WDFWKrNVcDJTP9 | 0.007032 |   5.539827 |   5.262836 |  0.276991 |
| KT1JJcydTkinquNqh6kE5HYgFpD2124qHbZp | 0.042888 |  33.787274 |  32.097911 |  1.689363 |
| KT1JoAP7MfiigepR332u6xJqza9CG52ycYZ9 | 0.000000 |   0.000157 |   0.000150 |  0.000007 |
| KT1JPeGNVarLsPZnSb3hG5xMVmJJmmBnrnpT | 0.003366 |   2.652015 |   2.519415 |  0.132600 |
| KT1JsHBFpoGRVXpcfC763YwvonKtNvaFotpG | 0.000003 |   0.002021 |   0.001920 |  0.000101 |
| KT1Jw925NVi4FzTVohZk5iLqagnhJGDEQoTS | 0.000479 |   0.376989 |   0.358140 |  0.018849 |
| KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt | 0.001236 |   0.974078 |   0.925375 |  0.048703 |
| KT1KeNNxEM4NyfrmF1CG6TLn3nRSmEGhP7Z2 | 0.001643 |   1.294659 |   1.229927 |  0.064732 |
| KT1KJ5Qt18yU9DrqN36tgyLtaSvFSZ5r6YL6 | 0.002900 |   2.284941 |   2.170694 |  0.114247 |
| KT1LfoE9EbpczdzUzowRckGUfikGcd5PyVKg | 0.000000 |   0.000031 |   0.000030 |  0.000001 |
| KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv | 0.035778 |  28.185759 |  26.776472 |  1.409287 |
| KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3 | 0.000023 |   0.017797 |   0.016908 |  0.000889 |
| KT1Lm4ZSyXSHod7U6znR7z9SGVmexntNQwAp | 0.000013 |   0.010614 |   0.010084 |  0.000530 |
| KT1Lnh39om2iqr4qb9AarF9T38ayNBLnfAVn | 0.000000 |   0.000034 |   0.000033 |  0.000001 |
| KT1MfT8XvQp9ZeGUx4cmCNF3wui55WLNYhq9 | 0.000000 |   0.000002 |   0.000002 |  0.000000 |
| KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj | 0.024467 |  19.274821 |  18.311080 |  0.963741 |
| KT1MSFeAGaWk8w7F1gmgUMaarU7mH385ueYC | 0.005919 |   4.662630 |   4.429499 |  0.233131 |
| KT1MX2TwjSBzPaSsBUeW2k9DKehpiuMGfFcL | 0.001302 |   1.025690 |   0.974406 |  0.051284 |
| KT1Na4maJ99GE6CGA1vEocWXrKRmxmsVUaTi | 0.001329 |   1.047182 |   0.994823 |  0.052359 |
| KT1NfMCxyzwev243rKk3Y6SN8GfmdLKwASFQ | 0.000039 |   0.030682 |   0.029148 |  0.001534 |
| KT1NGd6RaRtmvwexYXGibtdvKBnNjjpBNknn | 0.000009 |   0.007249 |   0.006887 |  0.000362 |
| KT1NmVtU3CNqzhNWwLhE5BqAopjkcmHpWzT2 | 0.000028 |   0.022144 |   0.021037 |  0.001107 |
| KT1NxnFWHW7bUxzks1oHVU2jn4heu48KC3eD | 0.000000 |   0.000004 |   0.000004 |  0.000000 |
| KT1PDBuQmFLVHfiWZjV248QdTrdcmAuSS7Tx | 0.000000 |   0.000028 |   0.000027 |  0.000001 |
| KT1PY2MMiTUkZQv7CPekXy186N1qmu7GikcT | 0.000000 |   0.000010 |   0.000010 |  0.000000 |
| KT1QB9UAT1okYfcPQLi4jBmZkYg7LHcepERV | 0.001330 |   1.047444 |   0.995072 |  0.052372 |
| KT1QLo7DzPZnYK2EhmWpejVUnFjQUuWFKHnc | 0.001329 |   1.047287 |   0.994923 |  0.052364 |
| KT1RbwPHzDwU9oPjnTWZrbCrMGjaFyj8dEtC | 0.000001 |   0.000711 |   0.000676 |  0.000035 |
| KT1REp3D8dkiVVi37TCSMJNgGeX6UigBtfaL | 0.000001 |   0.001174 |   0.001116 |  0.000058 |
| KT1Re5utTU2hrujXgZ3Ux5BgjN8rbru4sns2 | 0.000574 |   0.452440 |   0.429818 |  0.022622 |
| KT1RuTPgQ6kdpnE3Adnw7Hr2KFN45uC3BdBy | 0.003767 |   2.967317 |   2.818952 |  0.148365 |
| KT1S9VbEnU8nj33ufxrGBYGxBCnqmeoAnKt4 | 0.000000 |   0.000048 |   0.000046 |  0.000002 |
| KT1T3dPMBm7D3kKqALKYnW2mViFqMMVCYtmo | 0.000028 |   0.022193 |   0.021084 |  0.001109 |
| KT1TDrRrdz6SLYLBw8ZDxLWwJpx7FVpC52bt | 0.001734 |   1.365962 |   1.297664 |  0.068298 |
| KT1TS49jiXxrnwhoJzAvCzGZCXLJs3XV1k6C | 0.001340 |   1.055731 |   1.002945 |  0.052786 |
| KT1Uh1G9tdq45N63ZBrreDKy7eZF8QVoydm1 | 0.000000 |   0.000024 |   0.000023 |  0.000001 |
| KT1UVUasDXH6mg8NCzRRgqvcjMoDUpETYEzH | 0.001329 |   1.047182 |   0.994823 |  0.052359 |
| KT1VUbpty8fER7npuvsfYDZXf2wVPhAHVqSx | 0.000004 |   0.002879 |   0.002736 |  0.000143 |
| KT1VzTs5piA7kYQkkfA9QNApVqGq1h6eMuV4 | 0.008260 |   6.507279 |   6.181916 |  0.325363 |
| KT1W3oiS6s9NgSxhZY1nCsazW2QbwkmjkET1 | 0.000000 |   0.000004 |   0.000004 |  0.000000 |
| KT1WBDsJhoRvsvsRCmJirz9AFhSySSvzTWVd | 0.008540 |   6.727756 |   6.391369 |  0.336387 |
| KT1Wp4tXL6GUtABkikB68fT7SaPQY2UuFkuE | 0.004274 |   3.367080 |   3.198726 |  0.168354 |
| KT1WQWXvRcMjJB1y6mYZytoS5QsFJyFNDCk5 | 0.000009 |   0.006824 |   0.006483 |  0.000341 |
| KT1XiGwpmguFEnZDtBDDGisGxXw6qKJHPjdB | 0.000000 |   0.000006 |   0.000006 |  0.000000 |
| KT1XPMJx2wuCbbzKZx5jJyKqLpPJMHv58wni | 0.000000 |   0.000000 |   0.000000 |  0.000000 |
| KT1XrBAocuiE3C2vvtgt7PFoazrC1KRi9ZF4 | 0.132739 | 104.571581 |  99.343002 |  5.228579 |
| tz1aX2DF3ioDjqDcTVmrxVuqkxhZh1pLtfHU | 0.000017 |   0.013635 |   0.012954 |  0.000681 |
| tz1dfUssfLfTBoYqsWxMu86ycmLUvfF2abng | 0.000003 |   0.002023 |   0.001922 |  0.000101 |
| tz1gxmCTN8BSwuPLghDydtDKTqnAKyD8QTv7 | 0.002620 |   2.064190 |   1.960981 |  0.103209 |
| tz1hbXhPVUX1fC8hN7fALyaUpdoC6EMgqM2h | 0.004060 |   3.198804 |   3.038864 |  0.159940 |
| tz1Lfs9xYtCvj1xe5UCPG8Gv78d3mFAJn4Dx | 0.321793 | 253.508708 | 240.833273 | 12.675435 |
| tz1LRir5SfRcC4LNfagetqzKRMRjGNBTiHNH | 0.088429 |  69.664269 |  66.181056 |  3.483213 |
| tz1Lv6nFvAWMvNRbQF7UcX4jobGLrAhKQLNN | 0.017043 |  13.426673 |  12.755340 |  0.671333 |
| tz1MWAyijbHHqwxA2zD8bjr75wpJJhwViqzW | 0.003123 |   2.460268 |   2.337255 |  0.123013 |
| tz1MXhttCeSJYpF3QRmPkMLCfNZaVufEuJmJ | 0.007493 |   5.902638 |   5.607507 |  0.295131 |
| tz1PB27kbPL64MWYoNZAfQAEmzCZFi9EvgBw | 0.000022 |   0.017191 |   0.016332 |  0.000859 |
| tz1R9vogbJQ4QpEnhFjut6SfyoopP17KkdMc | 0.005425 |   4.273754 |   4.060067 |  0.213687 |
| tz1RoDhaKjJjqcVy9MCN85bVCvbXHEnAFC7j | 0.000748 |   0.589201 |   0.559741 |  0.029460 |
| tz1RomjUZ1j9F2vqE24h2Am8UeGUpcrf6vvJ | 0.000043 |   0.033532 |   0.031856 |  0.001676 |
| tz1SnvfwMUYfD2uJrHBiaj4XPstW3eUE9RJU | 0.007168 |   5.646608 |   5.364278 |  0.282330 |
| tz1Vcu87ZuUK2e8BcoCBUWUhu2s2hPAabStm | 0.000066 |   0.051776 |   0.049188 |  0.002588 |
| tz1VESLfEAEwDEKhyLZJYXVoervFk5ABPUUD | 0.003192 |   2.514661 |   2.388928 |  0.125733 |
| tz1VeiAS5wvYgNdri6vwDUrctQ5XhhaXY3K9 | 0.000086 |   0.067751 |   0.064364 |  0.003387 |
| tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH | 0.000366 |   0.287945 |   0.273548 |  0.014397 |
| tz1Z48RMPT1vjqNyUASCexnCEvEEE93J1pwL | 0.000515 |   0.405388 |   0.385119 |  0.020269 |
+--------------------------------------+----------+------------+------------+-----------+
|                                                     TOTAL    | 664.453233 | 34.971180 |
+--------------------------------------+----------+------------+------------+-----------+
```

## Roadmap:
* inlucde fiat price of XTZ in reports
* tax reporting

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details