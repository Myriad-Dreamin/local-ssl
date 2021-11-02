package ssl

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

type OpenSSLEnv struct {
	CmdName       string
	Out           *bytes.Buffer
	Err           *bytes.Buffer
	hasErr        bool
	userReportErr func(err error)
	wdStack       []string
}

func (env *OpenSSLEnv) reportErr(err error) {
	env.hasErr = true
	panic(err)
}

func (env *OpenSSLEnv) getErrBuf() error {
	if env.Err.Len() != 0 {
		return errors.New(env.Err.String())
	}
	return nil
}

func (env *OpenSSLEnv) validateExec() bool {
	if env.hasErr {
		return false
	}
	return true
}

func (env *OpenSSLEnv) output() string {
	return env.Out.String()
}

func (env *OpenSSLEnv) openSSLForced(params []string) *exec.Cmd {
	cmd := exec.Command(env.CmdName, params...)
	env.Out.Reset()
	env.Err.Reset()
	cmd.Stdout = env.Out
	cmd.Stderr = env.Err
	cmd.Dir = env.wdStack[len(env.wdStack)-1]
	err := cmd.Start()
	if err != nil {
		env.reportErr(err)
		return nil
	}
	err = cmd.Wait()

	if err != nil {
		if err = env.getErrBuf(); err != nil {
			env.reportErr(err)
			return nil
		}
		env.reportErr(err)
		return nil
	}
	code := cmd.ProcessState.ExitCode()
	if code != 0 {
		env.reportErr(fmt.Errorf("code: %v", code))
		return nil
	}
	return cmd
}

func (env *OpenSSLEnv) openSSL(params []string) *exec.Cmd {
	if env.hasErr {
		return nil
	}
	return env.openSSLForced(params)
}

func (env *OpenSSLEnv) HasErr() bool {
	return env.hasErr
}

type OpenSSLVersion struct {
	Version string
	Date    string
}

func (env *OpenSSLEnv) CheckVersion() *OpenSSLVersion {
	cmd := env.openSSL([]string{"version"})
	if cmd == nil {
		return nil
	}
	out := env.output()
	matched := regexp.MustCompile(`OpenSSL\s*([\w.]*)\s*(.*)`).FindStringSubmatch(out)
	if len(matched) == 0 {
		env.reportErr(fmt.Errorf("version string broken"))
	}
	return &OpenSSLVersion{
		Version: matched[1],
		Date:    matched[2],
	}
}

func (env *OpenSSLEnv) CheckBin() bool {
	v := env.CheckVersion()
	return v != nil
}

func (env *OpenSSLEnv) WriteConf(confLoc string, callback func(f *os.File) error) {
	if !env.validateExec() {
		return
	}
	f, err := os.Create(confLoc)
	if err != nil {
		env.reportErr(err)
		return
	}
	err = callback(f)
	if err != nil {
		env.reportErr(err)
		return
	}
	err = f.Close()
	if err != nil {
		env.reportErr(err)
		return
	}
}

func (env *OpenSSLEnv) WriteSSLConf(confLoc string, args *SSLTemplateArgs) {
	env.WriteConf(confLoc, func(f *os.File) error {
		return RenderSSLConf(f, args)
	})
}

func (env *OpenSSLEnv) WriteSignSSLConf(confLoc string, args *SignSSLTemplateArgs) {
	env.WriteConf(confLoc, func(f *os.File) error {
		return RenderSignSSLConf(f, args)
	})
}

func (env *OpenSSLEnv) GenerateRSAKey(loc string) {
	cmd := env.openSSL([]string{"genrsa", "-out", loc, "4096"})
	if cmd == nil {
		return
	}
	out := env.output()
	fmt.Print(out)
}

func (env *OpenSSLEnv) GenerateRootCACrt(confLoc string, priLoc string, crtLoc string) {
	cmd := env.openSSL([]string{"req", "-new", "-x509", "-days", "365", "-config", confLoc,
		"-key", priLoc, "-out", crtLoc})
	if cmd == nil {
		return
	}
	out := env.output()
	fmt.Print(out)
}

func (env *OpenSSLEnv) GenerateCSR(confLoc string, priLoc string, csrLoc string) {
	cmd := env.openSSL([]string{"req", "-new", "-config", confLoc, "-key", priLoc, "-out", csrLoc})
	if cmd == nil {
		return
	}
	out := env.output()
	fmt.Print(out)
}

func (env *OpenSSLEnv) CreateSignedCrt(confLoc string, csrLoc string, crtLoc string, caCrtLoc string, caPriLoc string) {
	cmd := env.openSSL([]string{"ca", "-batch", "-extensions", "v3_req",
		"-config", confLoc, "-in", csrLoc, "-out", crtLoc,
		"-cert", caCrtLoc, "-keyfile", caPriLoc})
	if cmd == nil {
		return
	}
	out := env.output()
	fmt.Print(out)
}

func (env *OpenSSLEnv) VerifyCrt(confLoc string, crtLoc string) {
	fmt.Println(confLoc, crtLoc)
	cmd := env.openSSL([]string{"req", "-noout", "-verify", "-config", confLoc, "-in", crtLoc})
	if cmd == nil {
		return
	}
	out := env.output()
	fmt.Print(out)
}

func (env *OpenSSLEnv) PushWd(wd string) {
	a, err := filepath.Abs(wd)
	if err != nil {
		env.reportErr(err)
		return
	}
	env.wdStack = append(env.wdStack, a)
	err = os.Chdir(a)
	if err != nil {
		env.reportErr(err)
		return
	}
}

func (env *OpenSSLEnv) MakeDir(path string) {
	if !env.validateExec() {
		return
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		env.reportErr(err)
		return
	}
}

func (env *OpenSSLEnv) SetFile(path string, content string) {
	f, err := os.Create(path)
	if err != nil {
		env.reportErr(err)
		return
	}
	_, err = f.WriteString(content)
	if err != nil {
		env.reportErr(err)
		return
	}
	err = f.Close()
	if err != nil {
		env.reportErr(err)
		return
	}
}
