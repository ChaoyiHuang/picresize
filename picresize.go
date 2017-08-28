package main

import (
    "bytes"
    "github.com/nfnt/resize"
    "github.com/namsral/flag"
    "image/png"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
    "path"
)


func ResizePng(r io.Reader, size uint, filepath string) {

    // decode png into image.Image
    img, err := png.Decode(r)
    if err != nil {
        panic(err)
    }

    // resize to width 60 using Lanczos resampling
    // and preserve aspect ratio
    m := resize.Resize(size, 0, img, resize.Lanczos3)

    out, err := os.Create(filepath)
    if err != nil {
         panic(err)
    }
    defer out.Close()

    // write new image to file
    png.Encode(out, m)
}


func postFile(filename string, targetUrl string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    fh, err := os.Open(filename)
    if err != nil {
        fmt.Println(filename)
        fmt.Println("error opening file")
        return err
    }

    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))

    fh.Close()

    os.Remove(filename)

    fmt.Println(filename)

    return nil
}

func uploadHandler(srcURL string, targetURL string, nodeName string) {

    // resize to 128 and 256 width respectively
    // open image on web server
    resp, err := http.Get(srcURL)
    if err != nil {
        fmt.Println(err)
        return
    }

    _, fileName := path.Split(srcURL)

    dest1 := "./" + fileName + "." + nodeName + ".128.png"
    ResizePng(resp.Body, 128, dest1)
    resp.Body.Close()
    postFile(dest1, targetURL)

    //should copy the content from the resp.Body
    //instead for the second resizing
    resp2, err2 := http.Get(srcURL)
    if err2 != nil {
        fmt.Println(err)
        return
    }
    dest2 := "./" + fileName + "." + nodeName + ".256.png"
    ResizePng(resp2.Body, 256, dest2)
    postFile(dest2, targetURL)
}


func main() {

    var targetURL, srcURL, nodeName string

    flag.StringVar(&targetURL, "TARGETURL",
        "http://114.115.138.63:8888/upload/", "the uploading URL")
    flag.StringVar(&srcURL, "SRCURL",
        "http://114.115.138.63:8888/static/Tricircle_view.png", "the file to be resized")
    flag.StringVar(&nodeName, "MY_NODE_NAME",
        "unknown", "the node where the container is running on")

    flag.Parse()

    fmt.Println("targetURL:", targetURL)
    fmt.Println("srcURL:", srcURL)
    fmt.Println("nodeName", nodeName)

    uploadHandler(srcURL, targetURL, nodeName)
}
