package barcomic

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mdp/qrterminal/v3"
)

type Config struct {
	addr              string
	port              string
	enableHttps       bool
	disableKeystrokes bool
	verbose           bool
}

type Server struct {
	config      Config
	mux         *http.ServeMux
	typeBarcode func(string) error
}

func NewServer(cfg Config) *Server {
	s := &Server{
		config:      cfg,
		mux:         http.NewServeMux(),
		typeBarcode: TypeBarcode,
	}
	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.HandleFunc("/barcode", s.barcodeHandler)
	s.mux.HandleFunc("/", s.otherHandler)
	return s
}

func Start(addr, port string, enableHttps, disableKeystrokes, verbose bool) {
	cfg := Config{
		addr:              addr,
		port:              port,
		enableHttps:       enableHttps,
		disableKeystrokes: disableKeystrokes,
		verbose:           verbose,
	}
	srv := NewServer(cfg)
	printQRCode(srv.config.addr, srv.config.port)
	srv.startRestAPI()
}

func printQRCode(addr, port string) {
	qrconfig := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	host := addr + ":" + port
	qrterminal.GenerateWithConfig(host, qrconfig)
	fmt.Printf("[*] Starting server using %s:%s\n", addr, port)
}
