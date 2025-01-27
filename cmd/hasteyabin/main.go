package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/duckfullstop/yabingo/pkg/yabin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] hastebin-data-dir \n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Please read the README.md for full information on how to use this program before continuing.")
	}
	flgYAAPIHost := flag.String("api-host", "", "yabin API host")
	flgYAAPIToken := flag.String("api-token", "", "yabin API token from cookie")
	flgYAAPIIgnoreTls := flag.Bool("api-ignore-tls", false, "ignore TLS validation issues")

	flgNoLangDetect := flag.Bool("no-autodetect", false, "don't detect paste content language")
	flag.Parse()

	argDir := flag.Arg(0)

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if argDir == "" {
		log.Fatal("🚫Argument 'directory' not provided")
	}

	if flgYAAPIHost == nil || *flgYAAPIHost == "" {
		log.Fatal("🚫You must provide a YABin API host")
	}

	if flgYAAPIToken == nil || *flgYAAPIToken == "" {
		log.Print("⚠️ YABin API token not provided. If your YABin installation has the envvar 'PUBLIC_CUSTOM_PATHS_ENABLED' set to 'false', all pastes will get randomly assigned keys. You have been warned!")
		log.Print(" - Press CTRL+C to exit in the next 5 seconds, otherwise proceeding...")
		time.Sleep(5 * time.Second)
	}

	if *flgYAAPIIgnoreTls == true {
		log.Print("⚠️ TLS validation disabled 😱")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		yabin.SetClient(client)
	}

	files, err := os.ReadDir(argDir)
	if err != nil {
		log.Fatal(err)
	}

	// Load paste key data into memory, but retain pointer to file contents
	// that way we don't have to load them all into memory in one go
	databuf := make(map[string]*os.File)

	for _, f := range files {
		file, err := os.OpenFile(filepath.Join(argDir, f.Name()), os.O_RDONLY, 0)
		if err != nil {
			log.Fatal(err)
		}
		databuf[f.Name()] = file
	}

	// Defer closing all file handlers
	defer func() {
		for _, f := range databuf {
			f.Close()
		}
	}()

	log.Printf("📚 Sending %d pastes...", len(databuf))

	// Requests to the YABin are intentionally NOT async to avoid flooding it
	yabin.SetURL(*flgYAAPIHost)
	//goland:noinspection GoDfaNilDereference
	yabin.SetAPIToken(*flgYAAPIToken)

	errbuf := make(map[string]error)
	for key, file := range databuf {
		// yabin.PasteWithKeyLanguage
		filestats, err := file.Stat()
		if err != nil {
			errbuf[key] = err
			continue
		}
		buf := make([]byte, filestats.Size())
		_, err = file.Read(buf)
		if err != nil {
			errbuf[key] = err
			continue
		}
		var newKey string
		if flgNoLangDetect != nil && *flgNoLangDetect == true {
			newKey, err = yabin.PasteWithKeyLanguage(string(buf), key, "plaintext", false)
		} else {
			newKey, err = yabin.PasteWithKey(string(buf), key, false)
		}

		if err != nil {
			errbuf[key] = err
			continue
		}
		if newKey != key {
			errbuf[key] = fmt.Errorf("requested key %s does not match new YABin key %s", newKey, key)
		}
	}

	if len(errbuf) > 0 {
		log.Println("😵‍💫 Errors reported when sending pastes to YABin:")
		for key, err := range errbuf {
			log.Printf(" - %s: %s", key, err)
		}
		if len(errbuf) < len(databuf) {
			log.Println("🎉 All other pastes have been sent to YABin successfully")
		}

		os.Exit(1)
	}

	log.Print("🎉 All pastes have been sent to YABin")
	return
}
