package downloader

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "log"
)

func RsyncFrom(remotePath, localPath string, down chan bool) {
    _, err := exec.Command("/usr/bin/rsync", remotePath, localPath).Output()
    if err != nil {
        log.Fatal(err)
    }

    //fmt.Println(string(out))

    down <- true
}

func DownloadFromUrl(url string, down chan bool) {
    tokens := strings.Split(url, "/")
    fileName := tokens[len(tokens)-1]
    //fmt.Println("Downloading", url, "to", fileName)

    // TODO: check file existence first with io.IsExist
    output, err := os.Create(fileName)
    if err != nil {
        fmt.Println("Error while creating", fileName, "-", err)
        return
    }
    defer output.Close()

    response, err := http.Get(url)
    if err != nil {
        fmt.Println("Error while downloading", url, "-", err)
        return
    }
    defer response.Body.Close()

    _, err = io.Copy(output, response.Body)
    if err != nil {
        fmt.Println("Error while downloading", url, "-", err)
        return
    }

    //fmt.Println(n, "bytes downloaded.")
    down <- true
}

//func main() {
//    url := "http://mail.cc.sandai.net/office/ip.txt"
//    downloadFromUrl(url)
//}
