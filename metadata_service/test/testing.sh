#!/bin/bash

# 发送一个test3到apiServer 当然会返回400 Bad Request，因为未在请求头部提供散列值
curl -v 10.29.2.2:12345/objects/test3 -XPUT -d"This is object test3"
# *   Trying 10.29.2.2:12345...
#* TCP_NODELAY set
#* Connected to 10.29.2.2 (10.29.2.2) port 12345 (#0)
#> PUT /objects/test3 HTTP/1.1
#> Host: 10.29.2.2:12345
#> User-Agent: curl/7.65.3
#> Accept: */*
#> Content-Length: 20
#> Content-Type: application/x-www-form-urlencoded
#>
#* upload completely sent off: 20 out of 20 bytes
#2020/04/20 10:58:09 error handling put method: missing object hash in digest header
#* Mark bundle as not supporting multiuse
#< HTTP/1.1 400 Bad Request
#< Date: Mon, 20 Apr 2020 02:58:09 GMT
#< Content-Length: 0
#<
#* Connection #0 to host 10.29.2.2 left intact

# 好吧，现在用openssl计算一下散列值（计算散列值方法很多）
echo -n "This is object test3" | openssl dgst -sha256 -binary | base64
# 6LYX5AfsHkCQxsqm5b4pBWU9N8xuFxkxhmniDeHTi+E=

curl -v 10.29.2.2:12345/objects/test3 -XPUT -d"This is object test3" -H "Digest: SHA-256=6LYX5AfsHkCQxsqm5b4pBWU9N8xuFxkxhmniDeHTi+E="

# 添加一个下新版本
echo -n "This is object test3 version 2" | openssl dgst -sha256 -binary | base64
# Vop5VVs9DjGGIxsPqIfSXYbNVnAPYylRcqBTpmIwbrA=

# 看看他们保存在哪个节点上？
curl 10.29.2.1:12345/locate/6LYX5AfsHkCQxsqm5b4pBWU9N8xuFxkxhmniDeHTi+E=

curl 10.29.2.1:12345/locate/Vop5VVs9DjGGIxsPqIfSXYbNVnAPYylRcqBTpmIwbrA=

curl 10.29.2.1:12345/versions/test3

curl 10.29.2.1:12345/objects/test3?version=2

curl --request DELETE -sL \
     --url 'http://10.29.2.1:12345/objects/test3'
