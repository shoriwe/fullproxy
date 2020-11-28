package Authentication

import (
	"encoding/base64"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"log"
	"os/exec"
)

/*
Pass base64 encoded username and password to a command, ExitCode = 0 means "Login success", any other way will mean it not
*/
func CommandAuth(executable string, defaultArgv []string) Types.AuthenticationMethod {
	return func(username []byte, password []byte) bool {
		command := exec.Command(executable,
			append(defaultArgv, []string{
				base64.StdEncoding.EncodeToString(username),
				base64.StdEncoding.EncodeToString(password),
			}...)...)
		startError := command.Start()
		if startError != nil {
			log.Print(startError)
			return false
		}
		waitError := command.Wait()
		if waitError != nil {
			if _, ok := waitError.(*exec.ExitError); !ok {
				log.Print(waitError)
			}
			return false
		}
		return true
	}
}
