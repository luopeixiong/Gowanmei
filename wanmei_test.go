package wanmei

import (
	"bufio"
	"os"
	"testing"
)

const base  = ``
func TestNewWmLib(t *testing.T) {
	LoadDll(`WmCode.dll`)
	x,err:=NewWmLib(`ne.dat`,"163",)
	if err!=nil{
		t.Error("err!=nil")
		t.FailNow()
	}
	t.Logf("%v",x)
}
func TestWanMeiLib_GetImageFromFile(t *testing.T) {
	LoadDll(base+`WmCode.dll`)
	x,_:=NewWmLib(base+`ne.dat`,"163")
	y,_:=x.GetImageFromFile(base+`temp.bmp`)
	t.Logf(y)
}
func TestWanMeiLib_GetImageFromBuffer(t *testing.T) {
	x,_:=NewWmLib(`ne.dat`,"163")
	f,_:=os.Open(`temp.bmp`)
	bf:=bufio.NewReader(f)
	buf:=make([]byte,100000,100000)
	i,_:=bf.Read(buf)
	buf=buf[:i]
	y,_:=x.GetImageFromBuffer(buf)
	t.Logf(y)
}