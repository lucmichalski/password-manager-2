package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nnachevv/passmag/crypt"
	"github.com/nnachevv/passmag/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

// the questions to ask
var editQs = []*survey.Question{
	{
		Name:   "name",
		Prompt: &survey.Input{Message: "Enter name for which you want to change password:"},
	},
	{
		Name:   "newname",
		Prompt: &survey.Input{Message: "Enter new name for your password:"},
	},
}

var edit = &cobra.Command{
	Use:   "edit",
	Short: "Initialize email, password and master password for your password manager",
	Long:  `Set master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := storage.FilePath()
		if err != nil {
			return err
		}

		if err := storage.VaultExist(path); err != nil {
			return err
		}

		var sessionKey string
		if !viper.IsSet("PASS_SESSION") {
			prompt := &survey.Input{Message: "Please enter your session key :"}
			survey.AskOne(prompt, &sessionKey, survey.WithValidator(survey.Required))
		} else {
			sessionKey = viper.GetString("PASS_SESSION")
		}

		var masterPassword string
		prompt := &survey.Password{Message: "Enter your  master password:"}
		survey.AskOne(prompt, &masterPassword, survey.WithValidator(survey.Required))

		vaultPwd := argon2.IDKey([]byte(masterPassword), []byte(sessionKey), 1, 64*1024, 4, 32)

		vaultData, err := crypt.DecryptFile(path, vaultPwd)

		if err != nil {
			return err
		}

		var s storage.Storage

		err = json.Unmarshal(vaultData, &s)
		if err != nil {
			return err
		}

		answers := struct {
			Name    string
			NewName string
		}{}

		pwd, err := s.Get(answers.Name)
		if err != nil {
			return err
		}

		err = s.Remove(answers.Name)
		if err != nil {
			return err
		}

		err = s.Add(answers.NewName, pwd)
		if err != nil {
			return err
		}
		fmt.Println("succesfuly moved your password")

		return nil
	},
}
