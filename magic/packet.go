package magic

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"regexp"
)

// Segment is 6 bytes
type Segment [6]byte

// MacAddress is 16 Segments
type MacAddress [16]Segment

// MagicPacket is constituted of 6 bytes of 0xFF followed by
// 16 groups of the destination MAC address.
type MagicPacket struct {
	header  Segment
	payload MacAddress
}

// Define globals for MacAddress parsing
var (
	delims = ":-"
	reMAC  = regexp.MustCompile(`^([0-9a-fA-F]{2}[` + delims + `]){5}([0-9a-fA-F]{2})$`)
)

// NewMagicPacket : This function accepts a MAC Address string, and returns a pointer to
// a MagicPacket object. A Magic Packet is a broadcast frame which
// contains 6 bytes of 0xFF followed by 16 repetitions of a given mac address.
func NewMagicPacket(mac string) (*MagicPacket, error) {
	var packet MagicPacket
	var macAddr Segment

	// We only support 6 byte MAC addresses
	if !reMAC.MatchString(mac) {
		return nil, errors.New("MAC address " + mac + " is not valid.")
	}

	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	// Copy bytes from the returned HardwareAddr -> a fixed size MACAddress
	for idx := range macAddr {
		macAddr[idx] = hwAddr[idx]
	}

	// Setup the header which is 6 repetitions of 0xFF
	for idx := range packet.header {
		packet.header[idx] = 0xFF
	}

	// Setup the payload which is 16 repetitions of the MAC addr
	for idx := range packet.payload {
		packet.payload[idx] = macAddr
	}

	return &packet, nil
}

// SendMagicPacket This function accepts a MAC address string,
// and sends a magic packet to the given MAC address
func SendMagicPacket(macAddr string) error {
	magicPacket, err := NewMagicPacket(macAddr)
	if err != nil {
		return err
	}

	// Fill our byte buffer with the bytes in our MagicPacket
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, magicPacket)

	// Get a UDPAddr to send the broadcast to
	udpAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		fmt.Printf("Unable to get a UDP address for %s:%d\n", udpAddr.IP.String(), udpAddr.Port)
		return err
	}

	// Open a UDP connection, and defer its cleanup
	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("Unable to dial UDP address for %s\n", "255.255.255.255:9")
		return err
	}
	defer connection.Close()

	// Write the bytes of the MagicPacket to the connection
	bytesWritten, err := connection.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Unable to write packet to connection\n")
		return err
	} else if bytesWritten != 102 {
		fmt.Printf("Warning: %d bytes written, %d expected!\n", bytesWritten, 102)
	}

	return nil
}
