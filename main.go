package main

import (
	"bytes"
	"fmt"
	ic "github.com/egirna/icap-client"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Parse command line arguments
	params := parseArgs()

	// Create ICAP request
	var icapReq *ic.Request
	var err error

	switch params.getMode() {
	// Response
	case ic.MethodRESPMOD:
		extResp, err := http.Get(params.Url)
		if err != nil {
			log.Fatal(err)
		}

		if icapReq, err = ic.NewRequest(ic.MethodRESPMOD, params.IcapUrl, nil, extResp); err != nil {
			log.Fatal(err)
		}

	// Request
	case ic.MethodREQMOD:
		file, err := os.Open(params.File)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			log.Fatal(err)
		}
		err = writer.Close()
		if err != nil {
			log.Fatal(err)
		}

		httpReq, err := http.NewRequest("POST", "https://example.com/upload", body)
		if err != nil {
			log.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", writer.FormDataContentType())

		if icapReq, err = ic.NewRequest(ic.MethodREQMOD, params.IcapUrl, httpReq, nil); err != nil {
			log.Fatal(err)
		}

	// Options
	case ic.MethodOPTIONS:
		if icapReq, err = ic.NewRequest(ic.MethodOPTIONS, params.IcapUrl, nil, nil); err != nil {
			log.Fatal(err)
		}
	}

	client := &ic.Client{
		Timeout: time.Duration(params.Timeout) * time.Second,
	}

	resp, err := client.Do(icapReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ICAP Mode:", params.getMode())
	fmt.Printf("Response Status: (%d) %s\n", resp.StatusCode, resp.Status)

	// Вывод заголовков ответа
	fmt.Println("Response Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	if params.Verbose {
		httpResp := resp.ContentResponse

		// Проверка на nil
		if httpResp == nil {
			log.Fatal("httpResp is empty")
		}

		// Вывод статуса ответа
		fmt.Println("Response status:", httpResp.Status)

		// Вывод заголовков ответа
		fmt.Println("Headers:")
		for key, values := range httpResp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}

		// Чтение тела ответа
		defer httpResp.Body.Close()
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			} else {
				log.Fatal(err)
			}
		}
		fmt.Println("Body:", string(body))
	}

	fmt.Println("Execution finished!")
}
