# 功能
* v2ray 订阅转 [clash](https://github.com/Dreamacro/clash)、[clashx](https://github.com/yichengchen/clashX)、[cfw](https://github.com/Fndroid/clash_for_windows_pkg/releases) 订阅
* v2ray 订阅转 QuantumultX 订阅 **给喜欢圈X的提供一个[IOS捷径](https://www.icloud.com/shortcuts/1fcda0a8ec9c4b48b899d5e9c9da7816)**
* ssr 订阅转 clashr 订阅, ssr兼容ss协议配置的订阅转clash订阅
# 主程序:
* macOS:    clashconfig-darwin-amd64
* windows:  clashconfig-windows-amd64
* Linux 64: clashconfig-linux-amd64
* 树莓派:    clashconfig-linux-armv7
# 使用方法:
1. 下载解压对应系统的主程序运行，主程序默认会监听 0.0.0.0:5050，也可以使用以下参数改变默认监听的地址和端口。**如果运行不成功，请手动下载[神机规则](https://raw.githubusercontent.com/ConnersHua/Profiles/master/Clash/Pro.yaml)保存为ConnersHua.yaml和主程序放在一起。**
```
Usage of ./clashconfig-linux-armv7:
  -h    this help
  -l string
        Listen address (default "0.0.0.0")
  -p string
        Listen Port (default "5050")
```
2. 使用下面格式设置订阅链接，QuantumultX、Clash设置可以直接用浏览器访问构造好的链接，点击导入即可。
#### v2ray转clash
```
http://127.0.0.1:5050/v2ray2clash?sub_link=此处换成需要转换的v2ray订阅链接
```
#### v2ray转QuantumultX
```
http://127.0.0.1:5050/v2ray2quanx?sub_link=此处换成需要转换的v2ray订阅链接
```
#### ssr转clashr
```
http://127.0.0.1:5050/ssr2clashr?sub_link=此处换成需要转换的ssr订阅链接
```
#### 测试地址(隐私自己考虑)
```
http://ne1l.tpddns.cn:5000/v2ray2clash?sub_link=此处换成需要转换的v2ray订阅链接
```

## 引用:
- [神机规则](https://github.com/ConnersHua/Profiles/tree/master)