package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spo-iitk/ras-backend/mail"
)

func companySignUpHandler(mail_channel chan mail.Mail) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var signupReq CompanySignUpRequest

		err := ctx.ShouldBindJSON(&signupReq)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := createCompany(ctx, &signupReq)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logrus.Infof("A Company %s made signUp request with id %d", signupReq.CompanyName, id)
		mail_channel <- mail.GenerateMail(signupReq.Email,
			"Registration requested on PhD RAS",
			"Dear "+signupReq.Name+",\n\nWe got your request for registration on PhD Recruitment Automation System, IIT Kanpur. We will get back to you soon. For any queries, please get in touch with us at phdplacement@iitk.ac.in.")

		mail_channel <- mail.GenerateMail("phdplacement@iitk.ac.in",
			"Registration requested on RAS",
			"Company "+signupReq.CompanyName+" has requested to be registered on PhD RAS. The details are as follows:\n\n"+
				"Name: "+signupReq.Name+"\n"+
				"Designation: "+signupReq.Designation+"\n"+
				"Email: "+signupReq.Email+"\n"+
				"Phone: "+signupReq.Phone+"\n"+
				"Comments: "+signupReq.Comments+"\n")

		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully Requested"})
	}
}
