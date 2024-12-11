package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

func checkAllSet(values ...string) bool {
	for _, v := range values {
		if v == "" {
			return false
		}
	}
	return true
}

func main() {

	host := flag.String("host", "", "SMPP server host")
	port := flag.String("port", "", "SMPP server port")
	user := flag.String("user", "", "SMPP system ID")
	password := flag.String("password", "", "SMPP password")
	source := flag.String("source", "", "Sender's number / source address")
	destination := flag.String("destination", "", "Recipient's number")
	message := flag.String("message", "", "Message to send")

	flag.Parse()

	if !checkAllSet(*host, *port, *user, *password, *source, *destination, *message) {
		fmt.Println("Error: All 'host', 'port', 'user', 'password', 'source', 'destination', and 'message' must be provided.")
		os.Exit(1)
	}

	tx := smpp.Transmitter{
		Addr:   fmt.Sprintf("%s:%s", *host, *port),
		User:   *user,
		Passwd: *password,
	}

	defer tx.Close()
	conn := <-tx.Bind()
	switch conn.Status() {
	case smpp.Connected:
	default:
		fmt.Printf("Failed to connect: %s\n", conn.Error())
	}

	sm, err := tx.Submit(&smpp.ShortMessage{
		Src:      *source,
		Dst:      *destination,
		Text:     pdutext.Raw(*message),
		Register: pdufield.NoDeliveryReceipt,
	})

	if err != nil {
		fmt.Printf("Failed to send message: %s\n", err)
	}

	fmt.Printf("Message sent: %s\n", sm.Text.Decode())

}
