package utils

import (
	"net"
	"netdash/internal/logger"
	"os"
)

func LogAccessInfo(port string) {
	hostEnv := os.Getenv("APP_HOST")
	if hostEnv != "" {
		logger.Log("SERVER", "Listening on http://%s%s", hostEnv, port)
		return
	}

	outboundIP, err := getOutboundIP()
	if err == nil {
		logger.Log("SERVER", "Listening on http://%s%s", outboundIP.String(), port)
	} else {
		logger.Log("SERVER", "Or access via localhost: http://127.0.0.1%s", port)
	}
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}