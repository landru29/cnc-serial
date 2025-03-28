# CNC Serial

Simple serial monitor to communicate with UART and control 3D printers or CNC. This is running on a text console.

It implements :
* helpers for G-code
* program management (read, send commands step by step)

![screen.png](screen.png "example")



## Prerequisite

You must have a sane instalation of :
* golang (minimum version: 1.18)
* protobuf-compiler
* protoc-gen-go

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Build

```
make build
```

## Usage

```
CNC monitor

Usage:
  cnc [command]

Available Commands:
  agent       manage the local agent
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  mock        CNC mock monitor
  rpc         CNC RPC monitor
  serial      CNC Serial monitor

Flags:
      --grpc        RPC connection
  -h, --help        help for cnc
  -l, --lang lang   language (available: en, fr) (default en)
```

```mermaid
graph TD;

CNC-->mock
CNC-->serial
CNC-->rpc
CNC-->agent
agent-->rpcAgent[rpc]
rpcAgent-->rpcAgentSerial[serial]
rpcAgent-->rpcAgentMock[mock]

```


`filename`, if specified, must be a valide G-Code.

## Debug

You must install Delve

```
go install github.com/go-delve/delve/cmd/dlv@latest
```

## Architecture

```mermaid
graph TD;
    GCodeProcessor(**GCodeProcessor**
    process G-Code);
    Transporter(**Transporter**
    transport command to CNC);
    Stacker(**Stacker**
    remember previous commands
    );
    Controller(**Controller**
    orchestrate processus
    );
    Programmer(**Programmer**
    Manage the execution of the program
    );
    Screen(**Screen**
    Layout display
    );

    serial([serial]);
    nop([nop]);
    rpc([RPC
    client]);

    Controller-->Stacker;
    Controller-->GCodeProcessor;
    Controller-->Programmer;
    Controller-->Transporter;
    Transporter-->serial;
    Transporter-->nop;
    Transporter-->rpc;
    Screen-->Controller;
    Controller-->Screen;
    Screen-->Stacker;

    rpcServer(RPC
    server)-->Transporter
```
