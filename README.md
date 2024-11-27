## Logrus trace otel exporter

Useful when you do not have an OpenTelemetry collector yet and your application is built on top of logrus.

```
go run example.go

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.12.0
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
time="2024-11-27T14:22:52+01:00" level=debug msg="Trace 61c2ad8ea1cc2d44f32c3111a028e3fd span a90782262958b20c name dbUser duration 0.0000s" duration="1.793µs" id=1 name=dbUser parent_id=04bb4a8802e375ce span_id=a90782262958b20c trace_id=61c2ad8ea1cc2d44f32c3111a028e3fd

time="2024-11-27T14:22:52+01:00" level=debug msg="Trace 61c2ad8ea1cc2d44f32c3111a028e3fd span 04bb4a8802e375ce name getUser duration 0.0000s" duration="10.6µs" id=1 name=getUser parent_id=90fffc3395d65b98 span_id=04bb4a8802e375ce trace_id=61c2ad8ea1cc2d44f32c3111a028e3fd

time="2024-11-27T14:22:52+01:00" level=debug msg="Trace 61c2ad8ea1cc2d44f32c3111a028e3fd span 90fffc3395d65b98 name /users/:id duration 0.0001s" duration="124.975µs" http.method=GET http.route="/users/:id" http.scheme=http http.target=/users/1 name="/users/:id" net.host.name=my-server net.protocol.version=1.1 net.sock.peer.addr="::1" parent_id=0000000000000000 span_id=90fffc3395d65b98 trace_id=61c2ad8ea1cc2d44f32c3111a028e3fd user_agent.original=curl/8.9.1
```
