package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

///////////////////////////

type LogWriter struct {
}

func NewLogWriter() *LogWriter {
    lw := &LogWriter{}
    return lw
}

func (lw LogWriter) Write (p []byte) (n int, err error) {
    log.Print(string(p))
    return len(p), nil
}

///////////////////////////

type Repository struct {
	FullName string `json:"full_name"`
}

type GithubJson struct {
	Repository Repository
	Ref        string
	After      string
	Deleted    bool
}

type Config struct {
	Hooks []Hook
}

type Hook struct {
	Repo   string
	Branch string
	Shell  string
}

var config Config

func loadConfig(configFile *string) {
	configData, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal(err)
	}

	addHandler()
}

func setLog(logFile *string) {
	if *logFile != "" {
		var (
			log_handler *os.File
			err error
		)
		if *logFile == "-" {
			log_handler = os.Stdout
		} else {
			log_handler, err = os.OpenFile(*logFile, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0777)
			if err != nil {
				panic("cannot write log")
			}
		}
		log.SetOutput(log_handler)
	}
	log.SetFlags(5)
}

func startWebserver() {
	if *scriptsPath != "" {
		log.Printf("Looking in %s for hook scripts", *scriptsPath)
	} else {
		log.Printf("Looking on $PATH for hook scripts")
	}
	log.Printf("Starting gohub on 0.0.0.0:%s", *port)
	http.ListenAndServe(":"+*port, nil)
}

func checkSignature(body []byte, r *http.Request) bool {
	header := r.Header.Get("X-Hub-Signature")
	if header == "" {
		log.Printf("No signature header")
		return false
	}
	splitHeader := strings.SplitN(header, "=", 2)
	algo, signature := splitHeader[0], splitHeader[1]
	if algo != "sha1" {
		log.Printf("Not sha1")
	}
	key := []byte(*githubSecret)
	mac := hmac.New(sha1.New, key)
	mac.Write(body)
	hash := mac.Sum(nil)
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		log.Printf("Bad hex value in signature")
		return false
	}
	checkResult := hmac.Equal(signatureBytes, hash)
	if !checkResult {
		log.Printf("Signature check failed %s != %s", hex.EncodeToString(signatureBytes), hex.EncodeToString(hash))
	}
	return checkResult
}

func addHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Body.Close()
		success := false
		defer func() {
			if success == false {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{\"status\":\"ERROR\"}"))
			}
		}()
		w.Header().Set("Content-Type", "application/json")

		if !checkSignature(body, r) {
			return
		}

		decoder := json.NewDecoder(bytes.NewBuffer(body))
		var data GithubJson
		err = decoder.Decode(&data)

		if err != nil {
			log.Println(err)
			return
		}

		var hook Hook
		for _, cfgHook := range config.Hooks {
	 		if cfgHook.Repo == data.Repository.FullName && data.Ref == "refs/heads/" + cfgHook.Branch {
				hook = cfgHook
				break
			}
		}

		if hook.Shell == "" {
			log.Printf("Shell command not set in webhook for %s branch %s.  Got:\n%s", data.Repository.FullName,
				data.Ref, string(body))
			return
		}

		project := hook.Repo[strings.LastIndex(hook.Repo, "/")+1:]
		if strings.HasPrefix(data.Ref, "refs/tags/") && !data.Deleted {
			go executeShell(hook.Shell, data.Repository.FullName, project, hook.Branch, "tag", data.Ref[10:])
		} else if data.Ref == "refs/heads/"+hook.Branch && !data.Deleted {
			go executeShell(hook.Shell, data.Repository.FullName, project, hook.Branch, "push", data.After)
		} else {
			log.Printf("Unhandled webhook for %s branch %s.  Got:\n%s", data.Repository.FullName,
				data.Ref, string(body))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\":\"OK\"}"))
		success = true
	})
}

func executeShell(shell string, args ...string) {
	log.Printf("Script starting")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	jobId := r.Uint32()

	commit := args[4]
	if args[3] == "push" {
		commit = commit[:6]
	}

	log.Println("Executing command repo=%s jobId=%s ref=%s ", args[0], strconv.FormatInt(int64(jobId), 10), commit)

	logStreamerOut := NewLogWriter()
	logStreameErr := NewLogWriter()

	shellPath := *scriptsPath + shell
	logStreamerOut.Write([]byte(fmt.Sprintf("Running %s %s\n", shellPath, strings.Join(args, " "))))
	cmd := exec.Command(shellPath, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = logStreamerOut
	cmd.Stderr = logStreameErr

	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if _, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Command finished with error: %v\n", err)
				return
			}
		} else {
			log.Printf("Command finished with error (2): %v\n", err)
			return
		}
	}
	log.Printf("Command finished successfully")
}

var (
	port         = flag.String("port", "7654", "port to listen on")
	configFile   = flag.String("config", "./config.json", "config")
	logFile      = flag.String("log", "", "log file")
	scriptsPath  = flag.String("scriptsPath", "", "path to hook scripts")
	githubSecret = flag.String("secret", "", "github hook secret")
)

func init() {
	flag.Parse()
}

func main() {
	setLog(logFile)
	loadConfig(configFile)
	startWebserver()
}
