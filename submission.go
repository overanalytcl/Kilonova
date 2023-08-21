package kilonova

import (
	"time"

	"github.com/shopspring/decimal"
)

const MaxScoreRoundingPlaces = 4

type Status string

const (
	StatusNone     Status = ""
	StatusCreating Status = "creating"
	StatusWaiting  Status = "waiting"
	StatusWorking  Status = "working"
	StatusFinished Status = "finished"
)

type Submission struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
	ProblemID int       `json:"problem_id"`
	Language  string    `json:"language"`
	Code      string    `json:"code,omitempty"`
	CodeSize  int       `json:"code_size"`
	Status    Status    `json:"status"`

	CompileError   *bool   `json:"compile_error"`
	CompileMessage *string `json:"compile_message,omitempty"`

	ContestID *int `json:"contest_id"`

	MaxTime   float64 `json:"max_time"`
	MaxMemory int     `json:"max_memory"`

	Score decimal.Decimal `json:"score"`

	ScorePrecision int32 `json:"score_precision"`
}

type SubmissionUpdate struct {
	Status Status
	Score  *decimal.Decimal

	CompileError   *bool
	CompileMessage *string

	MaxTime   *float64
	MaxMemory *int
}

type SubmissionFilter struct {
	ID        *int   `json:"id"`
	IDs       []*int `json:"ids"`
	UserID    *int   `json:"user_id"`
	ProblemID *int   `json:"problem_id"`
	ContestID *int   `json:"contest_id"`

	Status       Status  `json:"status"`
	Lang         *string `json:"lang"`
	CompileError *bool   `json:"compile_error"`

	Score *decimal.Decimal `json:"score"`

	Look        bool       `json:"-"`
	LookingUser *UserBrief `json:"-"`

	FromAuthors bool `json:"from_authors"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`

	Ordering  string `json:"ordering"`
	Ascending bool   `json:"ascending"`
}

type SubTest struct {
	ID           int             `json:"id"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	Done         bool            `json:"done"`
	Verdict      string          `json:"verdict"`
	Time         float64         `json:"time"`
	Memory       int             `json:"memory"`
	Percentage   decimal.Decimal `json:"percentage"`
	TestID       *int            `db:"test_id" json:"test_id"`
	UserID       int             `db:"user_id" json:"user_id"`
	SubmissionID int             `db:"submission_id" json:"submission_id"`

	ContestID *int `db:"contest_id" json:"contest_id"`

	VisibleID int `db:"visible_id" json:"visible_id"`

	Score decimal.Decimal `json:"score"`
}

type SubTestUpdate struct {
	Memory     *int
	Time       *float64
	Percentage *decimal.Decimal
	Verdict    *string
	Done       *bool
}

type SubmissionSubTask struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	SubmissionID int  `json:"submission_id"`
	UserID       int  `json:"user_id"`
	SubtaskID    *int `json:"subtask_id"`

	ProblemID int  `json:"problem_id"`
	ContestID *int `json:"contest_id"`
	VisibleID int  `json:"visible_id"`

	Score           decimal.Decimal  `json:"score"`
	FinalPercentage *decimal.Decimal `json:"final_percentage,omitempty"`

	ScorePrecision int `json:"score_precision"`

	Subtests []int `json:"subtests"`
}

type SubmissionPaste struct {
	ID         string      `json:"id"`
	Submission *Submission `json:"sub"`
	Author     *UserBrief  `json:"author"`
}

type FullSubmission struct {
	Submission
	Author   *UserBrief `json:"author"`
	Problem  *Problem   `json:"problem"`
	SubTests []*SubTest `json:"subtests"`

	SubTasks []*SubmissionSubTask `json:"subtasks"`

	// ProblemEditor returns whether the looking user is a problem editor
	ProblemEditor bool `json:"problem_editor"`

	CodeTrulyVisible bool `json:"truly_visible"`
}
