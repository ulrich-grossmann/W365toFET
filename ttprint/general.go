package ttprint

import (
	"W365toFET/base"
	"W365toFET/ttbase"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Tile struct {
	Day        int      `json:"day"`
	Hour       int      `json:"hour"`
	Duration   int      `json:"duration,omitempty"`
	Fraction   int      `json:"fraction,omitempty"`
	Offset     int      `json:"offset,omitempty"`
	Total      int      `json:"total,omitempty"`
	Subject    string   `json:"subject"`
	Groups     []string `json:"groups,omitempty"`
	Teachers   []string `json:"teachers,omitempty"`
	Rooms      []string `json:"rooms,omitempty"`
	Background string   `json:"background,omitempty"`
}

type Timetable struct {
	TableType string
	Info      map[string]any
	Typst     map[string]any `json:",omitempty"`
	Pages     []ttPage
}

type ttDay struct {
	Name  string
	Short string
}

type ttHour struct {
	Name  string
	Short string
	Start string
	End   string
}

type ttPage struct {
	Name       string
	Short      string
	Activities []Tile
}

func GenTypstData(
	ttinfo *ttbase.TtInfo,
	datadir string,
	stemfile string,
) []string {
	typst_files := []string{}
	printTables := ttinfo.Db.PrintOptions.PrintTables
	if len(printTables) == 0 {
		printTables = []string{
			"Class", "Teacher", "Room",
			"Class_overview", "Teacher_overview", "Room_overview",
		}
	}
	for _, ptable := range printTables {
		p, overview := strings.CutSuffix(ptable, "_overview")
		var f string
		switch p {
		case "Class":
			f = genTypstClassData(ttinfo, datadir, stemfile)
		case "Teacher":
			f = genTypstTeacherData(ttinfo, datadir, stemfile)
		case "Room":
			f = genTypstRoomData(ttinfo, datadir, stemfile)
		default:
			base.Error.Printf("\n", ptable)
			continue
		}
		typst_files = append(typst_files, f)
		if overview {
			typst_files = append(typst_files, f+"_overview")
		}
	}
	return typst_files
}

func makeTypstJson(tt Timetable, datadir string, outfile string) {
	b, err := json.MarshalIndent(tt, "", "  ")
	if err != nil {
		base.Error.Fatal(err)
	}
	// os.Stdout.Write(b)
	outdir := filepath.Join(datadir, "_data")
	if _, err := os.Stat(outdir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outdir, os.ModePerm)
		if err != nil {
			base.Error.Fatal(err)
		}
	}
	jsonpath := filepath.Join(outdir, outfile+".json")
	err = os.WriteFile(jsonpath, b, 0666)
	if err != nil {
		base.Error.Fatal(err)
	}
	base.Message.Printf("Wrote: %s\n", jsonpath)
}

func MakePdf(
	script string,
	datadir string,
	stemfile string,
	outfile string,
	typst string,
) {
	outdir := filepath.Join(datadir, "_pdf")
	if _, err := os.Stat(outdir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outdir, os.ModePerm)
		if err != nil {
			base.Error.Fatalln(err)
		}
	}
	outpath := filepath.Join(outdir, outfile+".pdf")

	cmd := exec.Command(typst, "compile",
		"--font-path", filepath.Join(datadir, "_fonts"),
		"--root", datadir,
		"--input", "ifile="+filepath.Join("/_data", stemfile+".json"),
		filepath.Join(datadir, "scripts", script),
		outpath)
	//fmt.Printf(" ::: %s\n", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		base.Error.Println("(Typst) " + string(output))
		base.Error.Fatal(err)
	}
	base.Message.Printf("Timetable written to: %s\n", outpath)
}
