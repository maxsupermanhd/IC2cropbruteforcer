package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type crop struct {
	Humidity       float64
	Nutrients      float64
	Air            float64
	Growth         int
	Gain           int
	Resist         int
	Tier           int
	GrowthDuration int
}

type cropEnvironment struct {
	Humidity  int
	Nutrients int
	Air       int
}

type fieldCell struct {
	What        string
	Environment cropEnvironment `json:"-"`
	GrowPoints  int
}

type cropField struct {
	Width                   int
	Height                  int
	widget                  *gtk.Box
	Cells                   []fieldCell
	CropStat                crop `json:"-"`
	CropYlevel              int  `json:"-"`
	CropSkyAccess           bool `json:"-"`
	CropWatering            bool `json:"-"`
	CropFertilizing         bool `json:"-"`
	CropDirtBelow           int  `json:"-"`
	CropBiomeHumidityBonus  int  `json:"-"`
	CropBiomeNutrientsBonus int  `json:"-"`
}

func (f *cropField) ChangeDimensions(w, h int) {
	newcells := make([]fieldCell, w*h)
	for i := 0; i < f.Height; i++ {
		for j := 0; j < f.Width; j++ {
			if i < w && j < h {
				newcells[i*w+j] = f.Cells[i*f.Height+j]
			}
		}
	}
	f.Cells = newcells
	f.Width = w
	f.Height = h
	f.widget.GetChildren().Foreach(func(item interface{}) { item.(*gtk.Widget).Destroy() })
	f.widget.Add(f.CreateWidget())
	f.widget.ShowAll()
}

func (f *cropField) CreateWidget() *gtk.Box {
	box := noerr(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))
	label := noerr(gtk.LabelNew(fmt.Sprintf("Field size: %dx%d (%d)", f.Width, f.Height, len(f.Cells))))
	// label.SetHExpand(true)
	btnadd := noerr(gtk.ButtonNewWithLabel("+"))
	btnadd.Connect("clicked", func() {
		f.ChangeDimensions(f.Width+1, f.Height+1)
		label.SetText(fmt.Sprintf("Field size: %dx%d (%d)", f.Width, f.Height, len(f.Cells)))
	})
	btnsub := noerr(gtk.ButtonNewWithLabel("-"))
	btnsub.Connect("clicked", func() {
		f.ChangeDimensions(f.Width-1, f.Height-1)
		label.SetText(fmt.Sprintf("Field size: %dx%d (%d)", f.Width, f.Height, len(f.Cells)))
	})
	labelgrid := noerr(gtk.GridNew())
	labelgrid.Add(label)
	labelgrid.Add(btnadd)
	labelgrid.Add(btnsub)
	box.Add(labelgrid)

	grid := noerr(gtk.GridNew())
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	for row := 0; row < f.Height; row++ {
		for col := 0; col < f.Width; col++ {
			button := noerr(gtk.ButtonNew())
			if f.Cells[row*f.Height+col].What == "air" {
				button.SetName("cropField-air")
			} else if f.Cells[row*f.Height+col].What == "block" {
				button.SetName("cropField-block")
			} else {
				button.SetName("cropField-crop")
			}
			torow := row
			tocol := col
			button.Connect("clicked", func() {
				name := noerr(button.GetName())
				if name == "cropField-air" {
					button.SetName("cropField-crop")
					f.Cells[torow*f.Height+tocol].What = "crop"
				} else if name == "cropField-crop" {
					button.SetName("cropField-block")
					f.Cells[torow*f.Height+tocol].What = "block"
				} else if name == "cropField-block" {
					button.SetName("cropField-air")
					f.Cells[torow*f.Height+tocol].What = "air"
				}
			})
			grid.Attach(button, row, col, 1, 1)
		}
	}
	box.Add(grid)
	return box
}

func NewCropField(s int) *cropField {
	f := cropField{Width: s, Height: s}
	box := noerr(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	f.widget = box
	f.Cells = make([]fieldCell, s*s)
	for i := 0; i < s*s; i++ {
		f.Cells[i].What = "crop"
	}
	f.widget.Add(f.CreateWidget())
	return &f
}

func main() {
	gtk.Init(nil)
	win := noerr(gtk.WindowNew(gtk.WINDOW_TOPLEVEL))
	win.SetTitle("IC2 ctop bruteforcer")
	win.SetDefaultSize(800, 700)
	win.Connect("destroy", func() { gtk.MainQuit() })

	mainbox := noerr(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 4))
	mainbox.SetVExpand(false)
	mainbox.SetHExpand(true)
	row0 := noerr(gtk.GridNew())
	row0.SetVExpand(false)
	row0.SetHExpand(true)
	row1 := noerr(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4))
	row1.SetVExpand(false)
	row1.SetHExpand(true)
	mainbox.Add(row0)
	mainbox.Add(row1)

	field := NewCropField(13)
	row0.Add(field.widget)

	envParams := noerr(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 3))
	envParams.Add(noerr(gtk.LabelNew("Field parameters:")))
	envParams.Add(intSelectorWithLabel(" Crop Y level", 0, 256, 1, func(newval int) {
		field.CropYlevel = newval
	}))
	envParams.Add(boolSelectorWithLabel("Sky access", func(newval bool) {
		field.CropSkyAccess = newval
	}))
	envParams.Add(boolSelectorWithLabel("Water supplied", func(newval bool) {
		field.CropWatering = newval
	}))
	envParams.Add(intSelectorWithLabel(" Additional dirt below", 0, 5, 1, func(newval int) {
		field.CropDirtBelow = newval
	}))
	envParams.Add(intSelectorWithLabel(" Biome humidity bonus", 0, 100, 1, func(newval int) {
		field.CropBiomeHumidityBonus = newval
	}))
	envParams.Add(intSelectorWithLabel(" Biome nutrients bonus", 0, 100, 1, func(newval int) {
		field.CropBiomeNutrientsBonus = newval
	}))
	row1.Add(envParams)

	cropParams := noerr(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 3))
	cropParams.Add(noerr(gtk.LabelNew("Crop parameters:")))
	cropParams.Add(floatSelectorWithLabel("Humidity modifier", 0, 5, 0.1, func(newval float64) {
		field.CropStat.Humidity = newval
	}))
	cropParams.Add(floatSelectorWithLabel("Nutrients modifier", 0, 5, 0.1, func(newval float64) {
		field.CropStat.Nutrients = newval
	}))
	cropParams.Add(floatSelectorWithLabel("Air modifier", 0, 5, 0.1, func(newval float64) {
		field.CropStat.Air = newval
	}))
	cropParams.Add(intSelectorWithLabel("Tier", 0, 255, 1, func(newval int) {
		field.CropStat.Tier = newval
	}))
	cropParams.Add(intSelectorWithLabel("Growth", 0, 255, 1, func(newval int) {
		field.CropStat.Growth = newval
	}))
	cropParams.Add(intSelectorWithLabel("Gain", 0, 255, 1, func(newval int) {
		field.CropStat.Gain = newval
	}))
	cropParams.Add(intSelectorWithLabel("Resist", 0, 255, 1, func(newval int) {
		field.CropStat.Resist = newval
	}))
	cropParams.Add(intSelectorWithLabel("Growth duration", 0, 200000, 1, func(newval int) {
		field.CropStat.GrowthDuration = newval
	}))
	row1.Add(cropParams)
	calcResults := noerr(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 3))
	calcResults.Add(noerr(gtk.LabelNew("Calculation results:")))
	calcAvgPoints := noerr(gtk.LabelNew("Not yet calculated"))
	calcResults.Add(calcAvgPoints)
	row1.Add(calcResults)

	cGrid := noerr(gtk.GridNew())
	row0.Add(cGrid)
	for row := 0; row < 13; row++ {
		for col := 0; col < 13; col++ {
			padding := noerr(gtk.LabelNew(""))
			padding.SetHExpand(true)
			padding.SetVExpand(true)
			cGrid.Attach(padding, row, col, 1, 1)
		}
	}
	cTrigger := noerr(gtk.ButtonNewWithLabel("Calculate"))
	cTrigger.Connect("clicked", func() {
		calcFieldEnvironment(field)
		report := fmt.Sprintf("Average growth per random tick: %.2f\n", calcFieldAverageGrouth(field))
		times, total := calcFieldGrowthTime(field)
		for k, v := range times {
			report += fmt.Sprintf("%d/%d will grow in %.2f random ticks\n", v, total, k)
		}
		calcAvgPoints.SetText(report)
	})
	cGrid.Attach(cTrigger, 6, 6, 1, 1)

	cSave := noerr(gtk.ButtonNewWithLabel("Save"))
	cSave.Connect("clicked", func() {
		d := noerr(gtk.FileChooserDialogNewWith2Buttons("Save file", win, gtk.FILE_CHOOSER_ACTION_SAVE, "Cancel", gtk.RESPONSE_CANCEL, "Save", gtk.RESPONSE_ACCEPT))
		if d.Run() == gtk.RESPONSE_ACCEPT {
			must(os.WriteFile(d.GetFilename(), noerr(json.Marshal(field)), 0666))
		}
		d.Close()
	})
	cGrid.Attach(cSave, 4, 4, 1, 1)

	cLoad := noerr(gtk.ButtonNewWithLabel("Load"))
	cLoad.Connect("clicked", func() {
		d := noerr(gtk.FileChooserDialogNewWith2Buttons("Open file", win, gtk.FILE_CHOOSER_ACTION_OPEN, "Cancel", gtk.RESPONSE_CANCEL, "Load", gtk.RESPONSE_ACCEPT))
		if d.Run() == gtk.RESPONSE_ACCEPT {
			must(json.Unmarshal(noerr(os.ReadFile(d.GetFilename())), &field))
			field.ChangeDimensions(field.Width, field.Height)
		}
		d.Close()
	})
	cGrid.Attach(cLoad, 8, 4, 1, 1)

	win.Add(mainbox)
	win.ShowAll()

	cssWdgScnBytes([]byte(`
	#cropField-air {
		background-color: white;
		background-image: none;
		color: black;
	}
	#cropField-crop {
		background-color: green;
		background-image: none;
	}
	#cropField-block {
		background-color: black;
		background-image: none;
		color: white;
	}`))

	gtk.Main()
}

func cssWdgScnBytes(data []byte) {
	cssProv := noerr(gtk.CssProviderNew())
	must(cssProv.LoadFromData(string(data)))
	gtk.AddProviderForScreen(noerr(gdk.ScreenGetDefault()), cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func intSelectorWithLabel(label string, min, max, step float64, f func(int)) *gtk.Grid {
	g := noerr(gtk.GridNew())
	l := noerr(gtk.LabelNew(label))
	s := noerr(gtk.SpinButtonNewWithRange(min, max, step))
	s.Connect("value-changed", func() { f(s.GetValueAsInt()) })
	g.Add(s)
	g.Add(l)
	return g
}

func floatSelectorWithLabel(label string, min, max, step float64, f func(float64)) *gtk.Grid {
	g := noerr(gtk.GridNew())
	l := noerr(gtk.LabelNew(label))
	s := noerr(gtk.SpinButtonNewWithRange(min, max, step))
	s.Connect("value-changed", func() { f(s.GetValue()) })
	g.Add(s)
	g.Add(l)
	return g
}

func boolSelectorWithLabel(label string, f func(bool)) *gtk.CheckButton {
	s := noerr(gtk.CheckButtonNewWithLabel(label))
	s.Connect("toggled", func() { f(s.GetActive()) })
	return s
}
