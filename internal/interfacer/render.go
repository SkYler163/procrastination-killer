package interfacer

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"

	"github.com/SkYler163/procrastination-killer/internal/model"
)

type buttonSwitcher struct {
	icon   fyne.Resource
	action func()
}

// Render renders interface.
func (i *Interfacer) Render() (fyne.Window, error) {
	myApp := app.New()
	myWindow := myApp.NewWindow("procrastination killer")

	f, err := os.Open("static/pomodoro.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to close pomodoro.png"))
		}
	}()

	img := canvas.NewImageFromReader(f, "pomodoro")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(75, 75))

	timeLeft := widget.NewLabel(fmt.Sprintf("%02d:00", i.workTimeMin))
	timeLeft.TextStyle.Bold = true
	pb := widget.NewProgressBar()

	playIcon, err := fyne.LoadResourceFromPath("static/play.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open play icon")
	}

	pauseIcon, err := fyne.LoadResourceFromPath("static/pause.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open pause icon")
	}

	stopIcon, err := fyne.LoadResourceFromPath("static/stop.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open stop icon")
	}

	var switchButton *widget.Button

	isTimerRunning := false
	bs := map[int8]buttonSwitcher{
		0: {
			icon: playIcon,
			action: func() {
				i.controlSignalsChan <- model.ControlSignalPause
			},
		},
		1: {
			icon: pauseIcon,
			action: func() {
				if !isTimerRunning {
					go i.controller(pb, timeLeft)
					i.timer.Run()
					isTimerRunning = !isTimerRunning
				} else {
					i.controlSignalsChan <- model.ControlSignalPlay
				}
			},
		},
	}
	switcher := int8(0)
	switchButton = widget.NewButtonWithIcon("", playIcon, func() {
		switcher = (switcher + 1) % 2
		switchButton.SetIcon(bs[switcher].icon)
		bs[switcher].action()
	})
	buttons := container.NewHBox(
		switchButton,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", stopIcon, func() {
			switcher = 0
			switchButton.SetIcon(playIcon)
			i.controlSignalsChan <- model.ControlSignalStop
			isTimerRunning = false
			timeLeft.SetText(fmt.Sprintf("%02d:00", i.workTimeMin))
			pb.SetValue(0)
		}),
	)

	timerContainer := container.NewHBox(
		timeLeft, widget.NewSeparator(), pb, widget.NewSeparator(), buttons,
	)
	authorLink := container.NewCenter(
		widget.NewHyperlink("author", &url.URL{Scheme: "https", Host: "github.com", Path: "skyler163"}),
	)

	myWindow.SetContent(container.NewVBox(img, timerContainer, authorLink))
	myWindow.SetFixedSize(true)
	myWindow.SetOnClosed(func() {
		if isTimerRunning {
			i.controlSignalsChan <- model.ControlSignalStop
		}
	})

	return myWindow, nil
}
