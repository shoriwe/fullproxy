package ArgumentParser

import (
	"fmt"
	"os"
)

func ShowGeneralHelpMessage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage:\n\t", os.Args[0], "PROTOCOL|TOOL *FLAGS\n\nProtocols available:\n\t - socks5\n\t - http\n\t - local-forward\n\t - remote-forward\n\t - master\n\t - translate")
}

func ShowTranslateHelpMessage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage:\n\t", os.Args[0], "translate TARGET *FLAGS\n\nTARGETS available:\n\t - port_forward-socks5\n\t")
}
