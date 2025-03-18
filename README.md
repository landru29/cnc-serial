# Serial

Simple serial monitor to communicate with UART. It implements helpers for G-code to control 3D printers or CNC

## Prerequisite

You must have a sane instalation of golang (minimum version: 1.18)

## Build

```
make build
```

## Usage

```
Serial monitor

Usage:
  serial [flags]

Flags:
  -b, --bit-rate int   Bit rate (default 115200)
  -d, --dry-run        Dry run (do not open serial port)
  -h, --help           help for serial
  -l, --lang lang      language (available: en, fr) (default en)
  -p, --port string    Port name
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
    Screen(**Screen**
    Layout display
    );
    Application(**Application**
    main entrypoint);

    Application-->Controller;
    Application-->GCodeProcessor;
    Application-->Screen;
    Controller-->Stacker;
    Controller-->GCodeProcessor;
    Controller-->Transporter;
    Transporter-->serial([serial]);
    Transporter-->nop([nop]);
    Screen-->Controller;
    Controller-->Screen;
    Screen-->Stacker;
```
