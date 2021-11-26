/*
 * In this question, we simulate the interaction between two clients and the FTP server in paralle.
 * One client will go for the correct user ID and password combination. Therefore, we should observe
 * "Welcome, <user_id>!" printed, where <user_id> is "userid" in this case. After that, a sequence of
 * actions is taken, which is 3 gets, 3 puts and bye. At the same time, another client types
 * incorrect user ID and password. "Sorry! User ID or password is incorrect." should be printed.
 * For the purpose of concision, we create a function for each process.
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
	s := make(chan ClientChan) // (νs)
	pid <- s                   // pid<s>

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

	select { //s|>{sorry: ... [] welcome: ...}
	case <-authAlts[0]: // sorry: 0
		fmt.Println("Client: Sorry! User ID or password is incorrect.")
	case <-authAlts[1]: // welcome: s <| get.s <| get.s <| get.s <| put.s <| put.s <| put.s <| bye
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
	for {
		z := <-pid // pid(z)
		clientChan := <-z
		loginInfo := clientChan.loginInfo // z(userid, passwd)

		s := make(chan ClientChan) // (νs')
		nis <- s                   // nis<s'>

		// branching
		authAlts := []chan int{
			make(chan int), // invalid
			make(chan int), // valid
		}
		s <- ClientChan{ // s'<userid, passwd>
			loginInfo: loginInfo,
			authAlts:  authAlts,
		}

		select { // s'|>{invalid : ... [] valid: ...}
		case r := <-authAlts[0]: // invalid : z<|sorry.0
			clientChan.authAlts[0] <- r
		case r := <-authAlts[1]: // valid: z <| welcome.Actions<z>
			clientChan.authAlts[1] <- r
			Actions(clientChan.actAlts)
		}
	}
}

// Actions(s) = s |> {get: Actions<s> [] put: Actions<s> [] bye: 0}
func Actions(actAlts []chan int) {
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
	for {
		x := <-nis // nis(x)
		clientChan := <-x
		loginInfo := clientChan.loginInfo // x(userid,passwd)

		if loginInfo.userid == "userid" && loginInfo.passwd == "pw" { // if ... then ... else ...
			var r int
			clientChan.authAlts[1] <- r
		} else {
			var r int
			clientChan.authAlts[0] <- r
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
