package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	shh "github.com/ethereum/go-ethereum/whisper/shhclient"
	whisper6 "github.com/ethereum/go-ethereum/whisper/whisperv6"
)

func main() {
	client, err := shh.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal("connect error:", err)
	}

	privateKeyID, err := client.NewKeyPair(context.Background())
	if err != nil {
		log.Fatal("new key pair error:", err)
	}
	fmt.Printf("private key id: %s\n", privateKeyID)

	filterID, err := client.NewMessageFilter(context.Background(), whisper6.Criteria{PrivateKeyID: privateKeyID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("filter id: %s\n", filterID)

	publicKey, err := client.PublicKey(context.Background(), privateKeyID)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("public key: %s\n", hexutil.Encode(publicKey))

	messageHash, err := client.Post(context.Background(), whisper6.NewMessage{
		TTL:       60,
		PowTime:   2,
		PowTarget: 2.5,
		Payload:   []byte("Hello"),
		PublicKey: publicKey,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("message hash: %s\n", messageHash)

	messages, err := client.FilterMessages(context.Background(), filterID)
	if err != nil {
		log.Fatal(err)
	}
	for _, message := range messages {
		fmt.Printf("message: %s", string(message.Payload))
	}
}
