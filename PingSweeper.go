package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: pingsweeper <subnet>")
		fmt.Println("Example: pingsweeper 192.168.1.0/24")
		return
	}

	subnet := os.Args[1]
	activeHosts := make(chan string)
	var wg sync.WaitGroup

	// Parse the subnet
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		fmt.Printf("Error parsing subnet: %v\n", err)
		return
	}

	// Iterate through all IPs in the subnet
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			if ping(ip) {
				activeHosts <- ip
			}
		}(ip.String())
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(activeHosts)
	}()

	// Print active hosts
	fmt.Println("Active hosts:")
	for host := range activeHosts {
		fmt.Println(host)
	}
}

// ping sends an ICMP echo request to the given IP and returns true if a reply is received
func ping(ip string) bool {
	// Create a raw socket for ICMP
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Printf("Error creating ICMP socket: %v\n", err)
		return false
	}
	defer conn.Close()

	// Create an ICMP echo message
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("PINGSWEEPER"),
		},
	}

	// Marshal the message into bytes
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		fmt.Printf("Error marshaling ICMP message: %v\n", err)
		return false
	}

	// Send the ICMP echo request
	start := time.Now()
	if _, err := conn.WriteTo(msgBytes, &net.IPAddr{IP: net.ParseIP(ip)}); err != nil {
		return false
	}

	// Set a timeout for the reply
	reply := make([]byte, 1500)
	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return false
	}

	// Wait for a reply
	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		return false
	}

	// Parse the reply
	duration := time.Since(start)
	parsedReply, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		return false
	}

	// Check if the reply is an echo reply
	switch parsedReply.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("%s is active (response time: %v)\n", ip, duration)
		return true
	default:
		return false
	}
}

// inc increments an IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
