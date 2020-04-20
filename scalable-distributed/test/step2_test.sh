#!/bin/bash

# 访问服务节点，创建一个叫test2的对象
curl -v 10.29.2.2:12345/objects/test2 -XPUT -d"This is object test2"

# Supposed Response:
#*   Trying 10.29.2.2:12345...
#* TCP_NODELAY set
#* Connected to 10.29.2.2 (10.29.2.2) port 12345 (#0)
#> PUT /objects/test2 HTTP/1.1
#> Host: 10.29.2.2:12345
#> User-Agent: curl/7.65.3
#> Accept: */*
#> Content-Length: 20
#> Content-Type: application/x-www-form-urlencoded
#>
#* upload completely sent off: 20 out of 20 bytes
#* Mark bundle as not supporting multiuse
#< HTTP/1.1 200 OK
#< Date: Tue, 14 Apr 2020 12:18:51 GMT
#< Content-Length: 0
#<
#* Connection #0 to host 10.29.2.2 left intact

# 看看test2被保存在哪个数据节点上？
curl 10.29.2.2:12345/locate/test2

# 换一个接口服务节点查看这个对象
curl 10.29.2.1:12345/objects/test2
