//package fileParser
package main

import (
    "os"
    "io"
    //"fmt"
    "log"
    "bufio"
    "strings"
    "sync"
    //"github.com/orcaman/concurrent-map"

)


func Parse() sync.Map {
//func main() {
    input, err := os.Open("/tmp/iplist.txt")

    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    defer input.Close()


    bufReader := bufio.NewReader(input)

    hostGroupMap := sync.Map{}
    /*
        logical format: 
        {"tw07112": "dcdnhub",
         "tw07113": "dcdnhub"}
    */

    for {
        line, _, err := bufReader.ReadLine()

        if err != nil && err != io.EOF {
            log.Fatal(err)
        } else if err == io.EOF {
            break
        } 

        lineSlice := strings.Fields(string(line))

        host := lineSlice[0]
        if strings.HasPrefix(host, "old") {
            continue
        }

        //_, ok := hostGroupMap.Load(name)
        //if  ok == false {
        //    if len(lineSlice) < 3 {
        //        hostGroupMap.Store(name, "")
        //        //fmt.Println(name)
        //    } else {
        //        hostGroupMap.Store(name, lineSlice[2])
        //    }
        //}

        if len(lineSlice) >= 3 {    // the line has group name within it
            if group := lineSlice[2]; groupIsIncluded(group) {
                _, ok := hostGroupMap.Load(host)
                if  ok == false {
                    hostGroupMap.Store(host, group)
                }
            }
        }
    }

    //fmt.Printf("%T\n%[1]v", hostGroupMap[1])
    //for item := range hostGroupMap.Iter() {
    //    //fmt.Println(item, "===")
    //    fmt.Println(item.Key)
    //}
    
    return hostGroupMap
}

func groupIsIncluded(group string) bool {
    for _, g := range config.IncludeGroup {
        if group == g {
            return true
        }
    } 

    return false
}


//func Parse_2() cmap.ConcurrentMap {
////func main() {
//    input, err := os.Open("/tmp/iplist.txt")
//
//    if err != nil {
//        log.Fatal(err)
//        os.Exit(1)
//    }
//
//    defer input.Close()
//
//
//    bufReader := bufio.NewReader(input)
//
//    hostGroupMap := cmap.New()
//    /*
//        logical format: 
//        {"tw07112": "dcdnhub",
//         "tw07113": "dcdnhub"}
//    */
//
//    for {
//        line, _, err := bufReader.ReadLine()
//
//        if err != nil && err != io.EOF {
//            log.Fatal(err)
//        } else if err == io.EOF {
//            break
//        } 
//
//        lineSlice := strings.Fields(string(line))
//
//        name := lineSlice[0]
//
//        
//        if hostGroupMap.Has(name) == false {
//            if len(lineSlice) < 3 {
//                hostGroupMap.Set(name, "")
//                //fmt.Println(name)
//            } else {
//                hostGroupMap.Set(name, lineSlice[2])
//            }
//        }
//    }
//
//    //fmt.Printf("%T\n%[1]v", hostGroupMap[1])
//    //for item := range hostGroupMap.Iter() {
//    //    //fmt.Println(item, "===")
//    //    fmt.Println(item.Key)
//    //}
//    
//    return hostGroupMap
//}

func Parse_1() map[string]string {
//func main() {
    input, err := os.Open("/tmp/iplist.txt")

    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    //output, err := os.Create("../hostname.txt")

    //if err != nil {
    //    log.Fatal(err)
    //    os.Exit(1)
    //}

    defer input.Close()
    //defer output.Close()

    bufReader := bufio.NewReader(input)
    //bufReader, bufWriter := bufio.NewReader(input), bufio.NewWriter(output)
    //defer bufWriter.Flush()

    hostGroupMap := make(map[string]string)
    /*
        format: 
        {"tw07112": "dcdnhub"}
    */


    for {
        line, _, err := bufReader.ReadLine()

        if err != nil && err != io.EOF {
            log.Fatal(err)
        } else if err == io.EOF {
            break
        } 

        //lineSlice := strings.Split(string(line), " ")
        lineSlice := strings.Fields(string(line))

        //fmt.Println(len(lineSlice), lineSlice)

        //if strings.HasPrefix(lineSlice[0], "old.") {
        //    name := lineSlice[0][4:]
        //}
        name := lineSlice[0]

        _, ok := hostGroupMap[name]
        if !ok {
            if len(lineSlice) < 3 {
                hostGroupMap[name] = ""
                //fmt.Println(name)
            } else {
                hostGroupMap[name] = lineSlice[2]
            }
        }

        //bufWriter.WriteString(lineSlice[0] + "\n")
        //bufWriter.WriteString(lineSlice[2] + "\n")
    }

    //for host, group := range hostGroupMap {
    //    fmt.Println(host, group)
    //}
    
    return hostGroupMap
}
