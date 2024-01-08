package runlevelx

type RunLevel uint

const (
	RunLevelUnknown RunLevel = iota
	RunLevelLocal
	RunLevelDev
	RunLevelEng
	RunLevelStg
	RunLevelPrd
	RunLevelCompose
	RunLevelTesting
)

var levelNames = map[RunLevel]string{
	RunLevelUnknown: "unknown",
	RunLevelLocal:   "local",
	RunLevelDev:     "dev",
	RunLevelEng:     "eng",
	RunLevelStg:     "stg",
	RunLevelPrd:     "prd",
	RunLevelCompose: "compose",
	RunLevelTesting: "testing",
}

func (rl RunLevel) RunLevelName() string {
	return levelNames[rl]
}

func (rl RunLevel) Is(testRl RunLevel) bool {
	return rl == testRl
}

func DetermineRunLevel(runLevelName string) RunLevel {
	runLevel := RunLevelLocal // default is local
	for k, v := range levelNames {
		if v == runLevelName {
			return k
		}
	}
	return runLevel
}

func (rl RunLevel) IsTesting() bool {
	return rl == RunLevelTesting
}

func (rl RunLevel) IsLocal() bool {
	return rl == RunLevelLocal
}

func (rl RunLevel) IsPrd() bool {
	return rl == RunLevelPrd
}

func (rl RunLevel) GetRunLevel() string {
	return rl.RunLevelName()
}
