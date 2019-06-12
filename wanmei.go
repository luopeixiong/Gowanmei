package wanmei

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

type wanMeiDLL struct {
	dll *syscall.LazyDLL
	loadWmFromFileEx,setWmOptionEx,
	getImageFromFileEx,getImageFromBufferEx *syscall.LazyProc

}
type WanMeiLib struct {
	libId uintptr
}
type WmOption struct {
	RetType byte  //返回方式	取值范围：0～1   默认为0,直接返回验证码,为1返回验证码字符和矩形范围形如：
	             // S,10,11,12,13|A,1,2,3,4 表示识别到文本 S 左边横坐标10,左边纵坐标11,右边横坐标,右边纵坐标12
	SegmentationType byte  //识别方式(图像分割方法)    取值范围：0～4   默认为0,0整体识别,1连通分割识别,2纵分割识别,3横分割识别,4横纵分割识别
	RecogType byte  //识别模式	取值范围：0～1   默认为0,0识图模式,1为识字模式
	AccelerationType byte //识别加速	取值范围：0～1   默认为0,0为不加速,1为使用加速
	AccelerationRet byte  //加速返回	取值范围：0～1   默认为0,0为不加速返回,1为使用加速返回
	MinSimilarity byte  //最小相似度	取值范围：0～100 默认为90 ,在go的实现中0同90,所以不允许出现0
	CharSpace int  //字符间隙    取值范围：-10～0   默认为0,如果字符重叠,根据实际情况填写,如-3允许重叠3像素,如果不重叠的话,直接写0
}
var lock sync.Mutex=sync.Mutex{}
var wmDll *wanMeiDLL=nil
func LoadDll(path string) error {
	lock.Lock()
	defer lock.Unlock()
	if wmDll==nil {
		wmDll=new(wanMeiDLL)
		wmDll.dll=syscall.NewLazyDLL(path)
		ret:=wmDll.dll.Load()
		if ret!=nil {
			return ret
		}
		wmDll.getImageFromFileEx=wmDll.dll.NewProc("GetImageFromFileEx")
		wmDll.getImageFromBufferEx=wmDll.dll.NewProc("GetImageFromBufferEx")
		wmDll.loadWmFromFileEx=wmDll.dll.NewProc("LoadWmFromFileEx")
		wmDll.setWmOptionEx=wmDll.dll.NewProc("SetWmOptionEx")
	}
	return nil
}
func NewWmLib(path ,password string) (ret *WanMeiLib,err error) {
	if wmDll==nil {
		e:=LoadDll("WmCode.dll")
		if e!=nil{
			return nil,e
		}
	}
	ret=new(WanMeiLib)
	strPath:=append([]byte(path),0)
	strPass:=append([]byte(password),0)
	ret.libId,_,_=wmDll.loadWmFromFileEx.Call(uintptr(unsafe.Pointer(&strPath[0])),uintptr(unsafe.Pointer(&strPass[0])))
	if int(ret.libId)==-1 {
		return nil,errors.New("load wm file error!")
	} else {
		return ret,nil
	}
}


func (self *WanMeiLib) SetWmOption(op * WmOption){
	setWmOption:=wmDll.setWmOptionEx
	if op.RetType!=0{
		setWmOption.Call(self.libId ,1,uintptr(op.RetType))
	}
	if op.SegmentationType!=0{
		setWmOption.Call(self.libId ,2,uintptr(op.SegmentationType))
	}
	if op.RecogType!=0{
		setWmOption.Call(self.libId ,3,uintptr(op.RecogType))
	}
	if op.AccelerationType!=0{
		setWmOption.Call(self.libId ,4,uintptr(op.AccelerationType))
	}
	if op.AccelerationRet!=0{
		setWmOption.Call(self.libId ,5,uintptr(op.AccelerationRet))
	}
	if op.MinSimilarity!=0 && op.MinSimilarity!=90{
		setWmOption.Call(self.libId ,6,uintptr(op.MinSimilarity))
	}
	if op.CharSpace!=0{
		setWmOption.Call(self.libId ,7,uintptr(op.CharSpace))
	}
}
func (self *WanMeiLib)  GetImageFromFile( filepath string) (ver string,err error){
	getImageFromFile:=wmDll.getImageFromFileEx
	strbuf:=append([]byte(filepath),0)
	retbuf:=make([]byte,5000,5000)
	ok,_,_:=getImageFromFile.Call(self.libId ,uintptr(unsafe.Pointer(&strbuf[0])),uintptr(unsafe.Pointer(&retbuf[0])))
	if ok==0 {
		return "",errors.New("GetImageFromFile error")
	}
	slen:=0
	for _,v:= range retbuf{
		slen+=1
		if v==0{
			break
		}
	}
	retbuf=retbuf[:slen]
	return string(retbuf),nil
}
func (self *WanMeiLib)  GetImageFromBuffer( buff []byte) (ver string,err error){
	getImageFromBuffer:=wmDll.getImageFromBufferEx
	retbuf:=make([]byte,5000,5000)
	ok,_,_:=getImageFromBuffer.Call(self.libId ,uintptr(unsafe.Pointer(&buff[0])),uintptr(len(buff)),uintptr(unsafe.Pointer(&retbuf[0])))
	if ok==0 {
		return "",errors.New("GetImageFromBuffer error")
	}
	slen:=0
	for _,v:= range retbuf{
		slen+=1
		if v==0{
			break
		}
	}
	retbuf=retbuf[:slen]
	return string(retbuf),nil
}