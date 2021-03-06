package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/nnachevv/passmag/cmd"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ = Describe("Add", func() {
	var (
		path   string
		addCmd *cobra.Command
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)

	BeforeEach(func() {
		addCmd = cmd.NewAddCmd()
		addCmd.SetArgs([]string{})
		addCmd.SetOut(&stdOut)
		addCmd.SetErr(&stdErr)

		fName, err := tempFile("fixtures/vault.bin")
		Expect(err).ShouldNot(HaveOccurred())
		path = filepath.Join(os.TempDir(), fName)
		viper.Set("path", path)
		viper.Set("PASS_SESSION", "fixed")
	})

	Context("When user set PASS_SESSION ,decline generation, pass right arguments for his passwords", func() {
		It("should add to vault his password", func() {
			c, state, err := vt10x.NewVT10XConsole()
			qf.SetStdio(terminal.Stdio{
				In:  c.Tty(),
				Out: c.Tty(),
				Err: c.Tty()})
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your  master password:")
				c.SendLine("master")
				c.ExpectString("Enter name for your password:")
				c.SendLine("dummy-name")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
			}()

			err := addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("want password to be generated, pass name for his password", func() {
		It("should add to vault his password", func() {
			c, state, err := vt10x.NewVT10XConsole()
			qf.SetStdio(terminal.Stdio{
				In:  c.Tty(),
				Out: c.Tty(),
				Err: c.Tty()})
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your  master password:")
				c.SendLine("master")
				c.ExpectString("Enter name for your password:")
				c.SendLine("dummy-name")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
			}()

			err := addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("trying to add already exist password and decline editing already existing password", func() {
		It("should add to vault his password", func() {
			c, state, err := vt10x.NewVT10XConsole()
			qf.SetStdio(terminal.Stdio{
				In:  c.Tty(),
				Out: c.Tty(),
				Err: c.Tty()})
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your  master password:")
				c.SendLine("master")
				c.ExpectString("Enter name for your password:")
				c.SendLine("exist")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
				c.ExpectString("This name with password already exist! Do you want to edit name with newly password")
				c.SendLine("N")
			}()

			err := addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})

	Context("trying to add already exist password and edit", func() {
		It("should add to vault his password", func() {
			c, state, err := vt10x.NewVT10XConsole()
			qf.SetStdio(terminal.Stdio{
				In:  c.Tty(),
				Out: c.Tty(),
				Err: c.Tty()})
			Expect(err).ShouldNot(HaveOccurred())
			defer c.Close()
			done := make(chan struct{})

			go func() {
				defer close(done)
				c.ExpectString("Enter your  master password:")
				c.SendLine("master")
				c.ExpectString("Enter name for your password:")
				c.SendLine("exist")
				c.ExpectString("Do you want to automatically generate password?")
				c.SendLine("N")
				c.ExpectString("Enter your password:")
				c.SendLine("dummy-password")
				c.ExpectString("This name with password already exist! Do you want to edit name with newly password")
				c.SendLine("y")
			}()

			err := addCmd.Execute()
			Expect(err).ShouldNot(HaveOccurred())

			c.Tty().Close()
			<-done
			fmt.Fprintf(ginkgo.GinkgoWriter, "--- Terminal ---\n%s\n----------------\n", expect.StripTrailingEmptyLines(state.String()))
		})
	})
})

// creates a new temporary file
func tempFile(path string) (string, error) {
	tar, err := os.Open(path)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(tar)
	Expect(err).ShouldNot(HaveOccurred())
	defer tar.Close()

	file, err := ioutil.TempFile(os.TempDir(), "fixture-file")
	Expect(err).ShouldNot(HaveOccurred())

	_, err = file.Write(bytes)
	Expect(err).ShouldNot(HaveOccurred())

	err = file.Sync()
	Expect(err).ShouldNot(HaveOccurred())

	_, err = file.Seek(0, io.SeekStart)
	Expect(err).ShouldNot(HaveOccurred())

	return file.Name(), nil
}
