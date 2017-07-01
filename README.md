# go-micro-calc
Calculator microservice in Go

## Running

First, start this Go daemon on port 3003, and let it know where to find the envoy egress.  Ideally, we'd set HTTP\_PROXY and let PLUS\_SVC\_URL point to the real url, but that's not working.

```
$ PLUS_SVC_URL=http://127.0.0.1:3001/plus PORT=3003 $GOPATH/bin/go-micro-calc 
Listening on :3003plus at=info 4.000000 + 4.000000 = 8.000000
2017/07/01 21:16:48 [e2c1db4e9a9e/sZNqZUpala-000002] "POST http://127.0.0.1:3001/plus HTTP/1.1" from 172.17.0.2 - 200 13B in 278.53µs
plus at=info 8.000000 + 4.000000 = 12.000000
2017/07/01 21:16:48 [e2c1db4e9a9e/sZNqZUpala-000003] "POST http://127.0.0.1:3001/plus HTTP/1.1" from 172.17.0.2 - 200 14B in 200.804µs
mul at=info 3.000000 * 4.000000 = 12.000000
2017/07/01 21:16:48 [e2c1db4e9a9e/sZNqZUpala-000001] "POST http://localhost:3000/mul HTTP/1.1" from 127.0.0.1:59832 - 200 14B in 12.368728ms
plus at=info 4.000000 + 4.000000 = 8.000000
```

Next, start up zipkin.  It's accessible at http://localhost:9411

```
$ docker run -d -p 9411:9411 openzipkin/zipkin
```

Finally, start up envoy with the provided proxy configuration.

```
Build and compile envoy, start it like..
$ ./bazel-bin/source/exe/envoy-static -c envoy-proxy.json -l debug --service-cluster microcalc
```

Now you can test by curling the envoy ingress listener.

```
$ curl http://localhost:3000/mul -X POST -d '{"a":3, "b": 4}'
{"answer":12}
```