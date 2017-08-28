package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)


const FIRST_EDGENODE_POS = 2
const MAX_EDGENODES = 100
const ANY_EDGENODE = "any"
const START_PICRESIZE_TEMPLATE = "start-picresize.yaml"
const SERVER_PORT = "8888"
const GRAFANA_PORT = "3000"
const KIBANA_PORT = "5601"
const PROMETHEUS_PORT = "9090"
const KUBERNETES_PORT = "9999"
const ADMIN_SERVICES_HOST = "114.115.138.63"

func enumNodes(destURL string) []string {

    var nodeList []string

    for count := FIRST_EDGENODE_POS; count < MAX_EDGENODES; count++ {
        cmdStr := fmt.Sprintf("/usr/bin/kubectl" +
            " get nodes | awk 'NR==%d {print $1}'", count)
        cmd := exec.Command("sh", "-c", cmdStr)
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Println(err)
        }
        nodeName := string(out)
        if nodeName == "" { break }
        nodeList = append(nodeList, nodeName)
    }

    return nodeList
}


func processPic(edgenode string, filename string) {

    srcpng := ""
    node_selector := ""
    hostname_selector := ""

    edgenode = strings.TrimSpace(edgenode)
    filename = strings.TrimSpace(filename)

    if edgenode == ANY_EDGENODE {
        //let kubernetes to descide which edge node will be schedule
        edgenode = ""
    }

    if filename != "" {
        srcpng = filename
    } else { return }

    if edgenode != "" {
        node_selector     = "nodeSelector:"
        hostname_selector = fmt.Sprintf("kubernetes.io/hostname: \"%s\"",
            edgenode)
    }

    uuidOut, err := exec.Command("sh", "-c", "uuidgen").CombinedOutput()
    if err != nil {
        fmt.Println(err)
    }
    uuid := string(uuidOut)
    uuid = strings.TrimSpace(uuid)

    yamlFile, err := ioutil.ReadFile(START_PICRESIZE_TEMPLATE)
    if err != nil {
        fmt.Println(err)
    }

    strYaml := string(yamlFile)
    strTargetYaml := strings.Replace(strYaml, "$$SRCPNG$$", srcpng, -1)
    strTargetYaml = strings.Replace(strTargetYaml, "$$NODE_SELECTOR$$",
        node_selector, -1)
    strTargetYaml = strings.Replace(strTargetYaml, "$$HOSTNAME_SELECTOR$$",
        hostname_selector, -1)
    strTargetYaml = strings.Replace(strTargetYaml, "$$PICRESIZE_JOB$$",
        "picresize-"+uuid, -1)
    strTargetYaml = strings.Replace(strTargetYaml, "$$PICRESIZE_CONTAINER$$",
        "picresize-"+uuid, -1)

    fmt.Println("strTargetYaml: ", strTargetYaml)

    tmpfile, err := ioutil.TempFile("./temp", "temp")
    if err != nil {
        log.Fatal(err)
    }

    // defer os.Remove(tmpfile.Name()) // clean up

    content := []byte(strTargetYaml)
    if _, err := tmpfile.Write(content); err != nil {
        log.Fatal(err)
    }

    cmdStr := fmt.Sprintf("/usr/bin/kubectl" +
            " create -f %s", tmpfile.Name())
    fmt.Println(cmdStr)

    cmd := exec.Command("sh", "-c", cmdStr)
    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("exec.command kubectl create -f ", tmpfile.Name(), out)

    if err := tmpfile.Close(); err != nil {
        log.Fatal(err)
    }
}


func reviewPic(w http.ResponseWriter, edgenode string,
    filename string, server_host string) {

    edgenode    = strings.TrimSpace(edgenode)
    filename    = strings.TrimSpace(filename)
    server_host = strings.TrimSpace(server_host)

    if edgenode == "" {
        edgenode = " *** all edge nodes ***"
    } else {
        edgenode = "edge node *** " + edgenode + " ***"
    }

    files_filter := "./static/" + filename + "*"

    fmt.Println("images files_filter:", files_filter)

    files, _ := filepath.Glob(files_filter)

    imgList := "<h1>The resized pictures handled by " + edgenode + "</h1>"
    for _, f := range files {
        img := fmt.Sprintf("<br><br><br><img src=http://%s/%s /img>", server_host, f)
        imgList = imgList + img
    }

    fmt.Println("imgList: ", imgList)

    fmt.Fprintf(w, imgList)
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Welcome to picture resize service at edge node</h1>" +
        "<ul>" +
        "<li><a href=\"http://%s/upload\">uploading a picture</a></li><br><br>" +
        "<li><a href=\"http://%s/review\">review resized pictures</a></li></ul>" +
        "<br><br><br><br><br><br>" +
        "<h1>!!! Links only for administration !!!</h1>" +
        "<ul><li><a href=\"http://%s/admin\">If you are administrator, click here</a></li></ul>",
        r.Host, r.Host, r.Host)
}


func adminHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Welcome to edge nodes administration</h1>" +
        "<ul><li><a href=\"http://"+ADMIN_SERVICES_HOST+":"+GRAFANA_PORT+"\">" +
            "Monitoring Dashboard of Grafana+Promethieus</a></li><br><br>" +
        "<li><a href=\"http://"+ADMIN_SERVICES_HOST+":"+PROMETHEUS_PORT+"\">" +
            "Monitoring Dashboard of Promethieus</a></li><br><br>" +
        "<li><a href=\"http://"+ADMIN_SERVICES_HOST+":"+KIBANA_PORT+"\">" +
            "Logging Dashboard of Kibana+ElasticSearch</a></li><br><br>" +
        "<li><a href=\"http://"+ADMIN_SERVICES_HOST+":"+KUBERNETES_PORT+"\">" +
            "Edge Nodes Management Dashboard of Kubernetes</a></li></ul>")
}


func uploadHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Println("uploadHandler method:", r.Method)

    if r.Method == "GET" {

        nodeList := enumNodes("")
        var buff bytes.Buffer
        buff.WriteString("<select name=\"edgenode\">")
        buff.WriteString("<option value=\"any\">any</option>")
        for _, node := range nodeList {
            buff.WriteString("<option value=\"" +
                node + "\">" + node + "</option>")
        }
        buff.WriteString("</select>")

        fmt.Fprintf(w, "<h1>please upload picture</h1>" +
            "<form enctype=\"multipart/form-data\" action=\"/upload/\" method=\"post\">" +
            "<input type=\"file\" name=\"uploadfile\" accept=\".png\" /><br><br><br>" +
            "please select an edge node to run" + buff.String() +
            "<br><br><br><input type=\"submit\" value=\"Submit\" />" + " </form>")
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Println("handler.Filename: ", handler.Filename)
        plainFilename := strings.TrimSpace(handler.Filename)
        plainFilename = strings.Replace(handler.Filename, " ", "_", -1)
        fmt.Println("plainFilename: ", plainFilename)

        f, err := os.OpenFile("./static/"+plainFilename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        io.Copy(f, file)

        f.Close()
        file.Close()

        edgenodes := r.MultipartForm.Value["edgenode"]
        fmt.Println("edgenodes: ", edgenodes)
        if len(edgenodes) < 1 { return }

        node_selected := strings.TrimSpace(edgenodes[0])
        fmt.Println("node_selected: ", node_selected)

        processPic(node_selected, plainFilename)

        time.Sleep(time.Second*5)
        reviewPic(w, node_selected, plainFilename, r.Host)
    }
}

func reviewHandler(w http.ResponseWriter, r *http.Request) {
    node_selected := ""
    filename      := ""
    reviewPic(w, node_selected, filename, r.Host)
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/upload/", uploadHandler)
    http.HandleFunc("/review/", reviewHandler)
    http.HandleFunc("/admin/", adminHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    http.ListenAndServe(":"+SERVER_PORT, nil)
}
