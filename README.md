# dnsVerifier

一款批量验证dns服务器可用性的工具

## 介绍

现在子域名暴力破解为了提升速度，很多工具都会用大量的dns服务器（如massdns），一般来说验证一个dns服务器是否可用只需找个存在的域名看能否解析以及结果是否正确即可，但是最近发现部分dns解析正确域名是没问题的，但是给一个不存在的域名就会超时，影响了枚举速度，还有些dns服务器会把不存在的域名解析到一个ip上，当对枚举结果进行递归枚举时，这种大量的错误结果会特别耽误时间，于是写了这个小工具可以批量验证dns服务器的可用性，除了上面提到的点外，还做了超时验证（默认300ms），这样比较慢的dns也会排除掉，提高子域名枚举的速度，大家可以利用fofa、zoomeye等寻找dns服务器，然后结合本程序验证dns，项目中的vaildResolvers_1.0.txt为验证后的dns服务器列表，可以直接使用

## 安装

```
git clone https://github.com/alwaystest18/dnsVerifier.git
cd dnsVerifier/
go install
go build dnsVerifier.go
```

## 使用

参数说明

```
Usage of ./dnsVerifier:
  -o string         //输出文件名称，默认为vaildResolvers.txt
        output vaild resolvers to file (default "vaildResolvers.txt")
  -r string         //dns服务器列表文件，每行一个，格式为1.1.1.1或1.1.1.1:53
        dns resolvers file
  -rate int         //协程数，默认500
        rate limit (default 500)
  -t int            //线程数，默认100
        Number of threads (default 100)
  -timeout int      //请求dns超时时间，单位毫秒
        DNS maximum connection time(ms) (default 300)
```

使用

```
$ ./dnsVerifier -r nameservers.txt 
```

实测10467个dns服务家庭宽带默认线程配置耗时7s

## 附带文件说明

vaildResolvers_1.0.txt  国内家庭宽带验证后的可用dns列表

resolvers_all.txt   收集的全世界dns列表，由于大家服务器可能在海外，因此可以通过本程序验证该列表找到适合所在地区的dns后使用

数据来源：

https://public-dns.info/nameserver/cn.txt

https://public-dns.info/nameservers.txt

https://github.com/neargle/public-dns-list


