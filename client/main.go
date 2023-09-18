package main
import (
	"os"

	 "github.com/alt-coder/ForwardIt/client/client"
)



func main(){
	if len(os.Args) < 2 {
		panic("Give host address")
	}
	if len(os.Args) < 2 {
		panic("Give port")
	}
	
	client.Client(50509,os.Args[1])
}

