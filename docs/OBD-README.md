# OmniBOLT Daemon | Smart Asset Lightning Network
[![](https://img.shields.io/badge/license-MIT-blue)](https://github.com/omnilaboratory/obd/blob/master/LICENSE) [![](https://img.shields.io/badge/standard%20readme-OK-brightgreen)](https://github.com/omnilaboratory/obd/blob/master/README.md) [![](https://img.shields.io/badge/golang-%3E%3D1.9.0-orange)](https://golang.org/dl/) [![](https://img.shields.io/badge/protocol-OmniBOLT-brightgreen)](https://github.com/omnilaboratory/OmniBOLT-spec) 
[![](https://img.shields.io/badge/API%20V0.3-Document-blue)](https://api.omnilab.online) 

OBD implements the [OmniBOLT](https://github.com/omnilaboratory/OmniBOLT-spec) specification, and it is an open source, off-chain decentralized platform, build upon BTC/OmniLayer network, implements basic HTLC payment, multi-currency atomic swap, and more off-chain contracts on the network of [smart assets enabled lightning channels](https://github.com/omnilaboratory/OmniBOLT-spec/blob/master/OmniBOLT-02-peer-protocol.md#omni-address).  

Clone, compile the source code and run the binary executable file, you will have a featured OmniBOLT deamon(OBD) to start the journey of lightning network.   


# Table of Contents
 * [Background](#background)
 * [Dependency](#dependency)
 * [Installation](#installation and minimum system requirement)
	* [Step 1: fetch the source code](#step-1-fetch-the-source-code)
	<!--* [Step 2: set up OmniCore node](#step-2)
		* [option 1: remote omnicore node](#option-1-remote-omnicore-node)	
		* [option 2: local omnicore node](#option-2-local-omnicore-node) --> 
	* [Step 2: Connect to a tracker](#step-2-connect-to-a-tracker)
	* [Step 3: compile and run OmniBOLT daemon](#step-3-compile-and-run-omnibolt-daemon)
	* [Step 4: test channel operations using GUI testing tool](#step-4-test-channel-operations-using-gui-testing-tool)
	* [Step 5: channel operations on test site](#step-5-channel-operations-on-test-site)
		* [sign up](#sign-up)	
		* [login](#login)
		* [normal operations](#normal-operations)

	<!-- Removed by Neo Carmack 2020-06-09 -->		
	<!-- 	* [create channel](https://github.com/omnilaboratory/obd#create-channel)
		* [deposit](https://github.com/omnilaboratory/obd#deposit)
		* [payments in channel](https://github.com/omnilaboratory/obd#payments-in-a-channel)
		* close channel (TBD) -->

	<!-- Added by Kevin Zhang 2019-11-19 -->
	<!-- Removed by Neo Carmack 2020-06-09 -->
	<!-- 	* [Step 6: transfer assets through HTLC](https://github.com/omnilaboratory/obd#step-6-transfer-assets-through-HTLC) -->

 * [API Document](#api-document)
 * [How to Contribute](#how-to-contribute)
 * [Current Features](#current-features)
 * [Comming Features](#comming-features)
 * [Experimental Features](#experimental-features)
 * [Related Projects](#related-projects)

# Background

Blockchain industry requires a much more flexible, extensible, and cheaper smart assets circulation solution to solve the main chain scalability problem. Lightning network is a solid technology for this problem.  

We propose [OmniBOLT](https://github.com/omnilaboratory/OmniBOLT-spec) to enable lightning network to be smart asset aware. OBD is the golang implementation. Interested readers please go to the [spec repository](https://github.com/omnilaboratory/OmniBOLT-spec) for further understanding of its advantages and how it works.  


# Dependency

If you run your own tracker, you should maintain an [Omnicore 0.18](https://github.com/OmniLayer/omnicore)(or later) full node, which integrates the latest BTC core 0.18 and which enables relative time locker used in RSM contracts and HTL contracts.

Running an obd node doesn't require a full BTC/Omni node. The obd node should connect a remote tracker for full node services. 
 

# Installation and minimum system requirement
The following instruction works for Ubuntu 14.04.4 LTS, golang 1.10 or later.

You could run obd on your desktop or a cloud server you trust. You could download the binary we compiled to major platforms at [here](https://github.com/omnilaboratory/obd/releases), and ignore steps 1 and 3. 

The following are the estimated system requirements for running an obd node:  
* 2.0 GHz 64-bit processor  
* 4 GB memory  
* 100 GB HDD(SSD would be better)  
* Ubuntu 14.04.4 LTS or later  
* golang 1.10 or later  


For running a tracker, the system requirement is higher, because a tracker requires a full BTC/omnicore node. We suggest you have:  
* 4 .0 GHz 64-bit processor  
* 16 GB memory  
* 500 GB HDD(SSD would be better)  for a btc/omnicore full node  
* Ubuntu 14.04.4 LTS or later  
* golang 1.10 or later  



## step 1: fetch the source code

on your terminal:

```
$ git clone https://github.com/omnilaboratory/obd.git
```

or if you already set up your local git repo, you just need to fetch the latest version: 

```
$ git pull origin master
```

check if all updated:

```
$ git remote -v
origin	https://github.com/omnilaboratory/obd.git (fetch)
origin	https://github.com/omnilaboratory/obd.git (push)
```

<!-- 
## Step 2: 
#### option 1: Remote Omnicore node 
Use our remote OmniCore node. Go to `\config\conf.ini`, you will see:
```
[chainNode]
host=62.234.216.108:18332
user=omniwallet
pass=cB3]iL2@eZ1?cB2?
```
This is a testing full node for our community to run/call/test omni commands remotely. The OmniBOLT daemon invokes Omni RPC commands from this node if you use this configuration. It is a convenient way to get started.  


#### option 2: Local Omnicore node 
[Install OmniCore](https://github.com/OmniLayer/omnicore#installation) on your local machine. Omnicore requires a full BTC core node, which may take days to synchronize the whole BTC database to your local device. After finishing synchronization, you can run omni/BTC commands for experiments, such as constructing raw transactions or generating new addresses.

Edit the configure file: `\config\conf.ini`
```
[chainNode]
host=127.0.0.1:port
user=your user name
pass=your password
```
-->
## Step 2: Connect to a tracker

Trackers offer such anonymous services: monitor node service quality, record channel balance if the channel is not private, update routing table for connected nodes, broadcast transactions, and help nodes to find paths for payments.

Anyone can be a tracker maintainer. Running a tracker requires a full omnilayer core node, which may take days to synchronize the whole database to your device. So that we deployed several trackers for our community. 

Edit the configure file: `\config\conf.ini`
```
[tracker]
host = 62.234.216.108:60060
```


## Step 3: Compile and Run OmniBOLT Daemon

Wait till all data downloaded.

```
$ go build obdserver.go
```
which generates the executable binary file `obdserver` under the source code directory. 

if you want to generate exe file for windows platform, use this:
```
$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build obdserver.go
```
you will see an obdserver.exe file generated under the same directory.

### Startup  
Run:
```
$ ./obdserver 
```

By default, the configuration file `./config/conf.ini` will be loaded. 
To override this set the optional `--configPath <path>` program argument:
```
$ ./obdserver --configPath "/config/conf.ini"  
```

The terminal displays:
```
2019/08/23 23:05:15 rpcclient.go:23: &{62.234.216.108:18332 omniwallet cB3]iL2@eZ1?cB2?}
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ws                       --> LightningOnOmni/routers.wsClientConnect (3 handlers)
```
Which tells us the daemon is running. We are going to use WebSocket testing tools to test our obd commands.

### Running obd in the exclusive mode

Once obd is launched through the steps mentioned above, you can then invoke the [Login](https://api.omnilab.online/?shell#login) gRPC API to run obd in the exclusive mode.

You will get the `login_token` (known as password) indicated by the red arrow in the screenshot of a running obd. The `login_token` is used to login to obd.

<p align="center">
  <img width="750" alt="An example of the launch an obd" src="assets/launch-an-obd.png">
</p>

And while in the exclusive mode, there are gRPC APIs used to interact with obd.

* Connect remote peer using [ConnectPeer](https://api.omnilab.online/?shell#connectpeer).
* Create a new channel using [OpenChannel](https://api.omnilab.online/?shell#openchannel)
* Fund to channel using [FundChannel](https://api.omnilab.online/?shell#fundchannel)
* Payment with RSMC using [RsmcPayment](https://api.omnilab.online/?shell#rsmcpayment) 
* Create a new invoice using [AddInvoice](https://api.omnilab.online/?shell#addinvoice)
* Payment with HTLC using [SendPayment](https://api.omnilab.online/?shell#sendpayment)

For more details on how to use these APIs, please refer to the online documentation at [API Website](https://api.omnilab.online/?shell#obd-grpc-api-reference).


## Step 4: Test channel operations using GUI testing tool.

**NOTE: You should replace all of relevant data by the exact value that your own OBD generates for you**

OmniBOLT deamon(OBD) exposes WebSocket services for client interaction. For ease of use, we released GUI tool to help users to get started. Go to the [GUI tool repository](https://github.com/omnilaboratory/DebuggingTool) to download and try it.

<p align="center">
  <img width="500" alt="Debugging Tool Screenshot" src="assets/image_screen.png">
</p>

If you don't want to deploy obd to start with, there is a list of testnet nodes in that GUI tool repo that you can connect to.  

Another option is to use web socket test client for Chrome to do experiments. Install it from:
```
https://chrome.google.com/webstore/detail/websocket-test-client/fgponpodhbmadfljofbimhhlengambbn?hl=en
```
Make sure your browser supports WebSocket, as displayed in this screenshot.

<p align="center">
  <img width="500" alt="Screenshot of Websocket online testing site" src="assets/WebSocketTestSite.png">
</p>

Input `ws://127.0.0.1:60020/wstest`, press `Open`. If in the right text box, displays `OPENED`, then we are ready to send messeages to obd.

The first message is to sign up as `Alice`. input the following request into the Request box, and press `SEND`:

```json
{
	"type":-102004
}
```

In the `Message Log` pannel, displays the response message from OBD:

*Return mnemonic words by hirarchecal deterministic wallet system.*

```json
{
    "type":-102004,
    "status":true,
    "from":"c2215a60-8b81-439f-8cb3-11ba51691076",
    "to":"c2215a60-8b81-439f-8cb3-11ba51691076",
    "result":"two ribbon knee leaf easy pottery hobby pony mule test bridge liar sand mirror decline gasp focus this park undo rough cricket portion ignore"
}
```

Then go to login as `Alice`. input the following message and press `SEND`:

*The mnemonic words is as a login name.*

```json
{
    "type":-102001,
    "data":{
        "mnemonic":"two ribbon knee leaf easy pottery hobby pony mule test bridge liar sand mirror decline gasp focus this park undo rough cricket portion ignore"
    }
}
```

In the `Message Log` pannel, displays the response message from OBD:

*A SHA256 string of mnemonic words as a user id.*

```json
{
    "type":-102001,
    "status":true,
    "from":"7da8d2441e0ad67040a274902f1965ee1a5c3fdd86f1ddc3280eda5230e006f2",
    "to":"all",
    "result":"7da8d2441e0ad67040a274902f1965ee1a5c3fdd86f1ddc3280eda5230e006f2 login"
}
```

It works.

## Step 5: Channel Operations on test site

For the convenience of brand new users, we suggest to connect our testnet nodes(for testing only). The URL is:

```
ws://62.234.216.108:60020/wstest
```
Open two chrom browsers, left is Alice and the right is Bob. Input URL and click `OPEN`, then both status will show `OPENED`.


### Sign up

1、Alice sign up

Websocket request:

```
{
	"type":-102004
}
```

OBD responses:

```json
{
    "type":-102004,
    "status":true,
    "from":"c2215a60-8b81-439f-8cb3-11ba51691076",
    "to":"c2215a60-8b81-439f-8cb3-11ba51691076",
    "result":"two ribbon knee leaf easy pottery hobby pony mule test bridge liar sand mirror decline gasp focus this park undo rough cricket portion ignore"
}
```

2、Bob sign up

Websocket request:

```
{
	"type":-102004
}
```

OBD responses:

```json
{
    "type":-102004,
    "status":true,
    "from":"cec4e1db-ef38-4508-a9bf-8c5976df1916",
    "to":"cec4e1db-ef38-4508-a9bf-8c5976df1916",
    "result":"outer exhibit burger screen onion dog ensure net depth scan steel field pizza group veteran doctor rhythm inch dawn rotate gravity index modify utility"
}
```

### Login

1、Alice login

Websocket request:

```json
{
	"type":-102001,
    "data":{
        "mnemonic":"two ribbon knee leaf easy pottery hobby pony mule test bridge liar sand mirror decline gasp focus this park undo rough cricket portion ignore"
    }
}
```

OBD responses:

```json
{
    "type":-102001,
    "status":true,
    "from":"7da8d2441e0ad67040a274902f1965ee1a5c3fdd86f1ddc3280eda5230e006f2",
    "to":"all",
    "result":"7da8d2441e0ad67040a274902f1965ee1a5c3fdd86f1ddc3280eda5230e006f2 login"
}
```

2、Bob login

Websocket request:

```json
{
    "type":-102001,
    "data":{
        "mnemonic":"outer exhibit burger screen onion dog ensure net depth scan steel field pizza group veteran doctor rhythm inch dawn rotate gravity index modify utility"
    }
}
```

OBD responses:

```json
{
    "type":-102001,
    "status":true,
    "from":"f38e72f6bf69c69ad1cdc0040550bafb86d5c4d35bd04542fcf5fc5ecb2135be",
    "to":"all",
    "result":"f38e72f6bf69c69ad1cdc0040550bafb86d5c4d35bd04542fcf5fc5ecb2135be login"
}
```

*A SHA256 string of mnemonic words as a user id.*

Alice's id is: 7da8d2441e0ad67040a274902f1965ee1a5c3fdd86f1ddc3280eda5230e006f2

Bob's   id is: f38e72f6bf69c69ad1cdc0040550bafb86d5c4d35bd04542fcf5fc5ecb2135be


### More channel operations
Instruction can be found in [GUI tool](https://omnilaboratory.github.io/obd/#/GUI-tool) or the [online API document.](https://api.omnilab.online):


# API Document
Please visit OBD [online API documentation](https://api.omnilab.online) for the lastest update.


# How to Contribute
OmniBOLT Daemon is MIT licensed open source software. Hopefully, you can get started by going through the above steps, but the lightning network is not that easy to develop. Anyone is welcome to join us in this journey, and please be nice to each other, don't bring any illegal/private stuff, abuse or racial into our community.

Please submit issues to this repo or help us with those open ones.

Guidelines:

  * read the [OmniBOLT](https://github.com/omnilaboratory/OmniBOLT-spec) spec. If you have any questions over there, raise issues in that repo.
  * ask questions or talk about things in Issues.
  * make branches and raise pull-request, even if working on the main repository.
  * don't copy/paste any code from anywhere else in contribution, because we have limited resources to compare source codes to avoid legal issues. What we can do is read your code, run tests of your newly developed modules and read your comments in your branch to see if it is solving a real problem. 
  * better running `go fmt` before committing any code.
  * add test to any package you commit.
  * write/contribute light client testing tools, such as an HTML page supporting WebSocket, so that new programmers can have an intuitive experience to get started. That helps. We will release our tools for testing.


Join us in [OmniBOLT slack channel](https://join.slack.com/t/omnibolt/shared_invite/enQtNzY2MDIzNzY0MzU5LTFlZTNlZjJhMzQxZTU2M2NhYmFjYjc1ZGZmODYwMWE3YmM0YjNhZWQyMDU2Y2VlMWIxYWFjN2YwMjlmYjUxNzA)
 

# Related projects: 

[https://github.com/OmniLayer/omniwallet](https://github.com/OmniLayer/omniwallet)

[https://github.com/OmniLayer/omnicore](https://github.com/OmniLayer/omnicore)

[https://github.com/OmniLayer/OmniJ](https://github.com/OmniLayer/OmniJ)

[https://github.com/OmniLayer/spec](https://github.com/OmniLayer/spec)

[https://github.com/omnilaboratory/OmniBOLT-spec](https://github.com/omnilaboratory/OmniBOLT-spec)

[https://github.com/lightningnetwork/lightning-rfc](https://github.com/lightningnetwork/lightning-rfc)

[https://github.com/lightningnetwork/lnd](https://github.com/lightningnetwork/lnd)

[https://github.com/omnilaboratory/OmniWalletMobile](https://github.com/omnilaboratory/OmniWalletMobile)





 


