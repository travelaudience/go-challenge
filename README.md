travel audience Go challenge
============================

Task
----

Write an HTTP service that exposes an endpoint "/numbers". This endpoint receives a list of URLs 
though GET query parameters. The parameter in use is called "u". It can appear 
more than once.

	http://yourserver:8080/numbers?u=http://example.com/primes&u=http://foobar.com/fibo

When the /numbers is called, your service shall retrieve each of these URLs if 
they turn out to be syntactically valid URLs. Each URL will return a JSON data 
structure that looks like this:

	{ "numbers": [ 1, 2, 3, 5, 8, 13 ] }

The JSON data structure will contain an object with a key named "numbers", and 
a value that is a list of integers. After retrieving each of these URLs, the 
service shall merge the integers coming from all URLs, sort them in ascending 
order, and make sure that each integer only appears once in the result. The 
endpoint shall then return a JSON data structure like in the example above with 
the result as the list of integers.

The endpoint needs to return the result as quickly as possible, but always 
within 500 milliseconds. It needs to be able to deal with error conditions when 
retrieving the URLs. If a URL takes too long to respond, it must be ignored. It 
is valid to return an empty list as result only if all URLs returned errors or 
took too long to respond. The timeout must be respected regardless of 
the size of the data.

The service should widely implement HTTP status codes and use them as documented in the [the specification](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html).
Consider different status codes for different cases.

Example
-------

The service receives an HTTP request:

	>>> GET /numbers?u=http://example.com/primes&u=http://foobar.com/fibo HTTP/1.0

It then retrieves the URLs specified as parameters.

The first URL returns this response:

	>>> GET /primes HTTP/1.0
	>>> Host: example.com
	>>> 
	<<< HTTP/1.0 200 OK
	<<< Content-Type: application/json
	<<< Content-Length: 34
	<<< 
	<<< { "numbers": [ 2, 3, 5, 7, 11, 13 ] }

The second URL returns this response:

	>>> GET /fibo HTTP/1.0
	>>> Host: foobar.com
	>>> 
	<<< HTTP/1.0 200 OK
	<<< Content-Type: application/json
	<<< Content-Length: 40
	<<< 
	<<< { "numbers": [ 1, 1, 2, 3, 5, 8, 13, 21 ] }

The service then calculates the result and returns it.

	<<< HTTP/1.0 200 OK
	<<< Content-Type: application/json
	<<< Content-Length: 44
	<<< 
	<<< { "numbers": [ 1, 2, 3, 5, 7, 8, 11, 13, 21 ] }


Completion Conditions
---------------------

Solve the task described above using Go. Only use what's provided in the Go 
standard library. The resulting program must run stand-alone with no other 
dependencies than the Go compiler.

Document your source code, both using comments and in a separate text file that 
describes the intentions and rationale behind your solution. Also write down 
any ambiguities that you see in the task description, and describe you how you 
interpreted them and why. If applicable, write automated tests for your code.

For testing purposes, you will be provided with an example server that, when 
run, listens on port 8090 and provides the endpoints /primes, /fibo, /odd and 
/rand.

Please return your working solution within 7 days of receiving the challenge.
