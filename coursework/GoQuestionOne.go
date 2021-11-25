/*
 */

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
	loginInfo ClientAccount
	authAlts  []chan int
	actAlts   []chan int
}

// Client(pid) = (νs) pid<s>.s<userid,passwd>.s|>{sorry: 0 [] welcome: s <| get.s <| get.s <| get.s <| put.s <| put.s <| put.s <| bye}
func Client(pid chan chan ClientChan, loginInfo ClientAccount) {
	fmt.Printf("Client: Initiating pid=%p\n", pid)
	s := make(chan ClientChan) // (νs)
	fmt.Println("Client: Channel s=", s)
	pid <- s // pid<s>
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
	s <- ClientChan{ // s<userid,passwd>
		loginInfo: loginInfo,
		authAlts:  authAlts,
		actAlts:   actAlts,
	}
	fmt.Printf("Client: s=%p <- ClientChan\n", s)

	fmt.Println("Client: Waiting for response...")
	select { //s|>{sorry: ... [] welcome: ...}
	case <-authAlts[0]: // sorry: 0
		fmt.Println("Client: Sorry! User ID or password is incorrect.")
	case <-authAlts[1]: // welcome: s <| get.s <| get.s <| get.s <| put.s <| put.s <| put.s <| bye
		fmt.Println()
		fmt.Printf("Client: Welcome, %s!\n", loginInfo.userid)
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
	go FtpThread(pid, nis) // FtpThread<pid, nis>
	go Authenticator(nis)  // Authenticator<nis>
}

// FtpThread(pid, nis) = !pid(z).z(userid, passwd).(νs')nis<s'>.s'<userid, passwd>.
// s'|>{invalid : z<|sorry.0 [] valid: z <| welcome.Actions<z>}
func FtpThread(pid chan chan ClientChan, nis chan chan ClientChan) {
	fmt.Printf("FtpThread: Initiating pid=%p, nis=%p\n", pid, nis)
	for {
		z := <-pid // pid(z)
		fmt.Printf("FtpThread: z=%p <- pid=%p\n", z, pid)
		clientChan := <-z
		fmt.Printf("FtpThread: clientChan <- z=%p\n", z)
		loginInfo := clientChan.loginInfo // z(userid, passwd)
		fmt.Printf("FtpThread: id=%s, pw=%s\n", loginInfo.userid, loginInfo.passwd)

		s := make(chan ClientChan) // (νs')
		fmt.Println("FtpThread: Channel s'=", s)
		nis <- s // nis<s'>
		fmt.Printf("FtpThread: nis=%p <- s'=%p\n", nis, s)

		// branching
		authAlts := []chan int{
			make(chan int), // invalid
			make(chan int), // valid
		}
		s <- ClientChan{ // s'<userid, passwd>
			loginInfo: loginInfo,
			authAlts:  authAlts,
		}
		fmt.Printf("FtpThread: s'=%p <- ClientChan\n", s)

		fmt.Println("FtpThread: Waiting for response...")
		select { // s'|>{invalid : ... [] valid: ...}
		case r := <-authAlts[0]: // invalid : z<|sorry.0
			fmt.Printf("FtpThread: Invalid selected %p\n", authAlts[0])
			clientChan.authAlts[0] <- r
			fmt.Printf("FtpThread: Select sorry %p\n", clientChan.authAlts[0])
		case r := <-authAlts[1]: // valid: z <| welcome.Actions<z>
			fmt.Printf("FtpThread: Valid selected %p\n", authAlts[1])
			clientChan.authAlts[1] <- r
			fmt.Printf("FtpThread: Select welcome %p\n", clientChan.authAlts[1])
			Actions(clientChan.actAlts)
		}
	}
}

// Actions(s) = s |> {get: Actions<s> [] put: Actions<s> [] bye: 0}
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

// Authenticator(nis) = !nis(x).x(userid,passwd).if passwd=pw then x<|valid.0 else x<|invalid.0
func Authenticator(nis chan chan ClientChan) {
	fmt.Printf("Authenticator: Initiating nis=%p\n", nis)
	for {
		x := <-nis // nis(x)
		fmt.Printf("Authenticator: x=%x <- nis=%p\n", x, nis)
		clientChan := <-x
		fmt.Printf("Authenticator: ClientChan <- x=%p\n", x)
		loginInfo := clientChan.loginInfo // x(userid,passwd)
		fmt.Printf("Authenticator: id=%s, pw=%s\n", loginInfo.userid, loginInfo.passwd)

		if loginInfo.userid == "userid" && loginInfo.passwd == "pw" { // if ... then ... else ...
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		Client(pid, ClientAccount{"userid", "pw"}) // should print: 'Client: Welcome, userid!'
	}()
	Client(pid, ClientAccount{"userid1", "pw1"})

	wg.Wait()
	fmt.Println("Exiting")
}
