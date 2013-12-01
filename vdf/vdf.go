package vdf

import (
  "encoding/json"
  "io"
  "io/ioutil"
  "log"
  "os"
)

var depth = 0

func parse(s []byte, i int) (map[string]interface{}, int) {
  deserialized := map[string]interface{}{}
  var nextIsValue bool
  var letter bool
  var lastStr string
  var ins string

  for i < len(s) {
    c := s[i]
    switch c {
    case '{':
      depth++
      l(i, "Opening node", c)
      letter = false
      nextIsValue = false
      var parsed map[string]interface{}
      parsed, ptr := parse(s, i+1)
      deserialized[lastStr] = parsed
      i = ptr
    case '}':
      depth--
      l(i, "Closing Node", c)
      letter = false
      return deserialized, i
    case '\r', '\n':
      l(i, "New Line", c)
      letter = false
      switch len(ins) {
      case 0:
      default:
        switch nextIsValue {
        case true:
          deserialized[lastStr] = ins
          nextIsValue = false
        case false:
          lastStr, ins = ins, ""
        }
      }
    case ' ', '\t':
      l(i, "Space/Tab", c)
      letter = false
      if len(ins) > 0 {
        lastStr, ins = ins, ""
        nextIsValue = true
      }
    case '"':
      l(i, "Level \"", c)
      letter = false
      ins, i = getToken(s, i+1)
    default:
      // Should be only letters not inside of " "
      l(i, "Default", c)
      switch letter {
      case true:
        l(i, "Add to ins", c)
        ins += string(c)
      case false:
        l(i, "Reset Ins", c)
        ins = string(c)
        letter = true
      }
    }
    i++
  }
  return deserialized, i
}

func getToken(s []byte, i int) (string, int) {
  var ins string
  var add int
  var b byte
  for add, b = range s[i:] {
    switch b {
    case str:
      return ins, i + add
    default:
      ins += string(b)
    }
  }
  log.Panic("Malformed VDF. Unclosed quotes.")
}

func ToJson(body io.ReadCloser) ([]byte, error) {
  b, err := ioutil.ReadAll(body)
  if err != nil {
    return nil, err
  }
  serialized, _ := parse(b, 0)
  return json.Marshal(serialized)
}

func l(i int, name string, c byte) {
  if os.Getenv("debug") == "true" {
    log.Printf("[%d] (%d) - %s: %s %v", i, depth, name, string(c), c)
  }
}
