package main

import (
  "fmt"
  "log"
  b64 "encoding/base64"
  "os"
  "os/exec"
  "strings"
  "net/http"
  "io/ioutil"
  "github.com/gofiber/fiber/v2"
)


func main() {
  app := fiber.New()
  fmt.Println("[DEBUG] process started")

  /*
    /api/bs64 endpoint removes the background from a given base64-encoded image content
    and returns the result image content in base64 (by default)
  */
  app.Post("/api/bs64", func (c *fiber.Ctx) error {
    c.Accepts("application/json")

    // get & parse body data
    payload := struct {
      Filename string `json:"filename"`
      Content string `json:"content"`
    }{}
    if err := c.BodyParser(&payload); err != nil {
      return c.JSON(err)
    }
    
    // decode the base64 encoded image content & write to a temporary file at /tmp
    rawDecodedContent, err := b64.StdEncoding.DecodeString(payload.Content)
    if err != nil {
      return c.JSON(err)
    }

    // to prevent the user from browsing in the server's directory
    filenameParts := strings.Split(payload.Filename, "/")
    safeFilename := filenameParts[len(filenameParts)-1]

    WriteContentToFile(rawDecodedContent, "../tmp/"+safeFilename)
    n, newFileData := rembg(safeFilename) // remove the background of the file

    return c.JSON(&fiber.Map{
      "filename": payload.Filename,
      "content": newFileData[:n],
    })
  })


  /*
    /api/file endpoint removes the background from a given (non-encoded) binary image content
    and returns the result image content in base64 (by default)
  */
  app.Post("/api/file", func(c *fiber.Ctx) error {
    c.Accepts("application/json")
    
    // get & parse body data
    payload := struct {
      Filename string `json:"filename"`
      Content string `json:"content"`
    }{}
    if err := c.BodyParser(&payload); err != nil {
      return c.JSON(err)
    }

    // to prevent the user from browsing in the server's directory
    filenameParts := strings.Split(payload.Filename, "/")
    safeFilename := filenameParts[len(filenameParts)-1]

    WriteContentToFile([]byte(payload.Content), "../tmp/"+safeFilename)
    n, newFileData := rembg(safeFilename) // remove the background of the file

    return c.JSON(&fiber.Map{
      "filename": payload.Filename,
      "content": newFileData[:n],
    })
  })


  /*
    /api/remote endpoint removes the background from an image given by an URI where it is stored at
    and returns the result image content in base64 (by default)
  */
  app.Post("/api/remote", func(c *fiber.Ctx) error {
    c.Accepts("application/json")

    // get & parse body data
    payload := struct {
      Filename string `json:"filename"`
      Uri string `json:"uri"`
    }{}
    if err := c.BodyParser(&payload); err != nil {
      return c.JSON(err)
    }
    
    // fetch the remote image file contents
    resp, err := http.Get(payload.Uri)
    if err != nil {
      return c.JSON(err)
    }
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return c.JSON(err)
    }

    // to prevent the user from browsing in the server's directory
    filenameParts := strings.Split(payload.Filename, "/")
    safeFilename := filenameParts[len(filenameParts)-1]

    WriteContentToFile(data, "../tmp/"+safeFilename)
    n, newFileData := rembg(safeFilename) // remove the background of the file

    return c.JSON(&fiber.Map{
      "filename": payload.Filename,
      "content": newFileData[:n],
    })
  })

  log.Fatal(app.Listen(":3000"))
}


/*
  creates a new file named $filename and writes $content([]byte) to it
  @params:
  * filename (string)
  * content ([]byte)
  @return:
  n (integer) - representing the number of total bytes written to the newly created file
*/
func WriteContentToFile(content []byte, filename string) int {
  f, err := os.Create(filename)
  isError(err)
  defer f.Close()

  n, err := f.Write(content)
  isError(err)

  return n
}


/*
  removes the background of a given image filename using RemBG(github.com/danielgatis/rembg),
  the function assumes that a file named $filename exists at ../tmp and of an image type.
  @params:
  * filename (string)
  @return:
  (
    n (integer) - representing the number of total bytes written to the result file
    newFileData ([]byte) - the result image content data
  )
*/
func rembg(filename string) (int, []byte) {
  err := exec.Command("rembg", "i", "../tmp/"+filename, "../tmp/output_"+filename).Run()
  isError(err)

  newFileData := make([]byte, 16000000) // allocate 16MB of memory for the output image
  f, err := os.Open("../tmp/output_"+filename)
  isError(err)

  n, err := f.Read(newFileData)
  isError(err)
  defer f.Close()

  /* delete the I/O files from /tmp */
  err = os.Remove("../tmp/"+filename)
  isError(err)
  err = os.Remove("../tmp/output_"+filename)
  isError(err)

  return n, newFileData
}


/* checks if a given error has occurred */
func isError(e error){
  if e != nil {
    panic(e)
  }
}
