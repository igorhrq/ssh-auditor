package main

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func hashKey(key ssh.PublicKey) string {
	hash := sha256.Sum256(key.Marshal())
	fp := base64.RawStdEncoding.EncodeToString(hash[:])
	return fp
}

func FetchSSHKeyFingerPrint(hostport string) string {

	var keyFingerprint string

	DumpHostkey := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		fp := hashKey(key)
		keyFingerprint = fp
		return nil
	}

	config := &ssh.ClientConfig{
		User: "security",
		Auth: []ssh.AuthMethod{
			ssh.Password("security"),
		},
		HostKeyCallback: DumpHostkey,
		Timeout:         2 * time.Second,
	}

	client, err := ssh.Dial("tcp", hostport, config)
	if err == nil {
		//This was supposed to fail
		client.Close()
		log.Printf("BADPW %s (%s): user=security password=security", hostport)
	}
	return keyFingerprint
}