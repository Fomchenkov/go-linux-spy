/*
 * ported from the logkeys project
 * https://github.com/kernc/logkeys
 *
 * ported by Chris Pergrossi
 * License is same as original project: this is free software
 * with no warranty whatsoever
 */

package main

import "log"

var (
	charKeys  = "1234567890-=qwertyuiop[]asdfghjkl;'`\\zxcvbnm,./<"
	shiftKeys = `!@#$%^&*()_+QWERTYUIOP{}ASDFGHJKL:"~|ZXCVBNM<>?>`
	altgrKeys = ""

	funcKeys = []string{"<Esc>", "<BckSp>", "<Tab>", "<Enter>", "<LCtrl>", "<LShft>", "<RShft>", "<KP*>", "<LAlt>",
		" ", "<CpsLk>", "<F1>", "<F2>", "<F3>", "<F4>", "<F5>", "<F6>", "<F7>", "<F8>", "<F9>", "<F10>", "<NumLk>",
		"<ScrLk>", "<KP7>", "<KP8>", "<KP9>", "<KP->", "<KP4>", "<KP5>", "<KP6>", "<KP+>", "<KP1>", "<KP2>", "<KP3>",
		"<KP0>", "<KP.>", "<F11>", "<F12>", "<KPEnt>", "<RCtrl>", "<KP/>", "<PrtSc>", "<AltGr>", "<Break>", "<Home>",
		"<Up>", "<PgUp>", "<Left>", "<Right>", "<End>", "<Down>", "<PgDn>", "<Ins>", "<Del>", "<Pause>", "<LMeta>",
		"<RMeta>", "<Menu>"}

	charOrFunc = "_fccccccccccccffccccccccccccffccccccccccccfcccccccccccffffffffffffffffffffffffffffff__cff_______ffffffffffffffff_______f_____fff"
)

const (
	Key1          = 2
	KeyEqual      = 13
	KeyQ          = 16
	KeyRightbrace = 27
	KeyA          = 30
	KeyGrave      = 41
	KeyBackslash  = 43
	KeySlash      = 53
	KeyEsc        = 1
	KeyBackspace  = 14
	KeyTab        = 15
	KeyEnter      = 28
	KeyLeftCtrl   = 29
	KeyLeftShift  = 42
	KeyRightShift = 54
	KeyLeftAlt    = 56
	KeySpace      = 57
	KeyKPDot      = 83
	KeyF11        = 87
	KeyF12        = 88
	KeyKPEnter    = 96
	KeyRightCtrl  = 97
	KeyRightAlt   = 100
	KeyDelete     = 111
	KeyPause      = 119
	KeyLeftMeta   = 125
	KeyCompose    = 127
)

func isCharKey(ch uint) bool {
	if ch >= uint(len(charOrFunc)) {
		log.Println("CharKey out of bounds: ", ch)
		return false
	}

	return charOrFunc[ch] == 'c'
}

func isFuncKey(ch uint) bool {
	if ch >= uint(len(charOrFunc)) {
		log.Println("FuncKey out of bounds: ", ch)
		return false
	}

	return charOrFunc[ch] == 'f'
}

func isUsedKey(ch uint) bool {
	if ch >= uint(len(charOrFunc)) {
		log.Println("UsedKey out of bounds: ", ch)
		return false
	}

	return charOrFunc[ch] != '_'
}

func toCharKeysIndex(keycode int) int {
	switch {
	case keycode >= Key1 && keycode <= KeyEqual: // keycodes 2 - 13
		return keycode - 2
	case keycode >= KeyQ && keycode <= KeyRightbrace: // keycodes 16 - 27
		return keycode - 4
	case keycode >= KeyA && keycode <= KeyGrave: // keycodes 30 - 41
		return keycode - 6
	case keycode >= KeyBackslash && keycode <= KeySlash: // keycodes 43 - 53
		return keycode - 7
	}

	return -1
}

func toFuncKeysIndex(keycode int) int {
	switch {
	case keycode == KeyEsc: // 1
		return 0
	case KeyBackspace <= keycode && keycode <= KeyTab: // 14 - 15
		return keycode - 13
	case KeyEnter <= keycode && keycode <= KeyLeftCtrl: // 28 - 29
		return keycode - 25
	case keycode == KeyLeftShift: // 42
		return keycode - 37
	case KeyRightShift <= keycode && keycode <= KeyKPDot: // 54 - 83
		return keycode - 48
	case KeyF11 <= keycode && keycode <= KeyF12: // 87 - 88
		return keycode - 51
	case KeyKPEnter <= keycode && keycode <= KeyDelete: // 96 - 111
		return keycode - 58
	case keycode == KeyPause: // 119
		return keycode - 65
	case KeyLeftMeta <= keycode && keycode <= KeyCompose: // 125 - 127
		return keycode - 70
	}

	return -1
}
