package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type PathNodes struct {
	PathNodes []Path
}

type Path struct {
	Path       string
	Permission string
	Size       int64
	Hash       string
	Verify     bool
}

func main() {
	filesIgnore := []string{"node_modules", ".cache", ".git", ".node-build"}
	root := "/home/pepper"
	paths := make([]Path, 0)
	filepath.Walk(root, func(pathStr string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for i := range filesIgnore {
			if regexp.MustCompile(filesIgnore[i]).MatchString(pathStr) {
				return nil
			}
		}
		content, _ := ioutil.ReadFile(pathStr)
		hash := md5.Sum([]byte(content))
		pathObj := Path{
			Path:       pathStr,
			Hash:       hex.EncodeToString(hash[:]),
			Size:       info.Size(),
			Permission: info.Mode().String(),
		}

		paths = append(paths, pathObj)
		// Perform some action on the file or directory
		return nil
	})
	toJson(paths)

	paths_old := fromJson("test_old.json")

	for i := range paths {
		for j := range paths_old {
			if paths[i].Path == paths_old[j].Path {
				paths[i].Verify = true
			}
		}
	}
	toJson(paths)
}

func fromJson(fileName string) []Path {
	file, _ := ioutil.ReadFile(fileName)
	data := []Path{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func toJson(paths []Path) {
	file, _ := json.MarshalIndent(paths, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
}
