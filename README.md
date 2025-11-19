# PingSweeper

PingSweeper is a lightweight subnet scanning tool written in Go.  
It scans an entire subnet by sending ICMP echo requests (pings) to every host and reports which ones respond.  
This project was built to practice low-level networking, raw ICMP sockets, and Go’s concurrency model.

---

## Features

- **Full subnet sweep**  
  Scans every host in a CIDR range such as `192.168.1.0/24`.

- **Raw ICMP sockets**  
  Builds and sends ICMP packets manually.

- **Concurrent scanning**  
  Uses goroutines for fast parallel scanning.

- **Latency measurement**  
  Shows response times for active hosts.

- **Zero external dependencies**  
  Uses Go's standard library + `golang.org/x/net/icmp`.

---

## Tech Used

- **Go (Golang)**  
- **icmp / ipv4 packages** for raw packet handling  
- **net** for IP manipulation  
- **sync.WaitGroup** for concurrency  

---

## How It Works

1. Accepts a subnet in CIDR format (e.g., `10.0.0.0/24`)  
2. Iterates through all possible host IPs using bitwise logic  
3. For each IP:
   - Constructs an ICMP Echo Request packet  
   - Sends it through a raw ICMP socket  
   - Waits for a reply (with timeout)  
4. Prints all responding hosts along with response times  
5. Concurrency makes the scan significantly faster  

