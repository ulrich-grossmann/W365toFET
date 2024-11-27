package base

const MAXWEIGHT = 100

func (db *DbTopLevel) addConstraint(c Constraint) {
	db.Constraints = append(db.Constraints, c)
}

// ++ LessonsEndDay

type LessonsEndDay struct {
	Constraint string
	Weight     int
	Course     Ref
}

func (c *LessonsEndDay) CType() string {
	return c.Constraint
}

func (db *DbTopLevel) NewLessonsEndDay() *LessonsEndDay {
	c := &LessonsEndDay{Constraint: "LessonsEndDay"}
	db.addConstraint(c)
	return c
}

// ++ BeforeAfterHour
// Permissible hours are before or after the specified hour, not including
// the specified hour.

type BeforeAfterHour struct {
	Constraint string
	Weight     int
	Courses    []Ref
	After      bool
	Hour       int
}

func (c *BeforeAfterHour) CType() string {
	return c.Constraint
}

func (db *DbTopLevel) NewBeforeAfterHour() *BeforeAfterHour {
	c := &BeforeAfterHour{Constraint: "BeforeAfterHour"}
	db.addConstraint(c)
	return c
}

// ++ AutomaticDifferentDays
// This Constraint applies to all courses (with more than one Lesson).
// If not present, all courses will by default apply it as a hard constraint,
// except for courses which have an overriding DAYS_BETWEEN constraint.

type AutomaticDifferentDays struct {
	Constraint           string
	Weight               int
	ConsecutiveIfSameDay bool
}

func (c *AutomaticDifferentDays) CType() string {
	return c.Constraint
}

func (db *DbTopLevel) NewAutomaticDifferentDays() *AutomaticDifferentDays {
	c := &AutomaticDifferentDays{Constraint: "AutomaticDifferentDays"}
	db.addConstraint(c)
	return c
}

// ++ DaysBetween
// This constraint applies between the lessons of the individual courses.
// It does not connect the courses. If DaysBetween = 1, this constraint
// overrides the global AutomaticDifferentDays constraint for these courses.

type DaysBetween struct {
	Constraint           string
	Weight               int
	Courses              []Ref // Courses or SuperCourses
	DaysBetween          int
	ConsecutiveIfSameDay bool
}

func (c *DaysBetween) CType() string {
	return c.Constraint
}

func (db *DbTopLevel) NewDaysBetween() *DaysBetween {
	c := &DaysBetween{Constraint: "DaysBetween"}
	db.addConstraint(c)
	return c
}

// ++ DaysBetweenJoin
// This constraint applies between the individual lessons of the two courses,
// not between the lessons of a course itself. That is, between course 1,
// lesson 1 and course 2 lesson 1; between course 1, lesson 1 and course 2,
// lesson 2, etc.

type DaysBetweenJoin struct {
	Constraint           string
	Weight               int
	Course1              Ref // Course or SuperCourse
	Course2              Ref // Course or SuperCourse
	DaysBetween          int
	ConsecutiveIfSameDay bool
}

func (c *DaysBetweenJoin) CType() string {
	return c.Constraint
}

func (db *DbTopLevel) NewDaysBetweenJoin() *DaysBetweenJoin {
	c := &DaysBetweenJoin{Constraint: "DaysBetweenJoin"}
	db.addConstraint(c)
	return c
}

//TODO ...

// ++ NotOnSameDay

type NotOnSameDay struct {
	Constraint string
	Weight     int
	Subjects   []Ref
}

func (c *NotOnSameDay) CType() string {
	return c.Constraint
}

func (db *DbTopLevel) NewNotOnSameDay() *NotOnSameDay {
	c := &NotOnSameDay{Constraint: "NotOnSameDay"}
	db.addConstraint(c)
	return c
}

// ++ ParallelCourses
// The lessons of the courses specified here should be at the same time.
// To avoid complications, it is required that the number and lengths of
// lessons be the same in each course.

type ParallelCourses struct {
	Constraint string
	Weight     int
	Courses    []Ref // Courses or SuperCourses
}

func (c *ParallelCourses) CType() string {
	return c.Constraint
}

func (db *DbTopLevel) NewParallelCourses() *ParallelCourses {
	c := &ParallelCourses{Constraint: "ParallelCourses"}
	db.addConstraint(c)
	return c
}
