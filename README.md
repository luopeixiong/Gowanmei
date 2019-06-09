# wanmei-go
完美验证码识别的go链接,如你所想,通过cgo实现
因原库只支持32位,因此请使用go386 编译
## 基本用法
//x,err:=NewWmLib(识别文件目录,识别文件密码)  使用缺省本地wmcode.dll
x,err:=NewWmLibWithDLLPath(识别文件目录,识别文件密码,dll文件位置)
	if err!=nil{
		panic("err!")
	}
y,err:=x.GetImageFromFile(文件名)
	if err!=nil{
		panic("err!")
	}
或者
f,_:=os.Open(文件名)
bf:=bufio.NewReader(f)
buf:=make([]byte,100000,100000)
i,_:=bf.Read(buf)
buf=buf[:i]
y,_:=x.GetImageFromBuffer(buf)
