package main

import (
	"bufio"
	"fmt"
	"flag"
	"github.com/mattn/go-xmpp"
	"github.com/mattn/go-iconv"
  "net/http"
	"log"
	"os"
	"strings"
  "io/ioutil"
)

// Hipchat Default 
var server   = flag.String("server", "conf.hipchat.com:5223", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")

func fromUTF8(s string) string {
	ic, err := iconv.Open("char", "UTF-8")
	if err != nil {
        fmt.Println("1")
		return s
	}
	defer ic.Close()
	ret, _ := ic.Conv(s)
	return ret
}

func toUTF8(s string) string {
	ic, err := iconv.Open("UTF-8", "char")
	if err != nil {
        fmt.Println("2")
		return s
	}
	defer ic.Close()
	ret, _ := ic.Conv(s)
	return ret
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: example [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *username == "" || *password == "" {
		flag.Usage()
	}

    // Print out the deets for debug
    fmt.Printf("Server: %s \nUsername: %s\nPassword: %s\n", *server, *username, *password)

	talk, err := xmpp.NewClient(*server, *username, *password)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
        msg, err := talk.Recv()
        if err != nil {
            log.Fatal(err)
        }
        // Assert that chat is of type xmpp.Chat
        if chat, ok := msg.(xmpp.Chat); ok {
          if chat.Text == "" {
            fmt.Println(chat.Remote +" is typing...")    
          } else if strings.Contains(chat.Text, "chuck")  {
            fmt.Println("Chuck Norris is prepping a fact")    

            resp, err := http.Get("http://api.icndb.com/jokes/random")

            if err != nil {
              return
            }

            test, err1 := ioutil.ReadAll(resp.Body)

            if err1 != nil {
              return
            }

            talk.Send(xmpp.Chat{Remote: chat.Remote, Type: "chat", Text: string(test)})

            err2 := resp.Body.Close()

            if err2 != nil {
              return
            }
          } else {
            fmt.Println(chat)
          }

        // Now it's a xmpp.Presence
        } else if presence, ok := msg.(xmpp.Presence); ok {
            fmt.Println(presence)
        }
		}
	}()
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimRight(line, "\n")

		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) == 2 {
			talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: toUTF8(tokens[1])})
		}
	}
}
