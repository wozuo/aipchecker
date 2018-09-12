package checker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type folder struct {
	name string
	path string
}

var wg sync.WaitGroup

func (f *folder) check() error {
	manifestPath, err := findManifest(f.path)

	if err != nil {
		return err
	}

	if manifestPath != "" {
		hasInternetPerm, err := checkPermissions(manifestPath)

		if err != nil {
			return err
		}

		if hasInternetPerm {
			err := f.rename()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkPermissions(manifestPath string) (bool, error) {
	manifestFile, err := os.Open(manifestPath)

	if err != nil {
		return false, err
	}
	defer func() {
		err := manifestFile.Close()

		if err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(manifestFile)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "android.permission.INTERNET") {
			return true, nil
		}
	}

	return false, nil
}

func (f *folder) rename() error {
	dir := filepath.Dir(f.path)
	newName := "[i] " + f.name
	newPath := filepath.Join(dir, newName)
	err := os.Rename(f.path, newPath)

	if err != nil {
		return err
	}

	return nil
}

func findManifest(path string) (string, error) {
	manifestPath := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "AndroidManifest.xml" {
			manifestPath = path
			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return "", err
	}

	return manifestPath, nil
}

// CheckProjects checks all projects in a directory for
// internet permission
func CheckProjects(path string) error {
	folders, err := getFolders(path)

	if err != nil {
		return err
	}

	fChan := make(chan folder, len(folders))
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(fChan)
	}

	for _, f := range folders {
		wg.Add(1)
		fChan <- f
	}

	wg.Wait()
	close(fChan)

	fmt.Println(" [DONE]")

	return nil
}

func getFolders(path string) ([]folder, error) {
	var folders []folder
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, folder{
				name: file.Name(),
				path: filepath.Join(path, file.Name()),
			})
		}
	}

	return folders, nil
}

func worker(fChan chan folder) {
	for folder := range fChan {
		err := folder.check()

		if err != nil {
			panic(err)
		}
		wg.Done()
	}
}
