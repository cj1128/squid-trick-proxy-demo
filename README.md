# Squid Demo

基于Squid的一个恶作剧代理，会将网站中的图片进行翻转。关于代理的知识，可以参考我的博客[正向代理与反向代理]。

<p align="center">
  <img src="http://ww1.sinaimg.cn/large/9b85365dgy1fu6u1kzebaj21ha0s0e82" width="800px" />
</p>

这个恶作剧代理旨在说明Squid的基本使用，需要显示配置客户端，如果真的要给你的邻居或者室友一点惊喜，可以开放一个公开热点，然后配置透明代理，此时，任何连上你的Wifi的人，都会情不自禁的大喊一声：What the fuck...

注意：**由于拦截HTTPS流量涉及到证书问题，该代理只能拦截HTTP流量，对于HTTPS的图片，是无法进行翻转的**。

## 如何使用

1. 将仓库克隆到`/tmp/squid`目录中

    ```bash
    $ git clone https://github.com/fate-lovely/squid-demo /tmp/squid
    $ cd /tmp/squid
    ```

2. 编译`rewrite`程序

    ```bash
    $ go build rewrite.go
    ```

3. 启动代理

    ```bash
    $ squid -N -f squid.conf
    ```

4. 启动一个服务器返回翻转以后的图片，这里我们选择使用`http-server`

    ```bash
    $ yan global add http-server
    $ http-server -c-1 cache
    ```

5. 安装[Proxy SwitchyOmega]，配置使用代理，协议选择`HTTP`，Server填写`127.0.0.1`，端口填写`3128`

访问任意含有图片的HTTP站点，比如[百度图片搜索]，Enjoy~ 😉。

[正向代理与反向代理]: http://cjting.me/misc/forward-proxy-and-reverse-proxy/
[Proxy SwitchyOmega]: https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif?utm_source=chrome-ntp-icon
[百度图片搜索]: http://image.baidu.com/search/index?tn=baiduimage&ipn=r&ct=201326592&cl=2&lm=-1&st=-1&fm=result&fr=&sf=1&fmq=1534061918470_R&pv=&ic=0&nc=1&z=&se=1&showtab=0&fb=0&width=&height=&face=0&istype=2&ie=utf-8&word=%E6%B1%A4%E5%94%AF
