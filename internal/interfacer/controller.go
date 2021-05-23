package interfacer

import "fyne.io/fyne/v2/widget"

func (i *Interfacer) controller(pb *widget.ProgressBar, timeLeftLabel *widget.Label) {
	for {
		select {
		case ticks := <-i.ticksChan:
			i.mu.Lock()
			pb.SetValue(ticks)
			i.mu.Unlock()
		case tl := <-i.timeLeftChan:
			i.mu.Lock()
			timeLeftLabel.SetText(tl)
			i.mu.Unlock()
		case pbReset := <-i.ticksReset:
			i.mu.Lock()
			pb.SetValue(0)
			pb.Max = pbReset
			i.mu.Unlock()
		case <-i.exitChan:
			return
		}
	}
}
