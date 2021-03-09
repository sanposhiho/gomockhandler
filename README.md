# gomockhandler

**This tool is currently under development.
If you find any bugs or have feature requests, please feel free to create an issue.**

gomockhandler is handler of [golang/mock](https://github.com/golang/mock), as the name implies.

Some of you may often create mock with `go generate` like below.
```
//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAG
```

But, it will take time to execute `go generate ./...` for projects with many files. And we cannot easily check if mock is up-to-date.

`gomockhandler` use one config file to generate all mocks.

With `gomockhandler`, 

- You can generate mocks more **quickly** :rocket:.
- You can check if mock is up-to-date :sparkles:.
- You can manage your mocks in one config file :books:.
- You can generate/edit the config of gomockhandler with CLI :wrench:.

Here is some example of the mock being generated in half the time with `gomockhandler`. (I ran `mockgen` to generate same mocks in `go generate ./...`)


<img width="825" alt="Screen Shot 2021-03-08 at 23 28 46" src="https://user-images.githubusercontent.com/44139130/110334403-1444ba00-8066-11eb-9377-0d8c98a84c9e.png">


![Screen Shot 2021-03-09 at 12 07 03](https://user-images.githubusercontent.com/44139130/110412721-fc098500-80cf-11eb-9de6-c98749b0a3fb.png)



## Install

You have to install `mockgen` first.

### Go version < 1.16
```
GO111MODULE=on go get github.com/golang/mock/mockgen
GO111MODULE=on go get github.com/sanposhiho/gomockhandler
```
### Go 1.16+
```
go install github.com/golang/mock/mockgen
go install github.com/sanposhiho/gomockhandler
```

## How to use

### [preparation] generate the config file

See [# Add a new mock to be generated](#add-a-new-mock-to-be-generated).

If you use `go:generate` to execute mockgen, you can generate the config file by rewriting `go:generate` comment a little bit.

replace from `mockgen` to `gomockhandler -project_root=/path/to/project_root`, and run `go generate ./...` in your project.

```
- //go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAG
+ //go:generate gomockhandler -project_root=/path/to/project -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAG
```

gomockhandler will generate a config file named `gomockhandler.json` in your project root directory.

### generate mock

```
gomockhandler -config=gomockhandler.json -concurrency=100 mockgen
```

### check if mock is up-to-date

It is useful for ci to check if all mocks are up-to-date

```
gomockhandler -config=gomockhandler.json check
```

You can see the error if some mocks are not up to date.

```
2021/03/06 02:37:16 mock is not up to date. source: user.go, destination: ../mock/user.go
```

## How to manage your mocks

### Add a new mock to be generated

You can add a new mock to be generated from CLI. You can use the same options as mockgen to add a new mock. If the config file does not exist in your project root directory, it will be created.


`mockgen` has two modes of operation: source and reflect, and, gomockhandler support both.

See [golang/mock#running-mockgen](https://github.com/golang/mock#running-mockgen) for more information.


Example(Source mode):
```
gomockhandler -project_root=/path/to/project -source=foo.go [other options]
```

Example(Reflect mode):
```
gomockhandler -project_root=/path/to/project [options] database/sql/driver Conn,Driver
```

### [WIP]Edit/Delete mocks

Currently, if you want to modify/delete the mock, you have to modify the config manually...

I'm working on developing it to be able to edit/delete it from the CLI.
