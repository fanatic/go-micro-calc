The ultimate goal is to run envoy via HTTP_PROXY= and transparently proxy requests from the application.

The HTTP Connection Manager expects the path to be relative, but when using HTTP_PROXY, the path is absolute.

https://github.com/lyft/envoy/blob/73c381b997e8e787fb783bc2d4b89be8338df2bb/source/common/http/conn_manager_impl.cc#L443-L453

So, apply this patch to strip the absolute part of the url and rebuild envoy.

```
--- a/source/common/http/http1/codec_impl.cc
+++ b/source/common/http/http1/codec_impl.cc
@@ -428,7 +428,12 @@ void ServerConnectionImpl::onMessageBegin() {
 }
 
 void ServerConnectionImpl::onUrl(const char* data, size_t length) {
+  ENVOY_CONN_LOG(debug, "length={} url={}", connection_, length, data);
   if (active_request_) {
+    if (length == 16) {
+       data += 11;     
+       length -= 11;
+    }
     active_request_->request_url_.append(data, length);
   }
 }
```

```
bazel build //source/exe:envoy-static
HTTP_PROXY=http://127.0.0.1:3001 PLUS_SVC_URL=http://plus/plus PORT=3003 $GOPATH/bin/go-micro-calc
./bazel-bin/source/exe/envoy-static -c proxy.json -l debug --service-cluster microcalc

```
