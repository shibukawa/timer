// Copyright 2022 Yoshiki Shibukawa.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image"
	"math"
	"time"

	"github.com/fogleman/gg"
)

func TimerImage(running bool, d time.Duration) image.Image {
	const (
		R       = 160
		Width   = 330
		Height  = 330
		CenterX = Width / 2
		CenterY = Height / 2
		M       = 15 //Other memory line segments
		M5      = 26 //5 minute memory line segment
	)
	dc := gg.NewContext(Width, Height)

	dc.SetHexColor("#ffffff")
	dc.Clear()

	//Arc
	if running {
		angle := float64(d) / float64(12*time.Minute) * math.Pi * 2
		dc.Push()
		dc.SetHexColor("#ff0000")
		dc.MoveTo(CenterX, CenterY)
		dc.DrawArc(CenterX, CenterY, R, -math.Pi/2, angle-math.Pi/2)
		dc.MoveTo(CenterX, CenterY)
		dc.Fill()
		dc.Pop()
	}

	//Memory drawing
	dc.Push()
	dc.SetHexColor("#000000")
	for i := 0; i < 60; i++ {
		dc.Push()
		var m float64 = M
		if i%5 == 0 {
			dc.SetLineWidth(2)
			m = M5
		}
		dc.MoveTo(CenterX, CenterY-R+m)
		dc.LineTo(CenterX, 0)
		dc.Stroke()
		dc.Pop()
		dc.RotateAbout(gg.Radians(float64(6)), CenterX, CenterY)
	}
	dc.Pop()

	dc.SetHexColor("#000")

	//Angle of the hour hand at 0 o'clock=Since we want to set it to 0 degrees, rotate it 90 degrees counterclockwise in advance.
	dc.RotateAbout(gg.Radians(-90), CenterX, CenterY)

	dc.DrawCircle(CenterX, CenterY, R)
	dc.Stroke()

	return dc.Image()
}
