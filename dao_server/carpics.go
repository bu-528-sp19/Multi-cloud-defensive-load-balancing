

//https://aws.amazon.com/getting-started/tutorials/backup-to-s3-cli/
//https://github.com/aws/aws-sdk-go
//https://mrwaggel.be/post/golang-transmit-files-over-a-nethttp-server-to-clients/
//https://astaxie.gitbooks.io/build-web-application-with-golang/en/04.5.html

package main

import (
	//"encoding/json"
	"net/http"
	"strconv"
  //"fmt"
	"io"
	"os"
	//"github.com/gorilla/mux"
)


func GetCarPics(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
  Filename := "car.jpg"
  //Filename := request.URL.Query().Get("file")
	getCarPics(Filename)
  Openfile, _ := os.Open(Filename)
	defer os.Remove(Filename)
  defer Openfile.Close()

  FileHeader := make([]byte, 512)
	Openfile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client
	return

}

func UploadCarPics(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	/*req.ParseMultipartForm(32 << 20)
  file, handler, err := req.FormFile("car2.jpg")
  if err != nil {
  	fmt.Println(err)
    return
  }
  defer file.Close()
  fmt.Fprintf(w, "%v", handler.Header)
  f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
  if err != nil {
  	fmt.Println(err)
    return
  }
  defer f.Close()
  io.Copy(f, file)
	*/
	uploadCarPics("car2.jpg")
	return
}
