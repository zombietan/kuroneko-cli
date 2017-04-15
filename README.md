# kuroneko

## installation
`$ go get github.com/zombietan/kuroneko/cmd`  

`$ cd $GOPATH/src/github.com/zombietan/kuroneko`  

`$ go install`  

## usage
`$ kuroneko 伝票番号`  
###### example
```sh:example
伝票番号 xxxx-xxxx-xxxx
配達完了
このお品物はお届けが済んでおります。

荷物受付　　　　　　　　　　　　　　　| 03/14 | 12:16 | ＸＸセンター　　　　　　　　　　　　　　| 012345 |
発送　　　　　　　　　　　　　　　　　| 03/14 | 12:16 | ＸＸセンター　　　　　　　　　　　　　　| 678901 |
持戻（ご不在）　　　　　　　　　　　　| 03/15 | 14:46 | ＸＸＹセンター　　　　　　　　　　　　　| 234567 |
配達完了　　　　　　　　　　　　　　　| 03/16 | 19:10 | ＹＸＸセンター　　　　　　　　　　　　　| 890123 |
---------------------------------------------------------------------------------------------------
```
