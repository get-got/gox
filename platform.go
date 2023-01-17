package main

import (
	"fmt"
	"log"
	"strings"

	version "github.com/hashicorp/go-version"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string

	// Default, if true, will be included as a default build target
	// if no OS/arch is specified. We try to only set as a default popular
	// targets or targets that are generally useful. For example, Android
	// is not a default because it is quite rare that you're cross-compiling
	// something to Android AND something like Linux.
	Default bool
	ARM     string
}

func PlatformFromString(os, arch string) Platform {
	if strings.HasPrefix(arch, "armv") && len(arch) >= 5 {
		return Platform{
			OS:   os,
			Arch: "arm",
			ARM:  arch[4:],
		}
	}
	return Platform{
		OS:   os,
		Arch: arch,
	}
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.GetArch())
}

func (p *Platform) GetArch() string {
	return fmt.Sprintf("%s%s", p.Arch, p.GetARMVersion())
}

func (p *Platform) GetARMVersion() string {
	if len(p.ARM) > 0 {
		return "v" + p.ARM
	}
	return ""
}

// addDrop appends all of the "add" entries and drops the "drop" entries, ignoring
// the "Default" parameter.
func addDrop(base []Platform, add []Platform, drop []Platform) []Platform {
	newPlatforms := make([]Platform, len(base)+len(add))
	copy(newPlatforms, base)
	copy(newPlatforms[len(base):], add)

	// slow, but we only do this during initialization at most once per version
	for _, platform := range drop {
		found := -1
		for i := range newPlatforms {
			if newPlatforms[i].Arch == platform.Arch && newPlatforms[i].OS == platform.OS {
				found = i
				break
			}
		}
		if found < 0 {
			panic(fmt.Sprintf("Expected to remove %+v but not found in list %+v", platform, newPlatforms))
		}
		if found == len(newPlatforms)-1 {
			newPlatforms = newPlatforms[:found]
		} else if found == 0 {
			newPlatforms = newPlatforms[found:]
		} else {
			newPlatforms = append(newPlatforms[:found], newPlatforms[found+1:]...)
		}
	}
	return newPlatforms
}

// addDrop appends all of the "add" entries and drops the "drop" entries, ignoring
// the "Default" parameter.
func addDrop(base []Platform, add []Platform, drop []Platform) []Platform {
	newPlatforms := make([]Platform, len(base)+len(add))
	copy(newPlatforms, base)
	copy(newPlatforms[len(base):], add)

	// slow, but we only do this during initialization at most once per version
	for _, platform := range drop {
		found := -1
		for i := range newPlatforms {
			if newPlatforms[i].Arch == platform.Arch && newPlatforms[i].OS == platform.OS {
				found = i
				break
			}
		}
		if found < 0 {
			panic(fmt.Sprintf("Expected to remove %+v but not found in list %+v", platform, newPlatforms))
		}
		if found == len(newPlatforms)-1 {
			newPlatforms = newPlatforms[:found]
		} else if found == 0 {
			newPlatforms = newPlatforms[found:]
		} else {
			newPlatforms = append(newPlatforms[:found], newPlatforms[found+1:]...)
		}
	}
	return newPlatforms
}

var (
	OsList = []string{
		"darwin",
		"dragonfly",
		"linux",
		"android",
		"solaris",
		"freebsd",
		"nacl",
		"netbsd",
		"openbsd",
		"plan9",
		"windows",
	}

	ArchList = []string{
		"386",
		"amd64",
		"amd64p32",
		"arm",
		"arm64",
		"mips64",
		"mips64le",
		"ppc64",
		"ppc64le",
	}

	Platforms_1_0 = []Platform{
		{OS: "darwin", Arch: "386", Default: true},
		{OS: "darwin", Arch: "amd64", Default: true},
		{OS: "linux", Arch: "386", Default: true},
		{OS: "linux", Arch: "amd64", Default: true},
		{OS: "linux", Arch: "arm", Default: true, ARM: "5"},
		{OS: "linux", Arch: "arm", Default: true, ARM: "6"},
		{OS: "linux", Arch: "arm", Default: true, ARM: "7"},
		{OS: "linux", Arch: "arm", Default: true, ARM: "8"},
		{OS: "freebsd", Arch: "386", Default: true},
		{OS: "freebsd", Arch: "amd64", Default: true},
		{OS: "openbsd", Arch: "386", Default: true},
		{OS: "openbsd", Arch: "amd64", Default: true},
		{OS: "windows", Arch: "386", Default: true},
		{OS: "windows", Arch: "amd64", Default: true},
	}

	Platforms_1_1 = addDrop(Platforms_1_0, []Platform{
		{OS: "freebsd", Arch: "arm", Default: true},
		{OS: "netbsd", Arch: "386", Default: true},
		{OS: "netbsd", Arch: "amd64", Default: true},
		{OS: "netbsd", Arch: "arm", Default: true},
		{OS: "plan9", Arch: "386", Default: false},
	}, nil)

	Platforms_1_3 = addDrop(Platforms_1_1, []Platform{
		{OS: "dragonfly", Arch: "386", Default: false},
		{OS: "dragonfly", Arch: "amd64", Default: false},
		{OS: "nacl", Arch: "amd64", Default: false},
		{OS: "nacl", Arch: "amd64p32", Default: false},
		{OS: "nacl", Arch: "arm", Default: false},
		{OS: "solaris", Arch: "amd64", Default: false},
	}, nil)

	Platforms_1_4 = addDrop(Platforms_1_3, []Platform{
		{OS: "android", Arch: "arm", Default: false},
		{OS: "plan9", Arch: "amd64", Default: false},
	}, nil)

	Platforms_1_5 = addDrop(Platforms_1_4, []Platform{
		{OS: "darwin", Arch: "arm", Default: false},
		{OS: "darwin", Arch: "arm64", Default: false},
		{OS: "linux", Arch: "arm64", Default: false},
		{OS: "linux", Arch: "ppc64", Default: false},
		{OS: "linux", Arch: "ppc64le", Default: false},
	}, nil)

	Platforms_1_6 = addDrop(Platforms_1_5, []Platform{
		{OS: "android", Arch: "386", Default: false},
		{OS: "android", Arch: "amd64", Default: false},
		{OS: "linux", Arch: "mips64", Default: false},
		{OS: "linux", Arch: "mips64le", Default: false},
		{OS: "nacl", Arch: "386", Default: false},
		{OS: "openbsd", Arch: "arm", Default: true},
	}, nil)

	Platforms_1_7 = addDrop(Platforms_1_5, []Platform{
		// While not fully supported s390x is generally useful
		{OS: "linux", Arch: "s390x", Default: true},
		{OS: "plan9", Arch: "arm", Default: false},
		// Add the 1.6 Platforms, but reflect full support for mips64 and mips64le
		{OS: "android", Arch: "386", Default: false},
		{OS: "android", Arch: "amd64", Default: false},
		{OS: "linux", Arch: "mips64", Default: true},
		{OS: "linux", Arch: "mips64le", Default: true},
		{OS: "nacl", Arch: "386", Default: false},
		{OS: "openbsd", Arch: "arm", Default: true},
	}, nil)

	Platforms_1_8 = addDrop(Platforms_1_7, []Platform{
		{OS: "linux", Arch: "mips", Default: true},
		{OS: "linux", Arch: "mipsle", Default: true},
	}, nil)

	// no new platforms in 1.9
	Platforms_1_9 = Platforms_1_8

	// unannounced, but dropped support for android/amd64
	Platforms_1_10 = addDrop(Platforms_1_9, nil, []Platform{{OS: "android", Arch: "amd64", Default: false}})

	Platforms_1_11 = addDrop(Platforms_1_10, []Platform{
		{OS: "js", Arch: "wasm", Default: true},
	}, nil)

	Platforms_1_12 = addDrop(Platforms_1_11, []Platform{
		{OS: "aix", Arch: "ppc64", Default: false},
		{OS: "windows", Arch: "arm", Default: true},
	}, nil)

	Platforms_1_13 = addDrop(Platforms_1_12, []Platform{
		{OS: "illumos", Arch: "amd64", Default: false},
		{OS: "netbsd", Arch: "arm64", Default: true},
		{OS: "openbsd", Arch: "arm64", Default: true},
	}, nil)

	Platforms_1_14 = addDrop(Platforms_1_13, []Platform{
		{OS: "freebsd", Arch: "arm64", Default: true},
		{OS: "linux", Arch: "riscv64", Default: true},
	}, []Platform{
		// drop nacl
		{OS: "nacl", Arch: "386", Default: false},
		{OS: "nacl", Arch: "amd64", Default: false},
		{OS: "nacl", Arch: "arm", Default: false},
	})

	Platforms_1_15 = addDrop(Platforms_1_14, []Platform{
		{OS: "android", Arch: "arm64", Default: false},
	}, []Platform{
		// drop i386 macos
		{OS: "darwin", Arch: "386", Default: false},
	})

	Platforms_1_16 = addDrop(Platforms_1_15, []Platform{
		{OS: "android", Arch: "amd64", Default: false},
		{OS: "darwin", Arch: "arm64", Default: true},
		{OS: "openbsd", Arch: "mips64", Default: false},
	}, nil)

	Platforms_1_17 = addDrop(Platforms_1_16, []Platform{
		{OS: "windows", Arch: "arm64", Default: true},
	}, nil)

	// no new platforms in 1.18
	Platforms_1_18 = Platforms_1_17

	PlatformsLatest = Platforms_1_18
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	// Use latest if we get an unexpected version string
	if !strings.HasPrefix(v, "go") {
		return PlatformsLatest
	}
	// go-version only cares about version numbers
	v = v[2:]

	current, err := version.NewVersion(v)
	if err != nil {
		log.Printf("Unable to parse current go version: %s\n%s", v, err.Error())

		// Default to latest
		return PlatformsLatest
	}

	var platforms = []struct {
		constraint string
		plat       []Platform
	}{
		{"<= 1.0", Platforms_1_0},
		{">= 1.1, < 1.3", Platforms_1_1},
		{">= 1.3, < 1.4", Platforms_1_3},
		{">= 1.4, < 1.5", Platforms_1_4},
		{">= 1.5, < 1.6", Platforms_1_5},
		{">= 1.6, < 1.7", Platforms_1_6},
		{">= 1.7, < 1.8", Platforms_1_7},
		{">= 1.8, < 1.9", Platforms_1_8},
		{">= 1.9, < 1.10", Platforms_1_9},
		{">= 1.10, < 1.11", Platforms_1_10},
		{">= 1.11, < 1.12", Platforms_1_11},
		{">= 1.12, < 1.13", Platforms_1_12},
		{">= 1.13, < 1.14", Platforms_1_13},
		{">= 1.14, < 1.15", Platforms_1_14},
		{">= 1.15, < 1.16", Platforms_1_15},
		{">= 1.16, < 1.17", Platforms_1_16},
		{">= 1.17, < 1.18", Platforms_1_17},
		{">= 1.18, < 1.19", Platforms_1_18},
	}

	for _, p := range platforms {
		constraints, err := version.NewConstraint(p.constraint)
		if err != nil {
			panic(err)
		}
		if constraints.Check(current) {
			return p.plat
		}
	}

	// Assume latest
	return PlatformsLatest
}
