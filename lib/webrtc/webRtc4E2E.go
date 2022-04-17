package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	stlist "github.com/hktalent/gohktools/lib/utils"
	"github.com/pion/webrtc/v3"
)

func signalCandidate(addr string, c *webrtc.ICECandidate) error {
	// payload := []byte(c.ToJSON().Candidate)
	// fmt.Println("c.ToJSON ", c.ToJSON())
	// resp, err := http.Post(fmt.Sprintf("http://%s/candidate", addr), "application/json; charset=utf-8", bytes.NewReader(payload)) //nolint:noctx
	// if err != nil {
	// 	return err
	// }

	// if closeErr := resp.Body.Close(); closeErr != nil {
	// 	return closeErr
	// }

	return nil
}

const (
	rtcpPLIInterval = time.Second * 3
)

func getPeerConnection() (*webrtc.PeerConnection, webrtc.Configuration) {
	aSt := stlist.StunList{}.GetStunList()
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: aSt,
			},
		},
	}
	// fmt.Println(aSt)
	// aSt := []string{"stun:stun.l.google.com:19302"}
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		panic(err)
	}

	return peerConnection, peerConnectionConfig
}

/*
https://github.com/pion/webrtc/blob/master/examples/ice-restart/main.go
https://github.com/pion/webrtc/tree/master/examples/rtp-to-webrtc
RTP: Real-time Transport Protocol
https://en.wikipedia.org/wiki/Real-time_Transport_Protocol
实时传输协议 (RTP) 是一种通过 IP 网络提供音频和视频的网络协议。RTP用于涉及流媒体的通信和娱乐系统，如电话、视频电话会议应用程序（包括WebRTC）、电视服务和基于网络的推送通话功能。
RTP通常通过用户数据报协议（UDP）运行。RTP与RTP控制协议（RTCP）一起使用。虽然RTP承载媒体流（例如音频和视频），但RTCP用于监控传输统计和服务质量（QoS），并帮助多个流的同步。
RTP是IP语音的技术基础之一，在这种情况下，通常与信令协议（如会话启动协议（SIP）一起使用，该协议在网络上建立连接。
RTP由互联网工程工作组（IETF）的音频视频传输工作组开发，并于1996年首次作为RFC 1889发布，然后于2003年被RFC 3550取代。
*/
func main() {
	peerConnection, _ := getPeerConnection()
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()
	// Allow us to receive 1 video track
	if _, err := peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}
	// localTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	// // Set a handler for when a new remote track starts, this just distributes all our packets
	// // to connected peers
	// peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	// 	// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
	// 	// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
	// 	go func() {
	// 		ticker := time.NewTicker(rtcpPLIInterval)
	// 		for range ticker.C {
	// 			if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remoteTrack.SSRC())}}); rtcpSendErr != nil {
	// 				fmt.Println(rtcpSendErr)
	// 			}
	// 		}
	// 	}()

	// 	// Create a local track, all our SFU clients will be fed via this track
	// 	localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "video", "pion")
	// 	if newTrackErr != nil {
	// 		panic(newTrackErr)
	// 	}
	// 	localTrackChan <- localTrack

	// 	rtpBuf := make([]byte, 1400)
	// 	for {
	// 		i, _, readErr := remoteTrack.Read(rtpBuf)
	// 		if readErr != nil {
	// 			panic(readErr)
	// 		}

	// 		// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
	// 		if _, err := localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
	// 			panic(err)
	// 		}
	// 	}
	// })
	// // offer := webrtc.SessionDescription{}
	// // // Set the remote SessionDescription
	// // err = peerConnection.SetRemoteDescription(offer)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // Create answer
	// // answer, err := peerConnection.CreateAnswer(nil)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // // Create channel that is blocked until ICE Gathering is complete
	// // gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	// // // Sets the LocalDescription, and starts our UDP listeners
	// // err = peerConnection.SetLocalDescription(answer)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // // Block until ICE Gathering is complete, disabling trickle ICE
	// // // we do this because we only can exchange one signaling message
	// // // in a production application you should exchange ICE Candidates via OnICECandidate
	// // <-gatherComplete
	// // Get the LocalDescription and take it to base64 so we can paste in browser
	// // fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
	// localTrack := <-localTrackChan
	// for {
	// 	fmt.Println("")
	// 	fmt.Println("Curl an base64 SDP to start sendonly peer connection")

	// 	recvOnlyOffer := webrtc.SessionDescription{}
	// 	// fmt.Println(recvOnlyOffer)
	// 	// signal.Decode(<-sdpChan, &recvOnlyOffer)

	// 	// Create a new PeerConnection
	// 	peerConnection, _ := getPeerConnection()
	// 	rtpSender, err := peerConnection.AddTrack(localTrack)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Read incoming RTCP packets
	// 	// Before these packets are returned they are processed by interceptors. For things
	// 	// like NACK this needs to be called.
	// 	go func() {
	// 		rtcpBuf := make([]byte, 1500)
	// 		for {
	// 			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
	// 				return
	// 			}
	// 		}
	// 	}()

	// 	// Set the remote SessionDescription
	// 	err = peerConnection.SetRemoteDescription(recvOnlyOffer)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Create answer
	// 	answer, err := peerConnection.CreateAnswer(nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Create channel that is blocked until ICE Gathering is complete
	// 	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// 	// Sets the LocalDescription, and starts our UDP listeners
	// 	err = peerConnection.SetLocalDescription(answer)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Block until ICE Gathering is complete, disabling trickle ICE
	// 	// we do this because we only can exchange one signaling message
	// 	// in a production application you should exchange ICE Candidates via OnICECandidate
	// 	<-gatherComplete

	// 	// Get the LocalDescription and take it to base64 so we can paste in browser
	// 	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
	// }
	///////////////////////////////////////////////////////

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})
	// Send the current time via a DataChannel to the remote peer every 3 seconds
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			for range time.Tick(time.Second * 3) {
				if err := d.SendText(time.Now().String()); err != nil {
					panic(err)
				}
			}
		})
	})
	/*start////////////////////////////
	var offer webrtc.SessionDescription
	if err = json.NewDecoder(strings.NewReader("{}")).Decode(&offer); err != nil {
		panic(err)
	}

	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}
	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	//  创建响应者
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	} else if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}
	<-gatherComplete
	////////////////////////////end/////*/

	// fmt.Println("wait gatherComplete ")
	// gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	// <-gatherComplete
	// fmt.Println("gatherComplete is ok")
	// 本地地址
	if nil != peerConnection.LocalDescription() {
		response, err := json.Marshal(*peerConnection.LocalDescription())
		if err != nil {
			panic(err)
		}
		fmt.Println(response)
	}

	var candidatesMux sync.Mutex
	// 存储E2E，P2P的地址、端口信息
	pendingCandidates := make([]*webrtc.ICECandidate, 0)
	answerAddr := "0.0.0.0:0"
	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		// 这里可以获得本地的udp 的外网、内网ip和端口信息
		// fmt.Println("OnICECandidate c.ToJSON ", c.ToJSON())
		// 远程节点来的时候就加入节点
		// 这里得到P2P UDP打洞的ip、端口信息
		if desc == nil {
			// xxj, err := json.Marshal(c.ToJSON().Candidate)
			// if nil == err {
			// 	fmt.Println("OnICECandidate c.ToJSON ", string(xxj))
			// }
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := signalCandidate(answerAddr, c); onICECandidateErr != nil {
			panic(onICECandidateErr)
		}
	})
	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("51pwn4E2E_P2P_hacker", nil)
	if err != nil {
		panic(err)
	}
	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			// os.Exit(0)
		}
	})
	// Register channel opening handling
	// dataChannel.OnOpen(func() {
	// 	fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

	// 	for range time.NewTicker(5 * time.Second).C {
	// 		message := signal.RandSeq(15)
	// 		fmt.Printf("Sending '%s'\n", message)

	// 		// Send the message as text
	// 		sendTextErr := dataChannel.SendText(message)
	// 		if sendTextErr != nil {
	// 			panic(sendTextErr)
	// 		}
	// 	}
	// })

	// // Register text message handling
	// dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
	// 	fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	// })
	// Create an offer to send to the other process
	offer1, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	if err = peerConnection.SetLocalDescription(offer1); err != nil {
		panic(err)
	}
	// Send our offer to the HTTP server listening in the other process
	// payload, err := json.Marshal(offer1)
	// if err != nil {
	// 	panic(err)
	// }
	// 这里发送 payload给E2E列表
	fmt.Println("payload ==> ", offer1)

	// Block forever
	select {}
}
