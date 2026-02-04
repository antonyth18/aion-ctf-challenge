package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const AdminToken = "AION_SUPER_SECRET_ADMIN_TOKEN_445"

// MEMORY POOL OPTIMIZATION
var sessionPool = sync.Pool{
	New: func() interface{} {
		return &Session{
			ID:             generateID(),
			CreatedAt:      time.Now(),
			PrivilegeLevel: "GUEST", // Default safe state
			History:        make([]string, 0, 10),
		}
	},
}

type Session struct {
	ID             string
	User           string
	PrivilegeLevel string
	CreatedAt      time.Time
	History        []string
}

// Reset clears the session for reuse.
func (s *Session) Reset() {
	s.ID = generateID()
	s.User = ""
	// Optimization: PrivilegeLevel is overwritten by next user anyway?
	s.CreatedAt = time.Now()
	s.History = s.History[:0]
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessionPool.Get().(*Session)
	
	defer func() {
		sess.Reset()
		sessionPool.Put(sess)
	}()

	authHeader := r.Header.Get("X-AION-Auth")
	if authHeader == AdminToken {
		sess.PrivilegeLevel = "ADMIN"
		sess.User = "Commander"
	} else if authHeader != "" {
		sess.PrivilegeLevel = "GUEST"
		sess.User = "Visitor"
	} else {
		// If no auth header, assume anonymous.
		if sess.User == "" {
			sess.User = "Anonymous"
		}
	}

	cmd := r.URL.Query().Get("cmd")
	if cmd != "" {
		sess.History = append(sess.History, cmd)
		
		if sess.PrivilegeLevel == "ADMIN" {
			output := runDiagnostics(cmd)
			fmt.Fprintf(w, "COMMANDER AUTHORIZED.\nResult: %s\n", output)
			return
		} else {
			fmt.Fprintf(w, "ACCESS DENIED. Current Level: %s. \n", sess.PrivilegeLevel)
			return
		}
	}

	fmt.Fprintf(w, "AION Interface. User: %s [%s]\n", sess.User, sess.PrivilegeLevel)
}

func runDiagnostics(input string) string {
	if strings.Contains(input, "flag") || strings.Contains(input, "cat") {
		return "MALICIOUS INPUT DETECTED"
	}
	
	// Vulnerable echo wrapper
	cmd := exec.Command("sh", "-c", "echo Running: "+input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Background Admin Bot
func adminBot() {
	for {
		time.Sleep(100 * time.Millisecond)
		sess := sessionPool.Get().(*Session)
		sess.PrivilegeLevel = "ADMIN"
		sess.User = "AUTOBOT_COMMANDER"
		time.Sleep(10 * time.Millisecond)
		sess.Reset()
		sessionPool.Put(sess)
	}
}

func main() {
	go adminBot()
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}