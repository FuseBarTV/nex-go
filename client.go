package nex

import (
	"crypto/rc4"
	"net"
	"time"
)

// Client represents a connected or non-connected PRUDP client
type Client struct {
	address                   *net.UDPAddr
	server                    *Server
	cipher                    *rc4.Cipher
	decipher                  *rc4.Cipher
	signatureKey              []byte
	signatureBase             int
	secureKey                 []byte
	serverConnectionSignature []byte
	clientConnectionSignature []byte
	sessionID                 int
	sessionKey                []byte
	sequenceIDIn              *Counter
	sequenceIDOut             *Counter
	pid                       uint32
	localStationUrl           string
	connectionId              uint32
	pingTimeoutTime           time.Time
	connected                 bool
}

// Reset resets the Client to default values
func (client *Client) Reset() {
	client.sequenceIDIn = NewCounter(0)
	client.sequenceIDOut = NewCounter(0)

	client.UpdateAccessKey(client.Server().AccessKey())
	client.UpdateRC4Key([]byte("CD&ML"))

	if client.Server().PrudpVersion() == 0 {
		client.SetServerConnectionSignature(make([]byte, 4))
		client.SetClientConnectionSignature(make([]byte, 4))
	} else {
		client.SetServerConnectionSignature([]byte{})
		client.SetClientConnectionSignature([]byte{})
	}

	client.SetConnected(false)
}

// Address returns the clients UDP address
func (client *Client) Address() *net.UDPAddr {
	return client.address
}

// Server returns the server the client is currently connected to
func (client *Client) Server() *Server {
	return client.server
}

// UpdateRC4Key sets the client RC4 stream key
func (client *Client) UpdateRC4Key(RC4Key []byte) {
	cipher, _ := rc4.NewCipher(RC4Key)
	client.cipher = cipher

	decipher, _ := rc4.NewCipher(RC4Key)
	client.decipher = decipher
}

// Cipher returns the RC4 cipher stream for out-bound packets
func (client *Client) Cipher() *rc4.Cipher {
	return client.cipher
}

// Decipher returns the RC4 cipher stream for in-bound packets
func (client *Client) Decipher() *rc4.Cipher {
	return client.decipher
}

// UpdateAccessKey sets the client signature base and signature key
func (client *Client) UpdateAccessKey(accessKey string) {
	client.signatureBase = sum([]byte(accessKey))
	client.signatureKey = MD5Hash([]byte(accessKey))
}

// SignatureBase returns the v0 checksum signature base
func (client *Client) SignatureBase() int {
	return client.signatureBase
}

// SignatureKey returns signature key
func (client *Client) SignatureKey() []byte {
	return client.signatureKey
}

// SetServerConnectionSignature sets the clients server-side connection signature
func (client *Client) SetServerConnectionSignature(serverConnectionSignature []byte) {
	client.serverConnectionSignature = serverConnectionSignature
}

// ServerConnectionSignature returns the clients server-side connection signature
func (client *Client) ServerConnectionSignature() []byte {
	return client.serverConnectionSignature
}

// SetClientConnectionSignature sets the clients client-side connection signature
func (client *Client) SetClientConnectionSignature(clientConnectionSignature []byte) {
	client.clientConnectionSignature = clientConnectionSignature
}

// ClientConnectionSignature returns the clients client-side connection signature
func (client *Client) ClientConnectionSignature() []byte {
	return client.clientConnectionSignature
}

// SequenceIDCounterOut returns the clients packet SequenceID counter for out-going packets
func (client *Client) SequenceIDCounterOut() *Counter {
	return client.sequenceIDOut
}

// SequenceIDCounterIn returns the clients packet SequenceID counter for incoming packets
func (client *Client) SequenceIDCounterIn() *Counter {
	return client.sequenceIDIn
}

// SetSessionKey sets the clients session key
func (client *Client) SetSessionKey(sessionKey []byte) {
	client.sessionKey = sessionKey
}

// SessionKey returns the clients session key
func (client *Client) SessionKey() []byte {
	return client.sessionKey
}

// SetPID sets the clients NEX PID
func (client *Client) SetPID(pid uint32) {
	client.pid = pid
}

// PID returns the clients NEX PID
func (client *Client) PID() uint32 {
	return client.pid
}

// SetLocalStationUrl sets the clients Local Station URL
func (client *Client) SetLocalStationUrl(localStationUrl string) {
	client.localStationUrl = localStationUrl
}

// LocalStationUrl returns the clients Local Station URL
func (client *Client) LocalStationUrl() string {
	return client.localStationUrl
}

// SetConnectionId sets the clients Connection ID
func (client *Client) SetConnectionId(connectionId uint32) {
	client.connectionId = connectionId
}

// ConnectionId returns the clients Connection ID
func (client *Client) ConnectionId() uint32 {
	return client.connectionId
}

// SetConnected sets the clients connection status
func (client *Client) SetConnected(connected bool) {
	client.connected = connected
}

// IncreasePingTimeoutTime adds a number of seconds to the timeout timer
func (client *Client) IncreasePingTimeoutTime(seconds int) {
	client.pingTimeoutTime = time.Now().Add(time.Second * time.Duration(seconds))
}

// StartTimeoutTimer begins the packet timeout timer
func (client *Client) StartTimeoutTimer() {
	for client.connected {
		if time.Now().After(client.pingTimeoutTime) {
			client.SetConnected(false)
			client.server.Kick(client)
		}
	}
}

// NewClient returns a new PRUDP client
func NewClient(address *net.UDPAddr, server *Server) *Client {
	client := &Client{
		address: address,
		server:  server,
	}

	client.Reset()

	return client
}
