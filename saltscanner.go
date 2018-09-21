package main 

import (
    "fmt"
    "encoding/json"
    "bufio"
    //"io"
    "log"
    "net"
    "os"
    "strconv"
    //"strings"
    "sync"
    //"time"

    "./conf"
    "./downloader"
    //"./fileParser"
    "github.com/robfig/cron"
    //"github.com/orcaman/concurrent-map"
)

var config = conf.ParseConfig()
//var conf Config{}
//var MasterCount int
var hostGroupMap sync.Map
var c = sync.NewCond(&sync.Mutex{})
//var wg sync.WaitGroup

type masterMinionsInfo struct{
    Master string
    Minions []string
}

func SocketServer(port int) {

    listen, err := net.Listen("tcp4", ":" + strconv.Itoa(port))
    defer listen.Close()

    if err != nil {
        fmt.Println("Fatal error ", err.Error())
    }
    log.Printf("Begin listen port: %d...\n", port)

    for {
        conn, err := listen.Accept()
        if err != nil {
            fmt.Println("Fatal error ", err.Error())
            continue
        }

        go handler(conn)
    }

}

func handler(conn net.Conn) {

    remoteAddr := conn.RemoteAddr().String()

    scanner := bufio.NewScanner(conn)

    ok := scanner.Scan()

    if !ok {
        log.Println("scan failed!")
        return
    } 

    data := masterMinionsInfo{}

    err := json.Unmarshal([]byte(scanner.Text()), &data)
    if err != nil {
        log.Println(err)
    }

    log.Printf("receive from %s(%s): %d minions", data.Master, remoteAddr, len(data.Minions))

    go filterMinion(data)   
}


func filterMinion(data masterMinionsInfo) {

    c.L.Lock()
    defer c.L.Unlock()

    // wait for the hostGroupMap ready
    for lengthOfSyncMap(hostGroupMap) == 0 {
        fmt.Println("waiting for hostGroupMap...\n")
        c.Wait()
    }

    // how many minions are found in hostGroupMap
    matchedCount := 0
    notmatchedCount := 0

    for _, minion := range data.Minions {
        if _, ok := hostGroupMap.Load(minion); ok {
            hostGroupMap.Delete(minion)
            fmt.Println("    matched: ", minion)
            matchedCount++
        } else {
            fmt.Println("not matched: ", minion)
            notmatchedCount++
        }
    }

    log.Printf("%d minions matched. %d not matched. %d left\n", matchedCount, notmatchedCount, lengthOfSyncMap(hostGroupMap))

    //wg.Done()
}

func initDataSrc() { 

    //wg.Add(MasterCount)

    done := make(chan bool, 1)
    go downloader.RsyncFrom(config.IPFileURL, "/tmp/iplist.txt", done)

    select {
        case <-done:
            log.Println("downloaded iplist\n")

            c.L.Lock()
            hostGroupMap = Parse()
            c.L.Unlock()

            c.Broadcast()
    }

    fmt.Println("hostGroupMap is ready.\n")
    //wg.Wait()
    //fmt.Println(client_conn, len(client_conn))

    //output()
}

// output result
func output() {
    log.Printf("output %d hosts to %s", lengthOfSyncMap(hostGroupMap), config.OutputFile)
    f, err := os.Create(config.OutputFile)
    if err != nil {
        log.Println(err)
        return
    }
    defer f.Close()

    bufWriter := bufio.NewWriter(f)
    defer bufWriter.Flush()

    hostGroupMap.Range(func(k, v interface{}) bool {
        //fmt.Println(k, "--", v)
        bufWriter.WriteString(k.(string) + " " + v.(string) + "\n")
        return true
    })
}

func main() {

    go SocketServer(config.Port)

    //time.Sleep(5 * time.Second)
    go initDataSrc()

    c := cron.New()

    c.AddFunc(config.Output_cron, output)
    c.AddFunc(config.InitDataSrc_cron, initDataSrc)

    c.Start()

    select {}
}


func lengthOfSyncMap(m sync.Map) int {
    length := 0
    hostGroupMap.Range(func(key, value interface{}) bool {
        length++
        return true
    })
    
    return length 
}
