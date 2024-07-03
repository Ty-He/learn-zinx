
package main

import (
    "fmt"
    "my_demo/protobufDemo/pb"
    "google.golang.org/protobuf/proto"
)

func main() {
    person := &pb.Person {
        Name : "ty",
        Age : 18,
        Emails : []string{"ty@163.com"},
        Phones : []*pb.PhoneNumber {
            {
                Number : "111111111",
                Type : pb.PhoneType_HOME,
            },
            {
                Number : "111111111",
                Type : pb.PhoneType_HOME,
            },
        },
    }
    fmt.Println("[1]==> ", person)

    data, err := proto.Marshal(person)
    if err != nil {
        fmt.Println("Marshal err", err)
        return
    }
    
    fmt.Println("bin data==>", data)

    var p pb.Person 
    err = proto.Unmarshal(data, &p)
    if err != nil {
        fmt.Println("Unmarshal err", err)
        return
    }

    fmt.Println("[2]==>", &p)
}
