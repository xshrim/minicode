package main

import (
	"bufio"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/params"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
)

var ethereumBootnodes = []string{
	"enode://91922b12115c067005c574844c6bbdb114eb262f90b6355cec89e13b483c3e4669c6d63ec66b6e3ca7a3a462d28edb3c659e9fa05ed4c7234524e582a8816743@120.27.164.92:13333", // CH
}

var (
	srv              *p2p.Server
	peerCh           chan p2p.PeerEvent
	wg               sync.WaitGroup
	hasPeerConnected bool

	w        *whisper.Whisper
	topic    whisper.TopicType
	symKey   []byte
	asymKey  *ecdsa.PrivateKey
	filterID string
)

var (
	argTopic  = flag.String("topic", "abcd", "less 4 bytes topic")
	argPasswd = flag.String("room-password", "123456", "password for generating symKey")
	argPoW    = flag.Float64("pow", 0.2, "The PoW of local node")
)

func whisperConfig() {
	var (
		err       error
		keyID     string
		asymKeyID string
	)

	w = whisper.New(&whisper.DefaultConfig)

	keyID, err = w.AddSymKeyFromPassword(*argPasswd)
	if err != nil {
		log.Panic("Failed AddSymKeyFromPassword : %s", err)
	}

	symKey, err = w.GetSymKey(keyID)
	if err != nil {
		log.Panic("Failed GetSymKey: %s", err)
	}

	asymKeyID, err = w.NewKeyPair()
	if err != nil {
		log.Panic("Failed to generate a new key pair: %s", err)
	}

	asymKey, err = w.GetPrivateKey(asymKeyID)
	if err != nil {
		log.Panic("Failed to retrieve a new key pair: %s", err)
	}

	/* Install Filter */
	topic = whisper.BytesToTopic([]byte(*argTopic))
	filter := whisper.Filter{
		KeySym: symKey,
		Topics: [][]byte{topic[:]},
	}

	filterID, err = w.Subscribe(&filter)
	if err != nil {
		log.Panic("Failed to install filter: %s", err)
	}
}

func serverConfig() {
	var peers []*discover.Node

	for _, node := range params.MainnetBootnodes {
		peer := discover.MustParseNode(node)
		peers = append(peers, peer)
	}
	for _, node := range ethereumBootnodes {
		peer := discover.MustParseNode(node)
		peers = append(peers, peer)
	}

	srv = &p2p.Server{
		Config: p2p.Config{
			PrivateKey:     asymKey,
			MaxPeers:       100,
			Protocols:      w.Protocols(),
			StaticNodes:    peers,
			BootstrapNodes: peers,
			TrustedNodes:   peers,
		},
	}
}

func peerMonitor() {

	subchan := make(chan *p2p.PeerEvent)
	sub := srv.SubscribeEvents(subchan)
	defer sub.Unsubscribe()

	for {
		select {
		case v := <-subchan:
			if v.Type == p2p.PeerEventTypeAdd {
				fmt.Printf("\nAdd Peer %s\n", v.Peer.String())
			} else if v.Type == p2p.PeerEventTypeDrop {
				fmt.Printf("\nDrop Peer %s\n", v.Peer.String())
			}

			if srv.PeerCount() > 0 {
				hasPeerConnected = true
			} else {
				hasPeerConnected = false
			}
		}
	}
}

func showInfo() {

	fmt.Printf("whisper v5\n")
	fmt.Printf("Topic: %s, PoW : %f\n", *argTopic, *argPoW)
	fmt.Printf("Peers:\n")
	peersInfo := srv.PeersInfo()
	for _, peer := range peersInfo {
		fmt.Printf(" ID %s\n", peer.ID)
	}
}

func txLoop(quit chan struct{}) {
	defer wg.Done()
	for {
		s := readInput()
		if s == "quit()" || s == "exit()" {
			fmt.Println("Program terminated")
			close(quit)
			break
		} else if s == "info()" {
			showInfo()
			continue
		}

		if len(s) == 0 {
			continue
		}

		if hasPeerConnected {
			msgSend([]byte(s))
		}
	}
}

func rxLoop(quit chan struct{}) {

	defer wg.Done()

	f := w.GetFilter(filterID)

	ticker := time.NewTicker(time.Millisecond * 20)

	for {
		select {
		case <-ticker.C:
			/* Retrive envelope from pool */
			mail := f.Retrieve()
			for _, msg := range mail {
				msgDisplay(msg)
			}
		case <-quit:
			return
		}
	}
}

func main() {

	/* Parse command line opt */
	flag.Parse()

	/* Configure whisper */
	whisperConfig()

	/* Configure Server  */
	serverConfig()

	/* Start Server */
	err := srv.Start()
	if err != nil {
		log.Panic("Failed to start Server %s", err)
	}
	defer srv.Stop()
	fmt.Println("Server Start...")

	/* Start Whisper background */
	err = w.Start(srv)
	if err != nil {
		log.Panic("Failed to start Whisper %s", err)
	}
	defer w.Stop()

	fmt.Println("Whisper Start...")

	go peerMonitor()

	var quit = make(chan struct{})

	wg.Add(1)
	go rxLoop(quit)

	wg.Add(1)
	go txLoop(quit)

	wg.Wait()

}

func msgSend(payload []byte) {

	params := whisper.MessageParams{
		Src:      asymKey,
		KeySym:   symKey,
		Payload:  payload,
		Topic:    topic,
		TTL:      whisper.DefaultTTL,
		PoW:      *argPoW,
		WorkTime: 5,
	}

	/* Craete message */
	msg, err := whisper.NewSentMessage(&params)
	if err != nil {
		log.Panic("failed to create new message: %s", err)
	}

	/* Wrap message into envelope */
	envelope, err := msg.Wrap(&params)
	if err != nil {
		fmt.Printf("failed to seal message: %v \n", err)
		return
	}

	/* Send envelope into pool */
	err = w.Send(envelope)
	if err != nil {
		fmt.Printf("failed to send message: %v \n", err)
		return
	}

	return
}

func msgDisplay(msg *whisper.ReceivedMessage) {
	payload := string(msg.Payload)
	timestamp := time.Unix(int64(msg.Sent), 0).Format("2006-01-02 15:04:05")
	var sender common.Address
	if msg.Src != nil {
		sender = crypto.PubkeyToAddress(*msg.Src)
	}

	if whisper.IsPubKeyEqual(msg.Src, &asymKey.PublicKey) {
		fmt.Printf("\n(%s PoW %f): %s\n", timestamp, msg.PoW, payload)
	} else {
		fmt.Printf("\n%x(%s PoW %f): %s\n", sender, timestamp, msg.PoW, payload)
	}
}

func readInput() string {

	if !hasPeerConnected {
		fmt.Printf("Connecting...(may take several minutes)\n")
	} else {
		fmt.Printf(">>")
	}

	f := bufio.NewReader(os.Stdin)

	input, _ := f.ReadString('\n')
	input = strings.TrimRight(input, "\n\r")

	return input
}
