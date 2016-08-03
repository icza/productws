# productws

[![GoDoc](https://godoc.org/github.com/icza/productws?status.svg)](https://godoc.org/github.com/icza/productws)

This project contains a [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) /
[JSON](https://en.wikipedia.org/wiki/JSON) web service demo in Go with an API to manage products.
The following operations are supported:

- `POST /create` Create a new product
- `GET /list` Get a list of all products
- `GET /details/<id>` Get details about a product
- `PUT /update` Update a product
- `PUT /setprices` Set price points for different currencies for a product

A product has the following attributes:

- `Product.ID` Product ID
- `Product.Name` Name
- `Product.Desc` Description: 
- `Product.Tags` Tags (optional) 
- `Product.Prices` One or more price points (at most one per currency, USD being default)

## Install

You can get it with:

	go get github.com/icza/productws

Nothing else is required. The demo is a web application. To start it type (in any folder):

	go run $GOPATH/src/github.com/icza/productws/cmd/proddemo/proddemo.go

You may also simply start the `$GOPATH/bin/proddemo` executable. (On Windows replace `$GOPATH` with `%GOPATH%`.)

The demo prints the address it's listening on (defaults to `":8081"`). You may override it with the `-addr` command line flag.

Test records are inserted on startup. To disable this, use the `-testdata=false` command line flag.


## Testing

For easy testing of the web service, check out the [tester-ui/tester.html](https://github.com/icza/productws/blob/master/doc.go)
simple HTML page built using [React](https://facebook.github.io/react/).
It provides UI for calling all the operations, allows you to edit data and see response. 

For automated testing you may use the [cURL](https://en.wikipedia.org/wiki/CURL) tool to query the web service.

To create a new product:

	curl -X POST -d "{\"Name\":\"JSCO Mouse\",\"Desc\":\"Computer Optical Noiseless Mouse\",\"Prices\":{\"USD\":{\"Value\":2782,\"Multiplier\":100}}}" localhost:8081/create

Example output:

	{"Op":"create","Success":true,"Data":{"ID":3}}

To list existing products:

	curl localhost:8081/list

Example response:

	{"Op":"list","Success":true,"Data":[1,2,3]}

To get the details of a product:

	curl localhost:8081/details/3

Example output:

	{"Op":"details","Success":true,"Data":{"ID":3,"Name":"JSCO Mouse","Desc":"Computer Optical Noiseless Mouse","Prices":{"USD":{"Value":2782,"Multiplier":100}}}}

To update a product (adding tags and GBP price):

	curl -X PUT -d "{\"ID\":3,\"Name\":\"JSCO Mouse\",\"Desc\":\"Computer Optical Noiseless Mouse\",\"Tags\":[\"Computer\",\"Mouse\"],\"Prices\":{\"USD\":{\"Value\":2782,\"Multiplier\":100},\"GBP\":{\"Value\":2093,\"Multiplier\":100}}}" localhost:8081/update

Example output:

	{"Op":"update","Success":true,"Data":{"ID":3}}

Let's verify the success of update with `curl localhost:8081/details/3`:

	{"Op":"details","Success":true,"Data":{"ID":3,"Name":"JSCO Mouse","Desc":"Computer Optical Noiseless Mouse","Tags":["Computer","Mouse"],"Prices":{"GBP":{"Value":2093,"Multiplier":100},"USD":{"Value":2782,"Multiplier":100}}}}

Set price points to different currencies (change GBP price and add HUF currency):

	curl -X PUT -d "{\"ID\":3,\"Prices\":{\"GBP\":{\"Value\":1999,\"Multiplier\":100},\"HUF\":{\"Value\":7717,\"Multiplier\":1}}}" localhost:8081/setprices

Example output:

	{"Op":"setprices","Success":true,"Data":{"ID":3}}

Let's verify the success of update with `curl localhost:8081/details/3`:

	{"Op":"details","Success":true,"Data":{"ID":3,"Name":"JSCO Mouse","Desc":"Computer Optical Noiseless Mouse","Tags":["Computer","Mouse"],"Prices":{"GBP":{"Value":1999,"Multiplier":100},"HUF":{"Value":7717,"Multiplier":1},"USD":{"Value":2782,"Multiplier":100}}}}

## Implementation details

The package documentation [doc.go](https://github.com/icza/productws/blob/master/doc.go) details the design choices
and gives an implementation overview.
It can also be viewed at [godoc.org](https://godoc.org/github.com/icza/productws).

## Authentication

The implementation does not include authentication.
If you need one, you could choose from the following options:

**Basic authentication**  
The client may use [Basic authentication](https://en.wikipedia.org/wiki/Basic_access_authentication)
which includes sending user+password with each request.  
Pros: Simple. Easy to implement. Supported by all browsers and clients.  
Cons: Password is sent unencrypted with all requests, should only be used over HTTPS.

**Using tokens** (OAuth 2.0 uses this too, see [RFC 6749](https://tools.ietf.org/html/rfc6749#section-7))   
A _token_ may be used instead of user+password. The token may be sent in HTTP headers, as request parameters
(or even in the request body).  
On first request (which may be a "special" authentication request or just a "regular" request)
the client sends authentication info. If they are valid, the server generates and sends a token back.
Subsequent requests only need to send this token.  
Pros: Tokens are independent from passwords. Tokens may have expiration time, they may be revoked arbitrarily,
they may be bound to IP etc.  
Cons: Slightly higher complexity; the server needs to maintain tokens (tell if a token is valid). 

## Making the service redundant

If we want the service to scale and / or to make it redundant, we have to replace the Store implementation
(obviously multiple service nodes needs to see the same data). Other than that, the service may be started
on multiple nodes without any problem. Multiple nodes may have and they may be reached at different addresses;
a load balancer / router may be started to coordinate requests and maintain equal distribution.

