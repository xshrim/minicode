package main

import "crypto/rand"
import "encoding/hex"
import "fmt"
import "log"


/*
 echo -e "/key/swarm/psk/1.0.0/\n/base16/\n `tr -dc 'a-f0–9' < /dev/urandom
 head -c64`" > ~/.ipfs/swarm.key
*/
func main() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalln("While trying to read random source:", err)
	}

	fmt.Println("/key/swarm/psk/1.0.0/")
	fmt.Println("/base16/")
	fmt.Print(hex.EncodeToString(key))
} 
