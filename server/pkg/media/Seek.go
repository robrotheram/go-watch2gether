package media

type Seek struct {
	ProgressPct float64 `json:"progress_percent"`
	ProgressSec float64 `json:"progress_seconds"`
}

func (s *Seek) Done() bool {
	return s.ProgressPct == float64(1)
}

var SEEK_FINISHED = Seek{ProgressPct: float64(1), ProgressSec: 0}
var SEEK_INIT = Seek{ProgressPct: float64(0), ProgressSec: 0}
