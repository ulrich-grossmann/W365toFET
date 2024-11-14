package w365tt

import (
	"W365toFET/base"
	"W365toFET/logging"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// The structures used for the "database", adapted to read from W365
//TODO: Currently dealing only with the elements needed for the timetable

type Ref = base.Ref // Element reference

type Info struct {
	Institution        string `json:"schoolName"`
	FirstAfternoonHour int    `json:"firstAfternoonHour"`
	MiddayBreak        []int  `json:"middayBreak"`
	Reference          string `json:"scenario"`
}

type Day struct {
	Id   Ref    `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"shortcut"`
}

type Hour struct {
	Id    Ref    `json:"id"`
	Name  string `json:"name"`
	Tag   string `json:"shortcut"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type TimeSlot struct {
	Day  int `json:"day"`
	Hour int `json:"hour"`
}

type Teacher struct {
	Id               Ref        `json:"id"`
	Name             string     `json:"name"`
	Tag              string     `json:"shortcut"`
	Firstname        string     `json:"firstname"`
	NotAvailable     []TimeSlot `json:"absences"`
	MinLessonsPerDay int        `json:"minLessonsPerDay"`
	MaxLessonsPerDay int        `json:"maxLessonsPerDay"`
	MaxDays          int        `json:"maxDays"`
	MaxGapsPerDay    int        `json:"maxGapsPerDay"`
	MaxGapsPerWeek   int        `json:"maxGapsPerWeek"`
	MaxAfternoons    int        `json:"maxAfternoons"`
	LunchBreak       bool       `json:"lunchBreak"`
}

func (t *Teacher) UnmarshalJSON(data []byte) error {
	// Customize defaults for Teacher
	t.MinLessonsPerDay = -1
	t.MaxLessonsPerDay = -1
	t.MaxDays = -1
	t.MaxGapsPerDay = -1
	t.MaxGapsPerWeek = -1
	t.MaxAfternoons = -1

	type tempT Teacher
	return json.Unmarshal(data, (*tempT)(t))
}

type Subject struct {
	Id   Ref    `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"shortcut"`
}

type Room struct {
	Id           Ref        `json:"id"`
	Name         string     `json:"name"`
	Tag          string     `json:"shortcut"`
	NotAvailable []TimeSlot `json:"absences"`
}

type RoomGroup struct {
	Id    Ref    `json:"id"`
	Name  string `json:"name"`
	Tag   string `json:"shortcut"`
	Rooms []Ref  `json:"rooms"`
}

type RoomChoiceGroup struct {
	Id    Ref    `json:"id"`
	Name  string `json:"name"`
	Tag   string `json:"shortcut"`
	Rooms []Ref  `json:"rooms"`
}

type Class struct {
	Id               Ref        `json:"id"`
	Name             string     `json:"name"`
	Tag              string     `json:"shortcut"`
	Year             int        `json:"level"`
	Letter           string     `json:"letter"`
	NotAvailable     []TimeSlot `json:"absences"`
	Divisions        []Division `json:"divisions"`
	MinLessonsPerDay int        `json:"minLessonsPerDay"`
	MaxLessonsPerDay int        `json:"maxLessonsPerDay"`
	MaxGapsPerDay    int        `json:"maxGapsPerDay"`
	MaxGapsPerWeek   int        `json:"maxGapsPerWeek"`
	MaxAfternoons    int        `json:"maxAfternoons"`
	LunchBreak       bool       `json:"lunchBreak"`
	ForceFirstHour   bool       `json:"forceFirstHour"`
}

func (t *Class) UnmarshalJSON(data []byte) error {
	// Customize defaults for Teacher
	t.MinLessonsPerDay = -1
	t.MaxLessonsPerDay = -1
	t.MaxGapsPerDay = -1
	t.MaxGapsPerWeek = -1
	t.MaxAfternoons = -1

	type tempT Class
	return json.Unmarshal(data, (*tempT)(t))
}

type Group struct {
	Id  Ref    `json:"id"`
	Tag string `json:"shortcut"`
}

type Division struct {
	Id     Ref    `json:"id"`
	Name   string `json:"name"`
	Groups []Ref  `json:"groups"`
}

type Course struct {
	Id             Ref   `json:"id"`
	Subjects       []Ref `json:"subjects,omitempty"`
	Subject        Ref   `json:"subject"`
	Groups         []Ref `json:"groups"`
	Teachers       []Ref `json:"teachers"`
	PreferredRooms []Ref `json:"preferredRooms,omitempty"`
	// Not in W365:
	Room Ref // Room, RoomGroup or RoomChoiceGroup Element
}

type SuperCourse struct {
	Id        Ref `json:"id"`
	Subject   Ref `json:"subject"`
	EpochPlan Ref `json:"epochPlan,omitempty"`
}

type SubCourse struct {
	Id0            Ref   `json:"id"`
	Id             Ref   `json:"-"`
	SuperCourses   []Ref `json:"superCourses"`
	Subjects       []Ref `json:"subjects,omitempty"`
	Subject        Ref   `json:"subject"`
	Groups         []Ref `json:"groups"`
	Teachers       []Ref `json:"teachers"`
	PreferredRooms []Ref `json:"preferredRooms,omitempty"`
	// Not in W365:
	Room Ref // Room, RoomGroup or RoomChoiceGroup Element
}

type Lesson struct {
	Id       Ref   `json:"id"`
	Course   Ref   `json:"course"` // Course or SuperCourse Elements
	Duration int   `json:"duration"`
	Day      int   `json:"day"`
	Hour     int   `json:"hour"`
	Fixed    bool  `json:"fixed"`
	Rooms    []Ref `json:"localRooms"` // only Room Elements
}

type EpochPlan struct {
	Id   Ref    `json:"id"`
	Tag  string `json:"shortcut"`
	Name string `json:"name"`
}

type DbTopLevel struct {
	Info             Info               `json:"w365TT"`
	Days             []*Day             `json:"days"`
	Hours            []*Hour            `json:"hours"`
	Teachers         []*Teacher         `json:"teachers"`
	Subjects         []*Subject         `json:"subjects"`
	Rooms            []*Room            `json:"rooms"`
	RoomGroups       []*RoomGroup       `json:"roomGroups"`
	RoomChoiceGroups []*RoomChoiceGroup `json:"roomChoiceGroups"`
	Classes          []*Class           `json:"classes"`
	Groups           []*Group           `json:"groups"`
	Courses          []*Course          `json:"courses"`
	SuperCourses     []*SuperCourse     `json:"superCourses"`
	SubCourses       []*SubCourse       `json:"subCourses"`
	Lessons          []*Lesson          `json:"lessons"`
	EpochPlans       []*EpochPlan       `json:"epochPlans,omitempty"`
	Constraints      map[string]any     `json:"constraints"`

	// These fields do not belong in the JSON object.
	Elements        map[Ref]any       `json:"-"`
	MaxId           int               `json:"-"` // for "indexed" Ids only
	SubjectTags     map[string]Ref    `json:"-"`
	SubjectNames    map[string]string `json:"-"`
	RoomTags        map[string]Ref    `json:"-"`
	RoomChoiceNames map[string]Ref    `json:"-"`
}

func (db *DbTopLevel) NewId() Ref {
	return Ref(fmt.Sprintf("#%d", db.MaxId+1))
}

func (db *DbTopLevel) AddElement(ref Ref, element any) {
	_, nok := db.Elements[ref]
	if nok {
		logging.Error.Printf("Element Id defined more than once:\n  %s\n", ref)
		return
	}
	db.Elements[ref] = element
	// Special handling if it is an "indexed" Id.
	if strings.HasPrefix(string(ref), "#") {
		s := strings.TrimPrefix(string(ref), "#")
		i, err := strconv.Atoi(s)
		if err == nil {
			if i > db.MaxId {
				db.MaxId = i
			}
		}
	}
}

func (db *DbTopLevel) checkDb() {
	// Initializations
	if db.Info.MiddayBreak == nil {
		db.Info.MiddayBreak = []int{}
	} else {
		// Sort and check contiguity.
		slices.Sort(db.Info.MiddayBreak)
		mb := db.Info.MiddayBreak
		if mb[len(mb)-1]-mb[0] >= len(mb) {
			logging.Error.Fatalln("MiddayBreak hours not contiguous")
		}

	}
	db.SubjectTags = map[string]Ref{}
	db.SubjectNames = map[string]string{}
	db.RoomTags = map[string]Ref{}
	db.RoomChoiceNames = map[string]Ref{}
	// Initialize the Ref -> Element mapping
	db.Elements = make(map[Ref]any)

	// Checks
	if len(db.Days) == 0 {
		logging.Error.Fatalln("No Days")
	}
	if len(db.Hours) == 0 {
		logging.Error.Fatalln("No Hours")
	}
	if len(db.Teachers) == 0 {
		logging.Error.Fatalln("No Teachers")
	}
	if len(db.Subjects) == 0 {
		logging.Error.Fatalln("No Subjects")
	}
	if len(db.Rooms) == 0 {
		logging.Error.Fatalln("No Rooms")
	}
	if len(db.Classes) == 0 {
		logging.Error.Fatalln("No Classes")
	}

	// More initializations
	for _, n := range db.Days {
		db.AddElement(n.Id, n)
	}
	for _, n := range db.Hours {
		db.AddElement(n.Id, n)
	}
	for _, n := range db.Teachers {
		db.AddElement(n.Id, n)
	}
	for _, n := range db.Subjects {
		db.AddElement(n.Id, n)
	}
	for _, n := range db.Rooms {
		db.AddElement(n.Id, n)
	}
	for _, n := range db.Classes {
		db.AddElement(n.Id, n)
	}
	if db.RoomGroups == nil {
		db.RoomGroups = []*RoomGroup{}
	} else {
		for _, n := range db.RoomGroups {
			db.AddElement(n.Id, n)
		}
	}
	if db.RoomChoiceGroups == nil {
		db.RoomChoiceGroups = []*RoomChoiceGroup{}
	} else {
		for _, n := range db.RoomChoiceGroups {
			db.AddElement(n.Id, n)
		}
	}
	if db.Groups == nil {
		db.Groups = []*Group{}
	} else {
		for _, n := range db.Groups {
			db.AddElement(n.Id, n)
		}
	}
	if db.Courses == nil {
		db.Courses = []*Course{}
	} else {
		for _, n := range db.Courses {
			db.AddElement(n.Id, n)
		}
	}
	if db.SuperCourses == nil {
		db.SuperCourses = []*SuperCourse{}
	} else {
		for _, n := range db.SuperCourses {
			db.AddElement(n.Id, n)
		}
	}
	if db.SubCourses == nil {
		db.SubCourses = []*SubCourse{}
	} else {
		for _, n := range db.SubCourses {
			// Add a prefix to the Id to avoid possible clashes with a
			// Course having the same Id.
			nid := "$$" + n.Id0
			n.Id = nid
			db.AddElement(nid, n)
		}
	}
	if db.Lessons == nil {
		db.Lessons = []*Lesson{}
	} else {
		for _, n := range db.Lessons {
			db.AddElement(n.Id, n)
		}
	}
	if db.Constraints == nil {
		db.Constraints = make(map[string]any)
	}
}

func (dbp *DbTopLevel) newSubjectTag() string {
	// A rather primitive new-subject-tag generator
	i := 0
	for {
		i++
		tag := "X" + strconv.Itoa(i)
		_, nok := dbp.SubjectTags[tag]
		if !nok {
			return tag
		}
	}
}

func (dbp *DbTopLevel) makeNewSubject(tag, name string) Ref {
	stag := tag
	if stag == "" {
		stag = dbp.newSubjectTag()
	}

	sref := dbp.NewId()
	sbj := &Subject{
		Id:   sref,
		Tag:  stag,
		Name: name,
	}
	dbp.Subjects = append(dbp.Subjects, sbj)
	dbp.AddElement(sref, sbj)
	dbp.SubjectTags[stag] = sref
	if tag == "" && name != "" {
		dbp.SubjectNames[name] = stag
	}
	return sref
}

// Block all afternoons.
func (dbp *DbTopLevel) handleZeroAfternoons(notAvailable *[]TimeSlot) {
	// Make an array and fill this in two passes, then remake list
	namap := make([][]bool, len(dbp.Days))
	nhours := len(dbp.Hours)
	for i := range namap {
		namap[i] = make([]bool, nhours)
		for h := dbp.Info.FirstAfternoonHour; h < nhours; h++ {
			namap[i][h] = true
		}
	}
	for _, ts := range *notAvailable {
		namap[ts.Day][ts.Hour] = true
	}
	*notAvailable = []TimeSlot{}
	for d, naday := range namap {
		for h, nahour := range naday {
			if nahour {
				*notAvailable = append(*notAvailable, TimeSlot{d, h})
			}
		}
	}
}

// Interface for Course and SubCourse elements
type CourseInterface interface {
	GetId() Ref
	GetGroups() []Ref
	GetTeachers() []Ref
	GetSubject() Ref
	getSubjects() []Ref       // not available externally
	getPreferredRooms() []Ref // not available externally
	GetRoom() Ref
	setSubject(Ref)
	setSubjects([]Ref)
	setPreferredRooms([]Ref)
	setRoom(Ref)
}

func (c *Course) GetId() Ref                    { return c.Id }
func (c *SubCourse) GetId() Ref                 { return c.Id }
func (c *Course) GetGroups() []Ref              { return c.Groups }
func (c *SubCourse) GetGroups() []Ref           { return c.Groups }
func (c *Course) GetTeachers() []Ref            { return c.Teachers }
func (c *SubCourse) GetTeachers() []Ref         { return c.Teachers }
func (c *Course) GetSubject() Ref               { return c.Subject }
func (c *SubCourse) GetSubject() Ref            { return c.Subject }
func (c *Course) getSubjects() []Ref            { return c.Subjects }
func (c *SubCourse) getSubjects() []Ref         { return c.Subjects }
func (c *Course) getPreferredRooms() []Ref      { return c.PreferredRooms }
func (c *SubCourse) getPreferredRooms() []Ref   { return c.PreferredRooms }
func (c *Course) GetRoom() Ref                  { return c.Room }
func (c *SubCourse) GetRoom() Ref               { return c.Room }
func (c *Course) setSubject(r Ref)              { c.Subject = r }
func (c *SubCourse) setSubject(r Ref)           { c.Subject = r }
func (c *Course) setSubjects(rr []Ref)          { c.Subjects = rr }
func (c *SubCourse) setSubjects(rr []Ref)       { c.Subjects = rr }
func (c *Course) setPreferredRooms(rr []Ref)    { c.PreferredRooms = rr }
func (c *SubCourse) setPreferredRooms(rr []Ref) { c.PreferredRooms = rr }
func (c *Course) setRoom(r Ref)                 { c.Room = r }
func (c *SubCourse) setRoom(r Ref)              { c.Room = r }
