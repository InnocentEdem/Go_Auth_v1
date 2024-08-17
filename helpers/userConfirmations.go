package helpers

import (
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/utils"
)
func ConfirmUser( client models.Client, user models.User) error {
    var confirmationMethods  = client.Confirmation

    if confirmationMethods.ConfirmEmail {
        utils.SendConfirmationEmail(client,user)
    }

    return nil
}
