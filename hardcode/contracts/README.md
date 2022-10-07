# Go code generation

Download `.abi` text files to your local directory,
then use the provided shell command template to generate Go code.

```sh
abigen \
    --abi=1inchlimitorder.abi \
    --pkg=_1inchlimitorder \
    --type=oneInchLimitOrder \
    --out=1inchlimitorder.go
```
