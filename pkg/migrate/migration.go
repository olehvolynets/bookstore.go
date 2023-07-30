package migrate

import (
	"fmt"
	"strconv"
	"time"
)

var (
	upDirection   = "up"
	downDirection = "down"
	printFormat   = `{
    Version: %s
    Name: %s
    UpFile: %s
    DownFile: %s
    AppliedAt: %v
}`
)

type Migration struct {
	Version   string
	Name      string
	UpFile    string
	DownFile  string
	AppliedAt time.Time
}

func (m *Migration) VersionUint() uint64 {
	v, _ := strconv.ParseUint(m.Version, 10, 64)

	return v
}

func (m *Migration) String() string {
	return fmt.Sprintf(printFormat, m.Version, m.Name, m.UpFile, m.DownFile, m.AppliedAt)
}

type Migrations []Migration

func (migs Migrations) Less(i, j int) bool {
	return migs[i].VersionUint() < migs[j].VersionUint()
}

func (migs Migrations) Len() int { return len(migs) }

func (migs Migrations) Swap(i, j int) {
	migs[i], migs[j] = migs[j], migs[i]
}
