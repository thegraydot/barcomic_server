package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/thegraydot/barcomic/internal"
	"github.com/spf13/cobra"
)

var Version = "dev"
var Hash = "none"

var (
	addr              string
	port              string
	enableHttps       bool
	disableKeystrokes bool
	interactive       bool
	verbose           bool
)

var rootCmd = &cobra.Command{
	Use:   "barcomic",
	Short: "An HTTP API for receiving comic book barcodes",
	RunE:  runServer,
}

func init() {
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "", "Address to listen on")
	rootCmd.Flags().StringVarP(&port, "port", "p", "9999", "Port to listen on")
	rootCmd.Flags().BoolVarP(&enableHttps, "https", "k", false, "Enable HTTPS (self-signed)")
	rootCmd.Flags().BoolVarP(&disableKeystrokes, "no-keystrokes", "s", false, "Disable keystroke injection")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "Run interactive configuration")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose information")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	fmt.Printf("[*] barcomic %s-%s\n", Version, Hash)

	// If address is provided, set interactive to false
	if addr != "" {
		interactive = false
	}

	// If requested, run the interactive server config
	if interactive {
		addr = interactiveNetworkConfiguration()
	}

	// Validate IP address and port before starting server
	if !validateAddr(addr) {
		fmt.Printf("[*] Error: Invalid IP address %s\n", addr)
		os.Exit(1)
	}
	if !validatePort(port) {
		fmt.Printf("[*] Error: Invalid port %s\n", port)
		os.Exit(1)
	}

	barcomic.Start(addr, port, enableHttps, disableKeystrokes, verbose)
	return nil
}

func interactiveNetworkConfiguration() string {
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("[*] Error: %+v\n", err.Error()))
		return "0.0.0.0"
	}

	var availableInterfaces [][2]string

	// Loop through interfaces
	for _, i := range interfaces {
		// Get interface name (e.g., wlan0)
		byNameInterface, err := net.InterfaceByName(i.Name)
		if err != nil {
			fmt.Println(err)
		}

		// Get addresses on interface and loop
		addresses, err := byNameInterface.Addrs()
		if err != nil {
			fmt.Println(err)
		}
		for _, address := range addresses {
			// Check address is IP4 and not loopback (127.0.0.1)
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					var availableInterface [2]string
					availableInterface[0] = ipnet.IP.String()
					availableInterface[1] = i.Name
					availableInterfaces = append(availableInterfaces, availableInterface)
				}
			}
		}
	}

	fmt.Printf("[*] The following addresses are available...\n")
	for i, availableInterface := range availableInterfaces {
		fmt.Printf("    [%d] %s (%s)\n", i, availableInterface[0], availableInterface[1])
	}

	// Get address selection from user
	fmt.Print("[*] Enter IP address [0.0.0.0]: ")
	var addrInput string
	fmt.Scanln(&addrInput)
	addrInput = strings.Trim(addrInput, " ")

	// Set IP to any address - default action
	if addrInput == "" {
		addrInput = "0.0.0.0"
	}

	// Check for valid input address
	if validateAddr(addrInput) {
		return addrInput
	}

	addrInputInt, err := strconv.Atoi(addrInput)
	if err != nil {
		fmt.Println("\n[*] Error. Could not determine network interface selection.")
		os.Exit(1)
	}
	if addrInputInt < 0 || addrInputInt >= len(availableInterfaces) {
		fmt.Println("\n[*] Error: Network interface selection error.")
		os.Exit(1)
	}
	return availableInterfaces[addrInputInt][0]
}

func validateAddr(a string) bool {
	if a == "" {
		return false
	}
	// Fast path: valid IP address (IPv4 or IPv6)
	if net.ParseIP(a) != nil {
		return true
	}
	// Fallback: attempt DNS resolution for hostnames
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupHost(ctx, a)
	return err == nil && len(addrs) > 0
}

func validatePort(p string) bool {
	portInt, err := strconv.ParseInt(p, 10, 0)
	if err != nil {
		return false
	}
	return portInt >= 0 && portInt <= 65535
}
