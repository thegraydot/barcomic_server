package barcomic

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var responseOK = []byte("OK")
var responseError = []byte("ERROR")

// barcodeRegexp matches strings that are entirely 12-17 digits
// Compiled once at package init. Anchors prevent partial matches
var barcodeRegexp = regexp.MustCompile(`^\d{12,17}$`)

func (s *Server) startRestAPI() {
	fmt.Println("[*] Generating TLS certificate...")
	if s.config.enableHttps {
		tlsCert := GenerateTLSCertificate(s.config.addr, s.config.port)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		}

		server := http.Server{
			Addr:         s.config.addr + ":" + s.config.port,
			Handler:      s.mux,
			TLSConfig:    tlsConfig,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Minute,
		}

		fmt.Printf("[*] Starting HTTPS server...\n\n")
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("ERROR: %v", err)
		}
	} else {
		server := http.Server{
			Addr:         s.config.addr + ":" + s.config.port,
			Handler:      s.mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Minute,
		}

		fmt.Printf("[*] Starting HTTP server...\n\n")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("ERROR: %v", err)
		}
	}
}

func (s *Server) verboseLoggingHandler(req *http.Request) {
	if s.config.verbose {
		fmt.Printf("INFO: %s %s %s\n", req.Method, req.RemoteAddr, req.RequestURI)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, req *http.Request) {
	s.verboseLoggingHandler(req)
	if req.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write(responseOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseError)
		return
	}
}

func (s *Server) barcodeHandler(w http.ResponseWriter, req *http.Request) {
	s.verboseLoggingHandler(req)
	if req.Method == "POST" {
		// UPC + EAN5 is longest barcode with 17 chars
		// So set max byte length to 20
		req.Body = http.MaxBytesReader(w, req.Body, 20)
		buffer, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseError)
			return
		}

		bufferString := string(buffer)

		// Check request body is digits with anchored match
		if !barcodeRegexp.MatchString(bufferString) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseError)
			return
		}

		// Check request body is a valid UPC
		isValidUpc, err := validateUpc(bufferString)
		if err != nil || !isValidUpc {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseError)
			return
		}

		if !s.config.disableKeystrokes {
			if err := s.typeBarcode(bufferString); err != nil {
				fmt.Printf("WARN: keystroke injection failed: %v\n", err)
			}
		}

		// Print barcode if verbose
		if s.config.verbose {
			fmt.Printf("INFO: %s\n", bufferString)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(buffer)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseError)
		return
	}
}

func (s *Server) otherHandler(w http.ResponseWriter, req *http.Request) {
	s.verboseLoggingHandler(req)
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
}

func validateUpc(barcode string) (bool, error) {
	// Extract first 11 digits of barcode
	barcodePrefix := barcode[0:11]
	// Extract last (check) digit from barcode
	checkDigit := barcode[11:12]

	// Sum all digits
	// Even digits are multiplied by 3
	sum := 0
	for i, v := range barcodePrefix {
		value, err := strconv.Atoi(string(v))
		if err != nil {
			return false, fmt.Errorf("validateUpc: non-digit at position %d: %w", i, err)
		}
		if i%2 == 0 {
			sum += 3 * value
		} else {
			sum += value
		}
	}

	result := (10 - sum%10) % 10
	checkDigitInt, err := strconv.Atoi(string(checkDigit))
	if err != nil {
		return false, fmt.Errorf("validateUpc: non-digit check character: %w", err)
	}
	return result == checkDigitInt, nil
}
