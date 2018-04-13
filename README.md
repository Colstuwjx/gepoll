# gepoll
Linux epoll test via Go

## test Kernel Message

```
# default device is /dev/kmsg
$ go build
$ ./gepoll
2018/04/13 12:09:54 data  12,25142,3000956334585,-;MESSAGE

# another terminal with below command, upper message would be output by gepoll
$ echo MESSAGE > /dev/kmsg
```
