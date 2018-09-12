package zipper

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type archive struct {
	name string
	path string
}

var wg sync.WaitGroup

func (a *archive) unzip() error {
	reader, err := zip.OpenReader(a.path)

	if err != nil {
		return err
	}
	defer func() {
		err := reader.Close()

		if err != nil {
			panic(err)
		}
	}()

	extractedDir := getExtractedDirName(0, string(a.path[0:len(a.path)-4]))
	os.Mkdir(extractedDir, 0755)

	for _, file := range reader.Reader.File {
		err := extract(extractedDir, file)

		if err != nil {
			return err
		}
	}

	return nil
}

func getExtractedDirName(i int, path string) string {
	if i > 0 {
		if i > 1 {
			path = string(path[0:len(path)-1]) + strconv.Itoa(i)
		} else {
			path = path + "_" + strconv.Itoa(i)
		}
	}

	if _, err := os.Stat(path); err == nil {
		i = i + 1
		path = getExtractedDirName(i, path)
	}

	return path
}

func extract(extractedDir string, file *zip.File) error {
	f, err := file.Open()

	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()

		if err != nil {
			panic(err)
		}
	}()

	extrFilePath := filepath.Join(extractedDir, file.Name)

	if file.FileInfo().IsDir() {
		err := os.MkdirAll(extrFilePath, file.Mode())

		if err != nil {
			return err
		}
	} else {
		extrFile, err := os.OpenFile(
			extrFilePath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			file.Mode(),
		)

		if err != nil {
			return err
		}
		defer func() {
			err := extrFile.Close()

			if err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(extrFile, f)

		if err != nil {
			return err
		}
	}

	return nil
}

// UnzipAll unpacks all zip files from a directory
func UnzipAll(path string) error {
	archives, err := getArchives(path)

	if err != nil {
		return err
	}

	aChan := make(chan archive, len(archives))
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(aChan)
	}

	for _, a := range archives {
		wg.Add(1)
		aChan <- a
	}

	wg.Wait()
	close(aChan)

	fmt.Println(" [DONE]")

	return nil
}

func getArchives(path string) ([]archive, error) {
	files, err := ioutil.ReadDir(path)
	var archives []archive

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".zip") {
			archives = append(archives, archive{
				name: file.Name(),
				path: filepath.Join(path, file.Name()),
			})
		}
	}

	return archives, nil
}

func worker(aChan chan archive) {
	for archive := range aChan {
		err := archive.unzip()

		if err != nil {
			panic(err)
		}
		wg.Done()
	}
}
