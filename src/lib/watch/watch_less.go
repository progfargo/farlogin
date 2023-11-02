package watch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"farlogin/src/app"

	"github.com/fsnotify/fsnotify"
)

var lessWatcher *fsnotify.Watcher

func WatchLess() {
	lessWatcher, _ = fsnotify.NewWatcher()

	if err := filepath.Walk(app.Ini.HomeDir+"/style/less", addLessWatchDir); err != nil {
		panic(err.Error())
	}

	var isActive bool

	go func() {
		for _ = range time.Tick(700 * time.Millisecond) {
			isActive = true
		}
	}()

	go func() {
		for {
			select {
			case event := <-lessWatcher.Events:
				fmt.Printf("%s - %s\n", event.Name, event.Op)
				if isActive {
					compileLess(event)
					refreshBrowser()
					isActive = false
				}
			case err := <-lessWatcher.Errors:
				fmt.Println(err.Error())
			}
		}
	}()

}

func addLessWatchDir(path string, fi os.FileInfo, err error) error {

	if fi.Mode().IsDir() {
		return lessWatcher.Add(path)
	}

	return nil
}

func compileLess(event fsnotify.Event) {
	source := event.Name

	//don't compile less files in lib directory. compile style.less instead.
	if filepath.Dir(source) == app.Ini.HomeDir+"/style/less/lib" {
		source = app.Ini.HomeDir + "/style/less/style.less"
	}

	target := strings.Replace(source, "/style/less/", "/asset/css/", -1)
	target = strings.Replace(target, ".less", ".css", -1)

	cmd := exec.Command("lessc", source, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("less compiled.\n")
}

func refreshBrowser() {
	cmd := exec.Command("xdotool", "search", "--onlyvisible", "--class", "Chrome", "windowfocus", "key", "F5")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("browser refreshed.\n")
}
