package pkg

import (
	"mime/multipart"
	"strings"
)

func GetFileExtention(header *multipart.FileHeader) string{
  arr := strings.Split(header.Filename, ".")
  if len(arr)==1{
    return ""
  }
  return arr[len(arr)-1]
}
