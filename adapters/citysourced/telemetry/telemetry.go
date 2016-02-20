package telemetry

import (
	"net"
	"time"

	"Gateway311/engine/logs"
)

var (
	log         = logs.Log
	chTQue      chan msgSender
	monitorAddr = "127.0.0.1:5051"
)

// SendRPC queues an RPC status message onto the send channel.
func SendRPC(id, status, route, url string, at time.Time) {
	statusMsg := AdpRPCMsgType{
		ID:     id,
		Status: status,
		Route:  route,
		URL:    url,
		At:     at,
	}
	chTQue <- msgSender(statusMsg)

}

// Shutdown should be called to gracefully stop the telemetry processes.
func Shutdown() {
	close(chTQue)
}

func init() {
	chTQue = make(chan msgSender, 100)

	tlmtryServer, err := net.ResolveUDPAddr("udp", monitorAddr)
	if err != nil {
		log.Errorf("Cannot start telemetry - %s", err.Error())
		return
	}

	conn, err := net.DialUDP("udp", nil, tlmtryServer)
	if err != nil {
		log.Errorf("Cannot start telemetry - %s", err.Error())
		return
	}

	go func() {
		log.Debug("Telemetry sender starting...")
		defer conn.Close()
		for m := range chTQue {
			msg, err := m.Marshal()
			if err != nil {
				log.Warning("unable to send message - %s", err.Error())
				continue
			}
			log.Debug(string(msg))
			if _, err := conn.Write(msg); err != nil {
				log.Warning(err.Error())
			}
		}
	}()
}
