# Platform Service's Rotating Shop Items Plugin gRPC Demo App

A CLI demo app to prepare required data and execute Rotating Shop Items Plugin gRPC for AGS's Platform Service.
Following diagram will explain how this CLI demo app works.
```mermaid
sequenceDiagram
    participant A as CLI Demo App
    participant I as AGS IAM
    participant P as AGS Platform
    participant G as Grpc Plugin Server
    
    A ->> I: user login
    I -->> A: auth token
    A ->> P: Set gRpc server target
    A ->> P: Create draft store and publish it
    A ->> P: Create item's category
    A ->> P: Create store's view
    A ->> P: Create items
    A ->> P: Create a section and put items in it
    A ->> P: Enable custom rotation for newly created section
    A ->> P: Get active sections
    P ->> G: Call rotating item function and/or backfill function
    G -->> P: Returns rotated item(s)
    P -->> A: Return active sections with rotated item(s)
    A ->> P: Delete store
    A ->> P: Unset gRpc server target
```

## Prerequsites

* Go 1.18 or latest

## Build

To build this CLI sample app, execute the following command in the working CLI directory.

```bash
$ go build .
```
and install the binary file
```bash
$ go install accelbyte.net/rotating-shop-items-cli
```
then run the command usage

## Usage

### Setup

The following environment variables are used by this CLI demo app.
```
export AB_BASE_URL='https://test.accelbyte.io'
export AB_CLIENT_ID='xxxxxxxxxx'
export AB_CLIENT_SECRET='xxxxxxxxxx'

export AB_NAMESPACE='namespace'
export AB_USERNAME='USERNAME'
export AB_PASSWORD='PASSWORD'
```
If these variables aren't provided, you'll need to supply the required values via command line arguments.

Also, you will need `Rotating Shop Items Plugin gRPC` server already deployed and accessible. If you want to use your local development environment, you can use tunneling service like `ngrok` to tunnel your grpc server port so it can be accessed by AGS.
> Current AGS deployment does not support mTLS and authorization for custom grpc plugin. Make sure you disable mTls and authorization in your deployed Grpc server.

### Example
CLI demo app requires 2 parameters. first is the grpc server url, and second is run mode. For run mode, only two availables: `rotation` and `backfill`. `rotation` mode will call custom rotation item function in grpc server, and `backfill` mode will call custom backfill function.

- With all environment variables setup
```bash
$ rotating-shop-items-cli rotatingShop -n accelbyte -b http://localhost:8081 -i clientId -s clientSecret -u username123 -p password123 -a /customitemtestexample -g localhost:6565 -m backfill
```
- Show usage help
```bash
 $ rotating-shop-items-cli rotatingShop -h
```