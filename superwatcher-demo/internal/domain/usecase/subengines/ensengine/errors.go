package ensengine

import "github.com/pkg/errors"

var (
	ErrLogLen = errors.New("invalid log topics length")
	ErrMapENS = errors.New("error mapping event log to entity.ENS")
)
