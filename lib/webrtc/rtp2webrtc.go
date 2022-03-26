package main

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/hktalent/gohktools/lib/internal/signal"
	stlist "github.com/hktalent/gohktools/lib/utils"
	"github.com/pion/webrtc/v3"
)

/*
https://github.com/pion/webrtc/tree/master/examples/rtp-to-webrtc
RTP: Real-time Transport Protocol
https://en.wikipedia.org/wiki/Real-time_Transport_Protocol
实时传输协议 (RTP) 是一种通过 IP 网络提供音频和视频的网络协议。RTP用于涉及流媒体的通信和娱乐系统，如电话、视频电话会议应用程序（包括WebRTC）、电视服务和基于网络的推送通话功能。
RTP通常通过用户数据报协议（UDP）运行。RTP与RTP控制协议（RTCP）一起使用。虽然RTP承载媒体流（例如音频和视频），但RTCP用于监控传输统计和服务质量（QoS），并帮助多个流的同步。
RTP是IP语音的技术基础之一，在这种情况下，通常与信令协议（如会话启动协议（SIP）一起使用，该协议在网络上建立连接。
RTP由互联网工程工作组（IETF）的音频视频传输工作组开发，并于1996年首次作为RFC 1889发布，然后于2003年被RFC 3550取代。
*/
func main() {
	// aSt := append(stlist.StunList{}.GetStunList(), "stun:stun.l.google.com:19302")
	aSt := stlist.StunList{}.GetStunList()
	fmt.Println(aSt)
	// aSt := []string{"stun:stun.l.google.com:19302"}
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: aSt,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	// Open a UDP Listener for RTP Packets on port 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 5004})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "51pwn4RTP")
	if err != nil {
		panic(err)
	}
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateFailed {
			if closeErr := peerConnection.Close(); closeErr != nil {
				panic(closeErr)
			}
		}
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	signal.Decode(signal.MustReadStdin(), &offer)

	// Set the remote SessionDescription
	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))

	// Read RTP packets forever and send them to the WebRTC Client
	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("error during read: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				// The peerConnection has been closed.
				return
			}

			panic(err)
		}
	}
}
