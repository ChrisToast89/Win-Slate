package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChrisToast89/slate-for-windows/setup/internal/audit"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/install"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/manifest"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/paths"
	"github.com/ChrisToast89/slate-for-windows/setup/internal/update"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: qainstall <installDir> <payloadSlate.exe>")
		os.Exit(2)
	}
	dest := os.Args[1]
	raw, err := os.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	fmt.Println("=== AUDIT ===")
	r := audit.Run()
	fmt.Printf("canProceed=%v summary=%s\n", r.CanProceed, r.Summary)
	for _, c := range r.Checks {
		mark := "FAIL"
		if c.OK {
			mark = "OK"
		}
		fmt.Printf("  [%s] %s — %s\n", mark, c.Label, c.Detail)
	}

	fmt.Println("=== UPDATE CHECK ===")
	u := update.Check()
	fmt.Printf("ok=%v installed=%v msg=%s\n", u.OK, u.Installed, u.Message)

	fmt.Println("=== INSTALL ===")
	res, err := install.Run(install.Options{
		InstallDir:      dest,
		DesktopShortcut: false,
		IsUpdate:        false,
		ReleaseTag:      "v" + paths.AppVersion,
		Payload:         raw,
	}, func(step, detail string, percent int) {
		fmt.Printf("  [%d%%] %s — %s\n", percent, step, detail)
	})
	if err != nil {
		fmt.Println("INSTALL ERROR:", err)
		os.Exit(1)
	}
	fmt.Println(res.Summary)

	m, err := manifest.Read(dest)
	if err != nil {
		panic(err)
	}
	fmt.Printf("manifest version=%s exe=%s smoke=%v\n", m.AppVersion, m.ExePath, m.SmokeOK)
	if _, err := os.Stat(filepath.Join(dest, "Slate.exe")); err != nil {
		fmt.Println("missing exe")
		os.Exit(1)
	}
	if !res.SmokeOK {
		fmt.Println("SMOKE FAIL:", res.SmokeDetail)
		os.Exit(1)
	}
	fmt.Println("=== QA PASS ===")
}
