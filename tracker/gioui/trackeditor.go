package gioui

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/vsariola/sointu/tracker"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const trackRowHeight = 16
const trackColWidth = 54
const patmarkWidth = 16

type TrackEditor struct {
	TrackVoices         *NumberInput
	NewTrackBtn         *widget.Clickable
	DeleteTrackBtn      *widget.Clickable
	AddSemitoneBtn      *widget.Clickable
	SubtractSemitoneBtn *widget.Clickable
	AddOctaveBtn        *widget.Clickable
	SubtractOctaveBtn   *widget.Clickable
	NoteOffBtn          *widget.Clickable
	trackPointerTag     bool
	trackJumpPointerTag bool
	tag                 bool
	focused             bool
	requestFocus        bool
}

func NewTrackEditor() *TrackEditor {
	return &TrackEditor{
		TrackVoices:         new(NumberInput),
		NewTrackBtn:         new(widget.Clickable),
		DeleteTrackBtn:      new(widget.Clickable),
		AddSemitoneBtn:      new(widget.Clickable),
		SubtractSemitoneBtn: new(widget.Clickable),
		AddOctaveBtn:        new(widget.Clickable),
		SubtractOctaveBtn:   new(widget.Clickable),
		NoteOffBtn:          new(widget.Clickable),
	}
}

func (te *TrackEditor) Focus() {
	te.requestFocus = true
}

func (te *TrackEditor) Focused() bool {
	return te.focused
}

func (te *TrackEditor) Layout(gtx layout.Context, t *Tracker) layout.Dimensions {
	for _, e := range gtx.Events(&te.tag) {
		switch e := e.(type) {
		case key.FocusEvent:
			te.focused = e.Focus
		case pointer.Event:
			if e.Type == pointer.Press {
				key.FocusOp{Tag: &te.tag}.Add(gtx.Ops)
			}
		case key.Event:
			switch e.State {
			case key.Press:
				switch e.Name {
				case key.NameDeleteForward, key.NameDeleteBackward:
					t.DeleteSelection()
					if !(t.NoteTracking() && t.Playing()) && t.Step.Value > 0 {
						t.SetCursor(t.Cursor().AddRows(t.Step.Value))
						t.SetSelectionCorner(t.Cursor())
					}
				case key.NameUpArrow, key.NameDownArrow:
					sign := -1
					if e.Name == key.NameDownArrow {
						sign = 1
					}
					cursor := t.Cursor()
					if e.Modifiers.Contain(key.ModShortcut) {
						cursor.Row += t.Song().Score.RowsPerPattern * sign
					} else {
						if t.Step.Value > 0 {
							cursor.Row += t.Step.Value * sign
						} else {
							cursor.Row += sign
						}
					}
					t.SetNoteTracking(false)
					t.SetCursor(cursor)
					if !e.Modifiers.Contain(key.ModShift) {
						t.SetSelectionCorner(t.Cursor())
					}
					//scrollToView(t.PatternOrderList, t.Cursor().Pattern, t.Song().Score.Length)
				case key.NameLeftArrow:
					cursor := t.Cursor()
					if !t.LowNibble() || !t.Song().Score.Tracks[t.Cursor().Track].Effect || e.Modifiers.Contain(key.ModShortcut) {
						cursor.Track--
						t.SetLowNibble(true)
					} else {
						t.SetLowNibble(false)
					}
					t.SetCursor(cursor)
					if !e.Modifiers.Contain(key.ModShift) {
						t.SetSelectionCorner(t.Cursor())
					}
				case key.NameRightArrow:
					if t.LowNibble() || !t.Song().Score.Tracks[t.Cursor().Track].Effect || e.Modifiers.Contain(key.ModShortcut) {
						cursor := t.Cursor()
						cursor.Track++
						t.SetCursor(cursor)
						t.SetLowNibble(false)
					} else {
						t.SetLowNibble(true)
					}

					if !e.Modifiers.Contain(key.ModShift) {
						t.SetSelectionCorner(t.Cursor())
					}
				case "+":
					if e.Modifiers.Contain(key.ModShortcut) {
						t.AdjustSelectionPitch(12)
					} else {
						t.AdjustSelectionPitch(1)
					}
				case "-":
					if e.Modifiers.Contain(key.ModShortcut) {
						t.AdjustSelectionPitch(-12)
					} else {
						t.AdjustSelectionPitch(-1)
					}
				}
				if e.Modifiers.Contain(key.ModShortcut) {
					continue
				}
				step := false
				if t.Song().Score.Tracks[t.Cursor().Track].Effect {
					if iv, err := strconv.ParseInt(e.Name, 16, 8); err == nil {
						t.NumberPressed(byte(iv))
						step = true
					}
				} else {
					if e.Name == "A" || e.Name == "1" {
						t.SetNote(0)
						step = true
					} else {
						if val, ok := noteMap[e.Name]; ok {
							if _, ok := t.KeyPlaying[e.Name]; !ok {
								n := tracker.NoteAsValue(t.OctaveNumberInput.Value, val)
								t.SetNote(n)
								step = true
								trk := t.Cursor().Track
								noteID := tracker.NoteIDTrack(trk, n)
								t.NoteOn(noteID)
								t.KeyPlaying[e.Name] = noteID
							}
						}
					}
				}
				if step && !(t.NoteTracking() && t.Playing()) && t.Step.Value > 0 {
					t.SetCursor(t.Cursor().AddRows(t.Step.Value))
					t.SetSelectionCorner(t.Cursor())
				}

				t.JammingPressed(e)
			case key.Release:
				t.JammingReleased(e)
			}
		}
	}

	if te.requestFocus {
		te.requestFocus = false
		key.FocusOp{Tag: &te.tag}.Add(gtx.Ops)
	}

	rowMarkers := layout.Rigid(t.layoutRowMarkers)

	for te.NewTrackBtn.Clicked() {
		t.AddTrack(true)
	}

	for te.DeleteTrackBtn.Clicked() {
		t.DeleteTrack(false)
	}

	//t.TrackHexCheckBoxes[i2].Value = t.TrackShowHex[i2]
	//cbStyle := material.CheckBox(t.Theme, t.TrackHexCheckBoxes[i2], "hex")
	//cbStyle.Color = white
	//cbStyle.IconColor = t.Theme.Fg

	for te.AddSemitoneBtn.Clicked() {
		t.AdjustSelectionPitch(1)
	}

	for te.SubtractSemitoneBtn.Clicked() {
		t.AdjustSelectionPitch(-1)
	}

	for te.NoteOffBtn.Clicked() {
		t.SetNote(0)
		if !(t.NoteTracking() && t.Playing()) && t.Step.Value > 0 {
			t.SetCursor(t.Cursor().AddRows(t.Step.Value))
			t.SetSelectionCorner(t.Cursor())
		}
	}

	for te.AddOctaveBtn.Clicked() {
		t.AdjustSelectionPitch(12)
	}

	for te.SubtractOctaveBtn.Clicked() {
		t.AdjustSelectionPitch(-12)
	}

	menu := func(gtx C) D {
		addSemitoneBtnStyle := LowEmphasisButton(t.Theme, te.AddSemitoneBtn, "+1")
		subtractSemitoneBtnStyle := LowEmphasisButton(t.Theme, te.SubtractSemitoneBtn, "-1")
		addOctaveBtnStyle := LowEmphasisButton(t.Theme, te.AddOctaveBtn, "+12")
		subtractOctaveBtnStyle := LowEmphasisButton(t.Theme, te.SubtractOctaveBtn, "-12")
		noteOffBtnStyle := LowEmphasisButton(t.Theme, te.NoteOffBtn, "Note Off")
		deleteTrackBtnStyle := IconButton(t.Theme, te.DeleteTrackBtn, icons.ActionDelete, t.CanDeleteTrack())
		newTrackBtnStyle := IconButton(t.Theme, te.NewTrackBtn, icons.ContentAdd, t.CanAddTrack())
		n := t.Song().Score.Tracks[t.Cursor().Track].NumVoices
		te.TrackVoices.Value = n
		in := layout.UniformInset(unit.Dp(1))
		voiceUpDown := func(gtx C) D {
			numStyle := NumericUpDown(t.Theme, te.TrackVoices, 1, t.MaxTrackVoices())
			gtx.Constraints.Min.Y = gtx.Px(unit.Dp(20))
			gtx.Constraints.Min.X = gtx.Px(unit.Dp(70))
			return in.Layout(gtx, numStyle.Layout)
		}
		t.TrackHexCheckBox.Value = t.Song().Score.Tracks[t.Cursor().Track].Effect
		hexCheckBoxStyle := material.CheckBox(t.Theme, t.TrackHexCheckBox, "Hex")
		dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D { return layout.Dimensions{Size: image.Pt(gtx.Px(unit.Dp(12)), 0)} }),
			layout.Rigid(addSemitoneBtnStyle.Layout),
			layout.Rigid(subtractSemitoneBtnStyle.Layout),
			layout.Rigid(addOctaveBtnStyle.Layout),
			layout.Rigid(subtractOctaveBtnStyle.Layout),
			layout.Rigid(noteOffBtnStyle.Layout),
			layout.Rigid(hexCheckBoxStyle.Layout),
			layout.Rigid(Label("  Voices:", white)),
			layout.Rigid(voiceUpDown),
			layout.Flexed(1, func(gtx C) D { return layout.Dimensions{Size: gtx.Constraints.Min} }),
			layout.Rigid(deleteTrackBtnStyle.Layout),
			layout.Rigid(newTrackBtnStyle.Layout))
		t.Song().Score.Tracks[t.Cursor().Track].Effect = t.TrackHexCheckBox.Value // TODO: we should not modify the model, but how should this be done
		t.SetTrackVoices(te.TrackVoices.Value)
		return dims
	}

	rect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	pointer.Rect(rect).Add(gtx.Ops)
	pointer.InputOp{Tag: &te.tag,
		Types: pointer.Press,
	}.Add(gtx.Ops)
	key.InputOp{Tag: &te.tag}.Add(gtx.Ops)

	return Surface{Gray: 24, Focus: te.focused}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return Surface{Gray: 37, Focus: te.focused, FitSize: true}.Layout(gtx, menu)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					rowMarkers,
					layout.Flexed(1, func(gtx C) D {
						return te.layoutTracks(gtx, t)
					}))
			}),
		)
	})
}

func (te *TrackEditor) layoutTracks(gtx C, t *Tracker) D {
	defer op.Save(gtx.Ops).Load()
	clip.Rect{Max: gtx.Constraints.Max}.Add(gtx.Ops)
	cursorSongRow := t.Cursor().Pattern*t.Song().Score.RowsPerPattern + t.Cursor().Row
	for _, ev := range gtx.Events(&te.trackJumpPointerTag) {
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if e.Type == pointer.Press {
			te.Focus()
			track := int(e.Position.X) / trackColWidth
			row := int((e.Position.Y-float32(gtx.Constraints.Max.Y-trackRowHeight)/2)/trackRowHeight + float32(cursorSongRow))
			cursor := tracker.SongPoint{Track: track, SongRow: tracker.SongRow{Row: row}}.Clamp(t.Song().Score)
			t.SetCursor(cursor)
			t.SetSelectionCorner(cursor)
			cursorSongRow = cursor.Pattern*t.Song().Score.RowsPerPattern + cursor.Row
		}
	}
	rect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	pointer.Rect(rect).Add(gtx.Ops)
	pointer.InputOp{Tag: &te.trackJumpPointerTag,
		Types: pointer.Press,
	}.Add(gtx.Ops)
	stack := op.Save(gtx.Ops)
	curVoice := 0
	for _, trk := range t.Song().Score.Tracks {
		gtx := gtx
		instrName := "?"
		firstIndex, err := t.Song().Patch.InstrumentForVoice(curVoice)
		lastIndex, err2 := t.Song().Patch.InstrumentForVoice(curVoice + trk.NumVoices - 1)
		if err == nil && err2 == nil {
			switch diff := lastIndex - firstIndex; diff {
			case 0:
				instrName = t.Song().Patch[firstIndex].Name
			default:
				n1 := t.Song().Patch[firstIndex].Name
				n2 := t.Song().Patch[firstIndex+1].Name
				if len(n1) > 0 {
					n1 = string(n1[0])
				} else {
					n1 = "?"
				}
				if len(n2) > 0 {
					n2 = string(n2[0])
				} else {
					n2 = "?"
				}
				if diff > 1 {
					instrName = n1 + "/" + n2 + "..."
				} else {
					instrName = n1 + "/" + n2
				}
			}
			if len(instrName) > 7 {
				instrName = instrName[:7]
			}
		}
		gtx.Constraints.Max.X = trackColWidth
		LabelStyle{Alignment: layout.N, Text: instrName, FontSize: unit.Dp(12), Color: mediumEmphasisTextColor}.Layout(gtx)
		op.Offset(f32.Pt(trackColWidth, 0)).Add(gtx.Ops)
		curVoice += trk.NumVoices
	}
	stack.Load()
	op.Offset(f32.Pt(0, float32(gtx.Constraints.Max.Y-trackRowHeight)/2)).Add(gtx.Ops)
	op.Offset(f32.Pt(0, (-1*trackRowHeight)*float32(cursorSongRow))).Add(gtx.Ops)
	if te.focused || t.OrderEditor.Focused() {
		x1, y1 := t.Cursor().Track, t.Cursor().Pattern
		x2, y2 := t.SelectionCorner().Track, t.SelectionCorner().Pattern
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		x2++
		y2++
		x1 *= trackColWidth
		y1 *= trackRowHeight * t.Song().Score.RowsPerPattern
		x2 *= trackColWidth
		y2 *= trackRowHeight * t.Song().Score.RowsPerPattern
		paint.FillShape(gtx.Ops, inactiveSelectionColor, clip.Rect{Min: image.Pt(x1, y1), Max: image.Pt(x2, y2)}.Op())
	}
	if te.focused {
		x1, y1 := t.Cursor().Track, t.Cursor().Pattern*t.Song().Score.RowsPerPattern+t.Cursor().Row
		x2, y2 := t.SelectionCorner().Track, t.SelectionCorner().Pattern*t.Song().Score.RowsPerPattern+t.SelectionCorner().Row
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		x2++
		y2++
		x1 *= trackColWidth
		y1 *= trackRowHeight
		x2 *= trackColWidth
		y2 *= trackRowHeight
		paint.FillShape(gtx.Ops, selectionColor, clip.Rect{Min: image.Pt(x1, y1), Max: image.Pt(x2, y2)}.Op())
		cx := t.Cursor().Track * trackColWidth
		cy := (t.Cursor().Pattern*t.Song().Score.RowsPerPattern + t.Cursor().Row) * trackRowHeight
		cw := trackColWidth
		if t.Song().Score.Tracks[t.Cursor().Track].Effect {
			cw /= 2
			if t.LowNibble() {
				cx += cw
			}
		}
		paint.FillShape(gtx.Ops, cursorColor, clip.Rect{Min: image.Pt(cx, cy), Max: image.Pt(cx+cw, cy+trackRowHeight)}.Op())
	}
	delta := (gtx.Constraints.Max.Y/2 + trackRowHeight - 1) / trackRowHeight
	firstRow := cursorSongRow - delta
	lastRow := cursorSongRow + delta
	if firstRow < 0 {
		firstRow = 0
	}
	if l := t.Song().Score.LengthInRows(); lastRow >= l {
		lastRow = l - 1
	}
	op.Offset(f32.Pt(0, float32(trackRowHeight*firstRow))).Add(gtx.Ops)
	for trkIndex, trk := range t.Song().Score.Tracks {
		stack := op.Save(gtx.Ops)
		for row := firstRow; row <= lastRow; row++ {
			pat := row / t.Song().Score.RowsPerPattern
			patRow := row % t.Song().Score.RowsPerPattern
			s := trk.Order.Get(pat)
			if s < 0 {
				op.Offset(f32.Pt(0, trackRowHeight)).Add(gtx.Ops)
				continue
			}
			if s >= 0 && patRow == 0 {
				paint.ColorOp{Color: trackerPatMarker}.Add(gtx.Ops)
				widget.Label{}.Layout(gtx, textShaper, trackerFont, trackerFontSize, patternIndexToString(s))
			}
			if s >= 0 && patRow == 1 && t.IsPatternUnique(trkIndex, s) {
				paint.ColorOp{Color: mediumEmphasisTextColor}.Add(gtx.Ops)
				widget.Label{}.Layout(gtx, textShaper, trackerFont, trackerFontSize, "*")
			}
			op.Offset(f32.Pt(patmarkWidth, 0)).Add(gtx.Ops)
			if te.focused && t.Cursor().Row == patRow && t.Cursor().Pattern == pat {
				paint.ColorOp{Color: trackerActiveTextColor}.Add(gtx.Ops)
			} else {
				paint.ColorOp{Color: trackerInactiveTextColor}.Add(gtx.Ops)
			}
			var c byte = 1
			if s >= 0 && s < len(trk.Patterns) {
				c = trk.Patterns[s].Get(patRow)
			}
			if trk.Effect {
				var text string
				switch c {
				case 0:
					text = "--"
				case 1:
					text = ".."
				default:
					text = fmt.Sprintf("%02x", c)
				}
				widget.Label{}.Layout(gtx, textShaper, trackerFont, trackerFontSize, strings.ToUpper(text))
			} else {
				widget.Label{}.Layout(gtx, textShaper, trackerFont, trackerFontSize, tracker.NoteStr(c))
			}
			op.Offset(f32.Pt(-patmarkWidth, trackRowHeight)).Add(gtx.Ops)
		}
		stack.Load()
		op.Offset(f32.Pt(trackColWidth, 0)).Add(gtx.Ops)
	}
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
