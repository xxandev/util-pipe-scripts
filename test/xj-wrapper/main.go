package main

import (
	"fmt"
	"log"
	"os"
	"util-pipe/internal/xj"
)

func main() {
	content, err := os.ReadFile("client.json")
	if err != nil {
		log.Fatalln("read test json error:", err)
	}
	object := xj.NewWrap()
	if err := object.Unmarshal(content); err != nil {
		log.Fatalln("unmarshal object json error:", err)
	}

	fmt.Println("===================================")
	fmt.Println("wrapper print:")
	object.Data.Range(func(key, value any) bool {
		fmt.Printf("%v = %#v\n", key, value)
		return true
	})

	fmt.Println("===================================")
	fmt.Println("wrapper load [address.street]:")
	fmt.Println(object.Data.Load("address.street"))

	fmt.Println("wrapper load [contact_info.phone]:")
	fmt.Println(object.Data.Load("contact_info.phone"))

	fmt.Println("wrapper store [contact_info.phone]: +234567891")
	object.Data.Store("contact_info.phone", "+234567891")

	fmt.Println("wrapper load [contact_info.phone]:")
	fmt.Println(object.Data.Load("contact_info.phone"))

	fmt.Println("wrapper store [contact_info.home_phone]: +987654321")
	object.Data.Store("contact_info.home_phone", "+987654321")

	fmt.Println("===================================")
	fmt.Println("wrapper marshal print:")
	tmp, err := object.MarshalIndent("", "    ")
	if err != nil {
		log.Fatalln("error wrapper marshal:", err)
	}
	fmt.Println(string(tmp))
}
