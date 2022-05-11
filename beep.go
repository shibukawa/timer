// Copyright 2022 Yoshiki Shibukawa.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	_ "embed"
	"time"

	"bytes"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

//go:embed beep-02.wav
var beepWav []byte

var audioContext *audio.Context
var player *audio.Player

func init() {
	audioContext = audio.NewContext(44100)
	s, err := wav.Decode(audioContext, bytes.NewReader(beepWav))
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(s)
	if err != nil {
		panic(err)
	}
	player = audioContext.NewPlayerFromBytes(b)
}

func playSound() {
	go func() {
		for i := 0; i < 5; i++ {
			player.Play()
			time.Sleep(3 * time.Second)
		}
	}()
}
