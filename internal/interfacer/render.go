package interfacer

import (
	"embed"
	"fmt"
	"log"
	"net/url"

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

const (
	width, height = 75, 75
)

// Render renders interface.
// nolint: funlen
func (i *Interfacer) Render(fs embed.FS) (fyne.Window, error) {
	myApp := app.New()
	myWindow := myApp.NewWindow("procrastination killer")

	pomodoro, err := fs.Open("static/pomodoro.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	play, err := fs.ReadFile("static/play.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	pause, err := fs.ReadFile("static/pause.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	stop, err := fs.ReadFile("static/stop.png")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	defer func() {
		err := pomodoro.Close()
		if err != nil {
			log.Println(errors.Wrap(err, "failed to close pomodoro.png"))
		}
	}()

	img := canvas.NewImageFromReader(pomodoro, "pomodoro")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(width, height))

	timeLeft := widget.NewLabel(fmt.Sprintf("%02d:00", i.workTimeMin))
	timeLeft.TextStyle.Bold = true
	pb := widget.NewProgressBar()

	playIcon := fyne.NewStaticResource("play", play)
	pauseIcon := fyne.NewStaticResource("pause", pause)
	stopIcon := fyne.NewStaticResource("stop", stop)

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
		switcher = (switcher + 1) % 2 // nolint: mnd
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
