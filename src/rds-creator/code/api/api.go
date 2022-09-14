package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"serverless-terraform-mysql-creator/code/internal/terraform"
)

type params struct {
	OssConfig  terraform.OssConfig `json:"oss_config,omitempty"`
	InvokeType int                 `json:"invoke_type,omitempty"`
}

func Invoke(c *gin.Context, logger *logrus.Entry, stop chan int) {

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	p := &params{}

	if err := json.Unmarshal(jsonData, p); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
	}

	client := terraform.NewTerraformClient(logger, p.InvokeType, stop)

	client.GetOSSAndSecret(c, &p.OssConfig)

	errMessage := client.Validate()
	if errMessage != "" {
		RespondError(c, http.StatusBadRequest, errMessage)
		return
	}
	outputs, err := client.Do()
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondOk(c, outputs)

}
