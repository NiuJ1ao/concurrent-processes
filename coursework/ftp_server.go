package main

import (
	"fmt"
	"sync"
	"time"
)

type ClientAccount struct {
	userid string
	passwd string
}

type ClientChan struct {
	login    ClientAccount
	authAlts []chan int
	actAlts  []chan int
}

// Client(pid) = (νs) pid<s>.s<userid,passwd>.s|>{sorry: ...[] welcome: ...}
func Client(pid chan chan ClientChan, wg sync.WaitGroup) {
	// defer wg.Done()
	fmt.Printf("Client: Initiating pid=%p\n", pid)
	userid := "userid"
	passwd := "pw"
	s := make(chan ClientChan)
	fmt.Println("Client: Channel s=", s)
	login := ClientAccount{userid, passwd}

	pid <- s
	fmt.Printf("Client: pid=%p <- s=%p\n", pid, s)

	// authentication branching
	authAlts := []chan int{
		make(chan int), // sorry
		make(chan int), // welcome
	}
	// action branching
	actAlts := []chan int{
		make(chan int), // get
		make(chan int), // put
		make(chan int), // bye,
	}
	s <- ClientChan{
		login:    login,
		authAlts: authAlts,
		actAlts:  actAlts,
	}
	fmt.Printf("Client: s=%p <- ClientChan\n", s)

	fmt.Println("Client: Waiting for response...")
	select {
	case <-authAlts[0]:
		fmt.Println("Client: Sorry! User ID or password is incorrect.")
	case <-authAlts[1]:
		fmt.Println()
		fmt.Println("Client: Welcome!")
		var r int
		actAlts[0] <- r
		actAlts[0] <- r
		actAlts[0] <- r
		actAlts[1] <- r
		actAlts[1] <- r
		actAlts[1] <- r
		actAlts[2] <- r
		time.Sleep(200 * time.Millisecond)
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
	fmt.Printf("FtpThread: Initiating pid=%p, nis=%p\n", pid, nis)
	for {
		z := <-pid
		fmt.Printf("FtpThread: z=%p <- pid=%p\n", z, pid)
		clientChan := <-z
		fmt.Printf("FtpThread: clientChan <- z=%p\n", z)
		login := clientChan.login
		fmt.Printf("FtpThread: id=%s, pw=%s\n", login.userid, login.passwd)

		s := make(chan ClientChan)
		fmt.Println("FtpThread: Channel s'=", s)
		nis <- s
		fmt.Printf("FtpThread: nis=%p <- s'=%p\n", nis, s)

		// branching
		authAlts := []chan int{
			make(chan int), // invalid
			make(chan int), // valid
		}
		s <- ClientChan{
			login:    login,
			authAlts: authAlts,
		}
		fmt.Printf("FtpThread: s'=%p <- ClientChan\n", s)

		fmt.Println("FtpThread: Waiting for response...")
		select {
		case r := <-authAlts[0]:
			fmt.Printf("FtpThread: Invalid selected %p\n", authAlts[0])
			clientChan.authAlts[0] <- r
			fmt.Printf("FtpThread: Select sorry %p\n", clientChan.authAlts[0])
		case r := <-authAlts[1]:
			fmt.Printf("FtpThread: Valid selected %p\n", authAlts[1])
			clientChan.authAlts[1] <- r
			fmt.Printf("FtpThread: Select welcome %p\n", clientChan.authAlts[1])
			Actions(clientChan.actAlts)
		}
	}
}

func Actions(actAlts []chan int) {
	fmt.Println("Actions: Interacting...")
	for {
		select {
		case <-actAlts[0]:
			fmt.Println("Actions: Get")
		case <-actAlts[1]:
			fmt.Println("Actions: Put")
		case <-actAlts[2]:
			fmt.Println("Actions: Bye")
			return
		}
	}

}

func Authenticator(nis chan chan ClientChan) {
	fmt.Printf("Authenticator: Initiating nis=%p\n", nis)
	for {
		x := <-nis
		fmt.Printf("Authenticator: x=%x <- nis=%p\n", x, nis)
		clientChan := <-x
		fmt.Printf("Authenticator: ClientChan <- x=%p\n", x)
		login := clientChan.login
		fmt.Printf("Authenticator: id=%s, pw=%s\n", login.userid, login.passwd)

		if login.userid == "userid" && login.passwd == "pw" {
			var r int
			clientChan.authAlts[1] <- r
			fmt.Printf("Authenticator: Select valid %p\n", clientChan.authAlts[1])
		} else {
			var r int
			clientChan.authAlts[0] <- r
			fmt.Printf("Authenticator: Select invalid %p\n", clientChan.authAlts[0])
		}
	}
}

func main() {
	var wg sync.WaitGroup

	pid := make(chan chan ClientChan)
	nis := make(chan chan ClientChan)
	Init(pid, nis)

	// wg.Add(1)
	Client(pid, wg)

	wg.Wait()
	fmt.Println("Exiting")
}
