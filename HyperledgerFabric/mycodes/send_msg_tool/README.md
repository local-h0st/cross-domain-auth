# send msg tool

用来向server_vs进程发送消息的，测试用，语法如下:
```
./sendMsg [-ip x.x.x.x] -p port -m message  // -ip缺省为localhost

```
localhost显示被拒绝，换成0.0.0.0试一下

不是这个问题，缺省改回去了，防火墙规则也改了一下，方便以后测试

重新go build之后放到~目录下了方便使用