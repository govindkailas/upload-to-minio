# upload-to-minio

## Getting started 
This is a simple go microservice that exposes a `upload` endpoint to upload a file to minio. 

## Deploying Minio
Official documentation on how to deploy minio is [here](https://min.io/docs/minio/kubernetes/upstream/operations/install-deploy-manage/deploy-operator-helm.html). For the impatients here's a quick start: `kubectl apply -f k8s/minio.yaml`
This would create a namespace called `minio` and deploy minio in it. The service can be accessed on port 9000. The default credentails are set as `minio` and `minio123`. Once minio is up and running we can create a bucket and upload a file to it using the `go client` 

To run the go microservice:
`go run main.go` 

If run locally it can be accessible on port# 8080. Let's do a heath check: `http://localhost:8080/ping`
```
$ curl http://localhost:8080/ping
{"message": "Yooh, I am alive responding pong!!"}
```

So, good. Now let's upload a file to minio. We'll use the go client to do this:
```
$ curl http://localhost:8080/upload -F 'file=@k8s.png' -vv
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> POST /upload HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.4.0
> Accept: */*
> Content-Length: 33718
> Content-Type: multipart/form-data; boundary=------------------------reN4am6EJmf7mE4lujmxA4
>
* We are completely uploaded and fine
< HTTP/1.1 200 OK
< Date: Thu, 21 Dec 2023 23:22:53 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

_Note: `file` is the name of the file in the form data. `k8s.png` is the file we're uploading. Its important that `file` and `r.FormFile("file")` in `go` should match_

If everything worked as expected, the go client would have uploaded the file to minio. Let's verify this from the go app log as well, 
```
2023/12/21 15:20:24 We already own testbucket
File type is png
2023/12/21 15:22:53 File saved to file_1703200973100231892.png
2023/12/21 15:22:53 Uploading file to minio
2023/12/21 15:22:53 Successfully uploaded file_1703200973100231892.png of size 32 KB
```


Let's check if the file is present in the bucket: 
![minio bucket list](minio_testbucket.jpg)
