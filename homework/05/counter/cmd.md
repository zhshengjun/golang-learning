
```
solc --abi --bin --overwrite \
  -o build \
  script/my_counter.sol

abigen \
  --abi build/MyCounter.abi \
  --bin build/MyCounter.bin \
  --pkg script \
  --type MyCounter \
  --out script/my_counter.go

gofmt -w script/my_counter.go
```