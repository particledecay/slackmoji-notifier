package build

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"text/template"
)

var (
	Version = ""
	Date    = ""
)

const versionTpl = `
Version:    {{ .Version }}
SHA:        {{ .Commit }}
Built On:   {{ .Date }}`

func getStringOrNotAvailable(value string) string {
	if value != "" {
		return fmt.Sprintf("v%s", value)
	} else {
		return "(n/a)"
	}
}

func getCommitSHA() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && setting.Value != "" {
				return setting.Value
			}
		}
	}

	return ""
}

// PrintVersion outputs the short version info
func PrintVersion() {
	fmt.Println(getStringOrNotAvailable(Version))
}

// PrintLongVersion outputs the full version info
func PrintLongVersion() error {
	data := struct {
		Version string
		Commit  string
		Date    string
	}{
		Version: getStringOrNotAvailable(Version),
		Commit:  getCommitSHA(),
		Date:    getStringOrNotAvailable(Date),
	}

	var tpl bytes.Buffer

	t, err := template.New("build").Parse(versionTpl)
	if err != nil {
		return err
	}

	if err := t.Execute(&tpl, data); err != nil {
		return err
	}

	fmt.Println(tpl.String())
	return nil
}
