package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "strings"
)

func main() {
  excludeDirs := flag.String("I", "", "Exclude directories (comma-separated)")
  outputFile := flag.String("o", "", "output file")
  flag.Parse()

  directory := "."
  args := flag.Args()
  if len(args) > 0 {
    directory = args[0]
  }
  var excludeList []string
  if *excludeDirs != "" {
    excludeList = strings.Split(*excludeDirs, ",")
  }

  var out *os.File
  var err error
  if *outputFile != "" {
    out, err = os.Create(*outputFile)
  } else {
    out, err = os.Create("output.txt")
  }
  if err != nil {
    fmt.Println(err)
    return
  }
  defer out.Close()

  walkFunc := func(path string, info os.FileInfo, err error) error {
    if err != nil {
      fmt.Printf("error accessing a path %q: %v\n", path, err)
      return err
    }

    relPath, _ := filepath.Rel(directory, path)

    if info.IsDir() {
      for _, excludeDir := range excludeList {
        if strings.Contains(relPath, excludeDir) {
          return filepath.SkipDir
        }
      }
    } else {
      if !info.IsDir() && isTextFile(path) {
        log.Println(path)
        data, err := ioutil.ReadFile(path)
        if err != nil {
          return err
        }
        fmt.Fprintln(out, path)
        fmt.Fprintln(out, "---------------------------------------------------------------------")
        fmt.Fprintln(out, string(data))
        fmt.Fprintln(out, "---------------------------------------------------------------------")
        fmt.Fprintln(out, "\n")
      }

    }
    return nil
  }

  err = filepath.Walk(directory, walkFunc)
  if err != nil {
    fmt.Printf("error walking the path %v: %v\n", directory, err)
  }
}

func isTextFile(path string) bool {
  ext := strings.ToLower(filepath.Ext(path))
  textExtensions := []string{
    //普通文件
    ".txt", ".md",
    //前端
    ".html", ".js", ".css", ".vue",
    //后端
    ".py", ".go", ".c", ".cpp", ".java", ".cs",
    //android
    ".kt", ".gradle", ".pro",
    //脚本语音
    ".php", ".rb", ".pl", ".lua",
    // 脚本文件
    ".bat", ".sh",
    //配置文件
    ".yaml", ".conf", ".json", ".xml", "*.properties",
    // git
    ".gitignore",
  }

  for _, e := range textExtensions {
    if ext == e {
      return true
    }
  }

  return false
}
