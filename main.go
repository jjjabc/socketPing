package main

import (
	"fmt"
	"net"
	"time"
	"bytes"
	"encoding/binary"
	"errors"
)
const chapLen  = 8
var (
	errChapPasswordForm = errors.New("ChapPassword form is error!")
	errNotReplyPackage = errors.New("Wrong Package！")
	errPaseWrong = errors.New("Pase Error!")
	errNullPointer = errors.New("NullPointer!")
)
type State struct  {
	alive bool
}
func (this State) Bytes()([]byte){
	if this.alive{
return []byte{0x01}
}else{
		return []byte{0x00}
	}
}
func main() {
	rAddr, err := net.ResolveUDPAddr("udp", "192.168.1.242:6007")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	udpConn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer udpConn.Close()
	sendUDP(udpConn)
	if err!=nil{
		fmt.Println(err.Error())
	}
	receiveUDP(udpConn,rAddr)
}
func heartBeat(udpConn *net.UDPConn,rAddr *net.UDPAddr,duration time.Duration) (t time.Duration,err error){
	chap:=[]byte{0x11,0x22,0x33,0x44,0xaa,0xbb,0xcc,0xdd}
	sn:=uint32(1)
	diTime:=time.Now()
	_,err=di(udpConn,rAddr,chap,sn,getState())
	if err!=nil{
		return 0,err
	}
	rState,err:=dong(udpConn,rAddr,chap,sn,1*time.Second)
	t=time.Now().Sub(diTime)
	setState(rState)
	return t,err
}
func getState()State{
	return State{alive:true}
}
func setState(rState *State)error{
	return nil
}
func di(udpConn *net.UDPConn,rAddr *net.UDPAddr,chap []byte,sn uint32,state State)(n int,err error){
	sendBytesBuffer:=bytes.NewBuffer([]byte{})
	binary.Write(sendBytesBuffer,binary.BigEndian,n)
	len,_:=sendBytesBuffer.Write(chap)
	if len!=chapLen{
		return 0,errChapPasswordForm
	}
	len,_=sendBytesBuffer.Write(state.Bytes())
	return udpConn.Write(sendBytesBuffer.Bytes())
}
func dong(udpConn *net.UDPConn,rAddr *net.UDPAddr,chap []byte,sn uint32,timeout time.Duration)(state *State,err error){
	b:=make([]byte,1024)
	udpConn.SetReadDeadline(time.Now().Add(timeout))
	n,Addr,err:=udpConn.ReadFromUDP(b)
	if err!=nil{
		return nil,err
	}
	if (udpAddrEqual(Addr,rAddr)==false)||(checkReply(b[:n],chap,sn)) {
		return nil, errNotReplyPackage
	}
	return getStateFromBytes(b)
}
func getStateFromBytes(b []byte)(state *State,err error){
	if (b==nil){
		return nil,errNullPointer
	}
	if (len(b)<1){
		return nil,errPaseWrong
	}
	if b[0]==0x01{
		return &State{alive:true},nil
	}else if b[0]==0x02{
		return &State{alive:false},nil
	}
	return nil,errPaseWrong
}
func checkReply(replyBuf []byte,chap []byte,sn uint32)bool{
	return true
}
func sendUDP(udpConn *net.UDPConn) {
	sendString:=`Hello World!/usr/local/go/bin/go build -o /Users/jiaxun/IdeaProjects/socketTest/MacPing /Users/jiaxun/IdeaProjects/socketTest/main.go /Users/jiaxun/IdeaProjects/socketTest/MacPing`
	_, err := udpConn.Write([]byte(sendString))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func udpAddrEqual(addr1 *net.UDPAddr,addr2 *net.UDPAddr)bool{
	if addr1==nil||addr2==nil{
		return false
	}
	return addr1.IP.Equal(addr2.IP)&&(addr1.Port==addr2.Port)
}
func receiveUDP(udpConn *net.UDPConn,rAddr *net.UDPAddr) {
	b := make([]byte, 1024)
	udpConn.SetReadDeadline(time.Now().Add(5*time.Second))

	n, Addr, err := udpConn.ReadFromUDP(b)
	//	n,err:=udpConn.Read(b)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n > 1024 {
		fmt.Println("Buff out!")
		return
	}
	fmt.Println(udpConn.RemoteAddr().Network())
		if !Addr.IP.Equal(rAddr.IP){
			fmt.Println("IP diff:%s-%s",Addr.IP.String(),rAddr.IP.String())
			return
		}else if Addr.Port!=rAddr.Port{
			fmt.Println("Port diff:%d-%d",Addr.Port,rAddr.Port)
			return
		}
	fmt.Printf("Receive from %s:%s", udpConn.RemoteAddr().String(),b[:n])

}
func sendIPPacket() {
	fmt.Println("go!")
	ipaddr, err := net.ResolveIPAddr("ip", "192.168.1.254")
	if err != nil {
		fmt.Println("IP address is wrong!")
		return
	}
	localAddrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Local interface is null!")
		return
	}
	lIPAddrs := getSendLocalIP(localAddrs, ipaddr)
	if lIPAddrs == nil {
		fmt.Println("Local interface is null!")
		return
	}
	fmt.Println(lIPAddrs.String())
	ipConn, err := net.DialIP("ip", &net.IPAddr{IP: lIPAddrs.IP, Zone: ""}, ipaddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	b := []byte{0xAA, 0xBB, 0xCC}
	_, err = ipConn.Write(b)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer ipConn.Close()
}
func getSendLocalIP(localAddrs []net.Addr, rIPAddr *net.IPAddr) *net.IPNet {
	for _, lAddr := range localAddrs {
		if lIPNet, ok := lAddr.(*net.IPNet); ok && !lIPNet.IP.IsLoopback() {
			if lIPNet.Contains(rIPAddr.IP) {
				return lIPNet
			}
		}
	}
	return nil
}
