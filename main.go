package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ncruces/zenity"
)

type Timer struct {
	finish  time.Time
	current time.Duration
	running bool
	quit    bool
	stop    chan struct{}
}

func (t *Timer) Update() error {
	if t.quit {
		return errors.New("stop")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		t.finish = time.Time{}
		t.running = false
		t.stop <- struct{}{}
	} else if t.running {
		now := time.Now()
		if t.finish.Before(now) {
			playSound()
			t.running = false
			t.stop <- struct{}{}
		} else {
			t.current = t.finish.Sub(now)
		}
	}
	return nil
}

func (t *Timer) Draw(screen *ebiten.Image) {
	m := TimerImage(t.running, t.current)
	em := ebiten.NewImageFromImage(m)
	screen.DrawImage(em, &ebiten.DrawImageOptions{})
}

func (g *Timer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 330, 330
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	timer := &Timer{
		stop: make(chan struct{}, 1),
	}
	timer.stop <- struct{}{}

	ebiten.SetWindowSize(330, 330)
	ebiten.SetWindowTitle("12分タイマー (escでリセット)")

	go func() {
		for {
			select {
			case <-timer.stop:
				for {
					t, err := zenity.Entry("時間 (1-12分)", zenity.Title("時間入力"))
					if err != nil {
						timer.quit = true
						return
					}
					m, err := strconv.ParseUint(t, 10, 64)
					if err != nil {
						continue
					}
					if m < 1 || 12 < m {
						continue
					}
					timer.finish = time.Now().Add(time.Minute * time.Duration(m))
					timer.running = true
					break
				}
			case <-ctx.Done():
				timer.quit = true
				return
			}
		}
	}()

	if err := ebiten.RunGame(timer); err != nil {
		fmt.Println(err)
	}
}
