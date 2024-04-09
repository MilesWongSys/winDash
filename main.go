package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"winDash/snippet"
)

func main() {
	if err := monitor(); err != nil {
		log.Fatal(err)
	}
}

func monitor() error {
	keyboardChan := make(chan types.KeyboardEvent, 100)
	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return err
	}
	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	log.Println("开始检测键盘输入：")

	var inputBuffer strings.Builder
	var detectKeyword bool

	for {
		select {
		case <-time.After(5 * time.Minute):
			log.Println("超过五分钟无输入，程序退出")
			return nil
		case <-signalChan:
			log.Println("收到终端信号，程序退出")
			return nil
		case k := <-keyboardChan:
			if k.Message == types.WM_KEYDOWN {
				continue
			}
			//log.Printf("Received %v %v\n", k.VKCode, k.ScanCode)
			if !isLegalCode(k.VKCode) {
				inputBuffer.Reset()
				detectKeyword = false
				continue
			}
			if k.VKCode == types.VK_OEM_3 {
				detectKeyword = true
				inputBuffer.Reset()
				continue
			}

			_, err := inputBuffer.WriteString(getStr(k.VKCode))
			if err != nil {
				log.Fatalf("输入缓存异常，%v", err)
			}
			if sp, isExist := snippet.GetSnippet("`" + inputBuffer.String()); detectKeyword && isExist {
				log.Println(sp)
				err = gui(sp)
				if err != nil {
					log.Fatalf("图形界面运行异常， %v", err)
				}
			}
			continue
		}
	}
}

func isLegalCode(code types.VKCode) bool {
	if code == types.VK_OEM_3 || (code <= types.VK_Z && code >= types.VK_0) {
		return true
	}
	return false
}

func getStr(code types.VKCode) string {
	return strings.ToLower(strings.TrimPrefix(code.String(), "VK_"))
}

func gui(sp *snippet.Snippet) error {
	uid := uuid.NewString()
	a := app.NewWithID(uid)
	window := a.NewWindow(uid)
	window.Show()
	var content *fyne.Container
	if len(sp.VariableList) == 0 {
		content = container.NewVBox(
			widget.NewLabel(sp.CMD),
			widget.NewButton("confirm", func() {
				window.Clipboard().SetContent(sp.CMD)
				window.Hide()
			}),
			widget.NewButton("cancel", func() {
				window.Hide()
			}),
		)
	} else {

	}
	window.SetContent(content)
	window.ShowAndRun()
	return nil
}
