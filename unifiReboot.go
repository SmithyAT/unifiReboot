/*
Copyright © 2023 Christian Schmied
*/
package main

import (
	"flag"
	"fmt"
	"github.com/melbahja/goph"
	"log"
	"os"
)

func main() {
	fmt.Println("Unifi Device Reboot Tool - Version 1.0 - © 2023 Christian Schmied")
	fmt.Println("-----------------------------------------------------------------")

	var ip, user, sshkey string
	var test bool
	flag.StringVar(&ip, "ip", "", "Unifi device IP-address")
	flag.StringVar(&user, "u", "", "Unifi device username")
	flag.StringVar(&sshkey, "i", "", "Unifi device ssh-key")
	flag.BoolVar(&test, "t", false, "Test connection")
	flag.Parse()

	// validate parameters
	if ip != "" && user != "" && sshkey != "" {
		// check if sshkey file exists
		if _, err := os.Stat(sshkey); os.IsNotExist(err) {
			fmt.Println("SSH-Key file does not exist")
			os.Exit(1)
		}

		auth, err := goph.Key(sshkey, "")
		if err != nil {
			panic(err)
		}

		fmt.Print("Connecting to device " + ip + ".....")
		client, err := goph.NewUnknown(user, ip, auth)
		if err != nil {
			fmt.Println("error, unable to connect to the device")
			os.Exit(1)
		}
		defer client.Close()
		fmt.Println("done, rebooting device")

		var cmd string
		if test {
			cmd = "mca-cli-op info"
		} else {
			cmd = "reboot"
		}
		out, err := client.Run(cmd)
		if err != nil {
			log.Fatal(err)
		}

		if test {
			fmt.Println(string(out))
		}

	} else {
		fmt.Println("Parameter -ip <IP-ADDRESS> or -u <USERNAME> or -i <SSH-KEY> missing!")
		os.Exit(1)
	}

}
