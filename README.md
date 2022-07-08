## build

```
go run *.go
```

## tests

```
go test ./... -v -short
```

## Custom source URI

Specify your own URI using the `CANDY_URI` env var
```
CANDY_URI=https://candystore.zimpler.net/ go run *.go
```