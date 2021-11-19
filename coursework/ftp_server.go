package main

import "fmt"

type Tag int

const (
	sorry   Tag = 0
	welcome Tag = 1
	get     Tag = 2
	put     Tag = 3
	buy     Tag = 4
)

type ClientAccount struct {
	userid string
	passwd string
}

type ClientChan struct {
	login ClientAccount
	p1    chan int
	p2    chan int
}

// Client(pid) = (νs) pid<s>.s<userid,passwd>.s|>{sorry: ...[] welcome: ...}
func Client(pid chan chan ClientChan) {
	userid := "userid"
	passwd := "pw"
	s := make(chan ClientChan)
	login := ClientAccount{userid, passwd}

	pid <- s
	s <- ClientChan{
		login: login,
		p1:    make(chan int), // sorry
		p2:    make(chan int), // welcome
	}
	sel := <-s
	select {
	case x := <-sel.p1:
		fmt.Printf("sorry: %d\n", x)
	case x := <-sel.p2:
		fmt.Printf("welcome: %d\n", x)
	}
}

// Init(pid, nis) = FtpThread<pid, nis> | Authenticator<nis>
func Init(pid chan chan ClientChan, nis chan chan ClientChan) {
	go FtpThread(pid, nis)
	go Authenticator(nis)
}

// FtpThread(pid, nis) = !pid(z).z(userid, passwd).(νs')nis<s'>.s'<userid, passwd>.
// s'|>{invalid : z<|sorry.0 [] valid: z <| welcome.Actions<z>}
func FtpThread(pid chan chan ClientChan, nis chan chan ClientChan) {
	for {
		z := <-pid
		clientChan := <-z
		login := clientChan.login

		s := make(chan ClientChan)
		nis <- s
		s <- ClientChan{
			login: login,
			p1:    make(chan int), // invalid
			p2:    make(chan int), // valid
		}

		// sel := <-s
		// select {
		// case x := <-sel.p1:
		// 	fmt.Printf("invalid: %d\n", x)
		// 	r := make(chan int)
		// 	clientChan.p1 <- r
		// case x := <-sel.p2:
		// 	fmt.Printf("valid: %d\n", x)

		// }
	}
}

func Actions(s chan int) {

}

func Authenticator(nis chan chan ClientChan) {
	for {
		x := <-nis
		clientChan := <-x
		login := clientChan.login

		if login.userid == "userid" && login.passwd == "pw" {

		} else {

		}
	}
}

func main() {

}
