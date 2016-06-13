package main
/******************** Testing Objective consensu: Sending Invoke Requests in Parallel ********
*   Setup: 8 node local docker peer network with security
*   0. Deploy chaincodeexample02 with 100000000, 1 as initial args
*   1. Send Invoke Requests in parallel on multiple peers using go routines.
    2. Logic is set so that one complete set of invoke requests on 8 peers would send 144 requests.
*   3. After each such one loop, verify query results match on PEER0 and PEER1 
*********************************************************************/


import (
	"fmt"
	"strconv"
	"time"

	"obcsdk/chaincode"
	"obcsdk/peernetwork"
)

func main() {

	//var MyNetwork peernetwork.PeerNetwork

	fmt.Println("Creating a local docker network")
	peernetwork.SetupLocalNetwork(5, true)

        //with the InvokeLoop designed with this 8 peer 4 req setup we would have 144 requests in one round
        numPeers := 5
	numReq := 4

	time.Sleep(60000 * time.Millisecond)
	peernetwork.PrintNetworkDetails()
	_ = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	//chaincode.Init()
	chaincode.RegisterUsers()


        //os.Exit(1)
	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer, _ := peernetwork.APeer(chaincode.ThisNetwork)
	url := "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

        fmt.Println("Peers on network ")
	chaincode.NetworkPeers(url)

	fmt.Println("\nPOST/Chaincode: Deploying chaincode at the beginning ....")


	dAPIArgs0 := []string{"example02", "init"}
	depArgs0 := []string{"a", "100000000", "b", "0"}
	chaincode.Deploy(dAPIArgs0, depArgs0)
	//fmt.Println("From Deploy error ", err)

        //os.Exit(1)

	time.Sleep(60000 * time.Millisecond);
	fmt.Println("\nPOST/Chaincode: Querying a and b after deploy >>>>>>>>>>> ")

	qAPIArgs0 := []string{"example02", "query"}
	qArgsa := []string{"a"}
	qArgsb := []string{"b"}
	A, _ := chaincode.Query(qAPIArgs0, qArgsa)
	B, _ := chaincode.Query(qAPIArgs0, qArgsb)
	myStr := fmt.Sprintf("\nA = %s B= %s", A,B)
	fmt.Println(myStr)


	InvokeLoop(numPeers, numReq)

}


func InvokeLoop(numPeers int, numReq int) {

	var inita, initb, curra, currb  int
	inita = 100000000
	initb = 0
	curra = inita
	currb = initb

	invArgs0 := []string{"a", "b", "1"}
	//i := 0
	for  {
		j := 0
		k := 1
		for ( j < numPeers  ) {
			k = 1
			//fmt.Println("Value in j ", j)
			currPeer := "PEER" + strconv.Itoa(j)
			iAPIArgsCurrPeer := []string{"example02", "invoke", currPeer}
			for (k <= numReq) {
			 	go chaincode.InvokeOnPeer(iAPIArgsCurrPeer, invArgs0)
			 	//fmt.Println("Invoking ", k)
			 	k++
	  	}
			m:=j-1
			for (m >=0 ) {
			 	loopPeer := "PEER" + strconv.Itoa(m)
			 	//fmt.Println("Value in loopPeer ", loopPeer)
			 	iAPIArgsLoopPeer := []string{"example02", "invoke", loopPeer}
			 	k = 1
		   	for (k <= numReq) {
					go chaincode.InvokeOnPeer(iAPIArgsLoopPeer, invArgs0)
					//fmt.Println("Invoking LoopPeer", k)
				  k++
		   	}
				m = m-1
			}
  		j++;
		}
		curra = curra - 144
		currb = currb + 144
		QueryMatch(curra, currb)
    //i++
	}
}

func QueryMatch(curra int, currb int) {

  fmt.Println("Inside Query match ********************************* &&&&&&&&&&&&&& %%%%%%%%%%%%%%%%")
	fmt.Println("Sleeping for 2 minutes ")
	time.Sleep(120000 * time.Millisecond)

	fmt.Println("\nPOST/Chaincode: Querying a and b after invoke >>>>>>>>>>> ")
	qAPIArgs00 := []string{"example02", "query", "PEER0"}
	qAPIArgs01 := []string{"example02", "query", "PEER1"}
	qArgsa := []string{"a"}
	qArgsb := []string{"b"}

	res0A, _ := chaincode.QueryOnHost(qAPIArgs00, qArgsa)
	res0B, _ := chaincode.QueryOnHost(qAPIArgs00, qArgsb)

	res0AI, _ := strconv.Atoi(res0A)
	res0BI, _ := strconv.Atoi(res0B)

	res1A, _ := chaincode.QueryOnHost(qAPIArgs01, qArgsa)
	res1B, _ := chaincode.QueryOnHost(qAPIArgs01, qArgsb)

	res1AI, _ := strconv.Atoi(res1A)
	res1BI, _ := strconv.Atoi(res1B)

	if (curra == res0AI) && (currb == res0BI) {
		fmt.Println("Results in a and b match with Invoke values on PEER0: PASS")
		valueStr := fmt.Sprintf(" curra : %d, currb : %d, resa : %d , resb : %d", curra, currb, res0AI, res0BI)
		fmt.Println(valueStr)
	} else {
		fmt.Println("******************************")
		fmt.Println("RESULTS DO NOT MATCH on PEER0 : FAIL")

		fmt.Println("******************************")
	}

	if (curra == res1AI) && (currb == res1BI) {
		fmt.Println("Results in a and b match with Invoke values on PEER1: PASS")
		valueStr := fmt.Sprintf(" curra : %d, currb : %d, resa : %d , resb : %d", curra, currb, res1AI, res1BI)
		fmt.Println(valueStr)
	} else {
		fmt.Println("******************************")
		fmt.Println("RESULTS DO NOT MATCH on PEER1 : FAIL")
		fmt.Println("******************************")
	}
}