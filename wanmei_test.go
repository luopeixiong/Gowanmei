package wanmei

import "testing"

func TestNewWmLib(t *testing.T) {
	x,err:=NewWmLibWithDLLPath(`ne.dat`,"163",`WmCode.dll`)
	if err!=nil{
		t.Error("err!=nil")
		t.FailNow()
	}
	t.Logf("%v",x)
}
func TestWanMeiLib_GetImageFromFile(t *testing.T) {
	x,_:=NewWmLibWithDLLPath(`ne.dat`,"163",`WmCode.dll`)
	y,_:=x.GetImageFromFile(`temp.bmp`)
	t.Logf(y)
}
