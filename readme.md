### log-proxy

a reverse proxy to a single address, logging request and response bodies

###### installation 

`$ go get -u go.senan.xyz/log-proxy`

###### usage 

```
  -listen-addr string
    	address to listen on, eg. :5050
  -to string
    	address to proxy to, eg. http://localhost:4040
```

eg. say I have a service running on my laptop on port 4040. it doesn't log requests and the usual client doesn't log anything

`$ log-proxy -listen-addr :5050 -to http://localhost:4040`

now point the client at 5050 instead, and see the requests and responses in sequence

```
2020/04/27 17:23:02 listening on ":5050"

######### (1) request #########
GET /rest/ping.view?u=admin&p=admin&c=c&v=1.15 HTTP/1.1
Host: localhost:4040
Accept: */*
User-Agent: curl/7.69.1
X-Forwarded-For: ::1
X-Forwarded-Host: 

######### (1) response #########
HTTP/1.1 200 
Content-Length: 125
Access-Control-Allow-Origin: *
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Content-Type: text/xml;charset=UTF-8
Date: Mon, 27 Apr 2020 16:23:08 GMT
Expires: 0
Pragma: no-cache
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Xss-Protection: 1; mode=block

<?xml version="1.0" encoding="UTF-8"?>
<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.15.0"/>

######### (2) request #########
GET /rest/getIndexes.view?u=admin&p=admin&c=c&v=1.15 HTTP/1.1
Host: localhost:4040
Accept: */*
User-Agent: curl/7.69.1
X-Forwarded-For: ::1
X-Forwarded-Host: 

######### (2) response #########
HTTP/1.1 200 
Content-Length: 328
Access-Control-Allow-Origin: *
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Content-Type: text/xml;charset=UTF-8
Date: Mon, 27 Apr 2020 16:23:12 GMT
Expires: 0
Pragma: no-cache
Set-Cookie: player-61646d696e=2; Max-Age=31536000; Expires=Tue, 27-Apr-2021 16:23:12 GMT; Path=/; HttpOnly
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Xss-Protection: 1; mode=block

<?xml version="1.0" encoding="UTF-8"?>
<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.15.0">
   <indexes lastModified="1587989053520" ignoredArticles="The El La Los Las Le Les">
      <index name="W">
         <artist id="1" name="Wagon Christ"/>
      </index>
   </indexes>
</subsonic-response>

######### (3) request #########
GET /rest/ping.view?u=admin&p=admin&c=c&v=1.15 HTTP/1.1
Host: localhost:4040
Accept: */*
User-Agent: curl/7.69.1
X-Forwarded-For: ::1
X-Forwarded-Host: 

######### (3) response #########
HTTP/1.1 200 
Content-Length: 125
Access-Control-Allow-Origin: *
Cache-Control: no-cache, no-store, max-age=0, must-revalidate
Content-Type: text/xml;charset=UTF-8
Date: Mon, 27 Apr 2020 16:23:15 GMT
Expires: 0
Pragma: no-cache
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Xss-Protection: 1; mode=block

<?xml version="1.0" encoding="UTF-8"?>
<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.15.0"/>
signal: interrupt
```
