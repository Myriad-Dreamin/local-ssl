# local-ssl

Create ssl certificate in one command.

## Step1: Create Local Certificate Authority (Local CA) Config File

create file `${ca}/ssl.config.json`:

```json
{
  "country": "\"CN\"",
  "organization": "\"Orignization\"",
  "state": "\"State\"",
  "locality": "\"City\"",
  "caOrganizationUnit": "\"Certificate\"",
  "caCommonName": "\"cert.ssl.com\"",
  "emailAddress": "no-reply@ssl.com"
}
```

## Step2: init Local CA Directory

```bash
local-ssl init -project=${ca}
```

## Step3: make certificate for your sites

```bash
local-ssl create -project=data -site=k8s.master01.io -unit=K8sMaster
```

## Command usage

```bash
local-ssl help                                                           
Usage: local-ssl [command]
available commands: help init create
Help Usage:
  local-ssl help [command] - get usage of `command`
```

## Build from source

```bash
go build -o local-ssl ./cmd
```
