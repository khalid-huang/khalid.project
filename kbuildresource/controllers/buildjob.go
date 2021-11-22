package controllers

import (
	"bryson.foundation/kbuildresource/async"
	"bryson.foundation/kbuildresource/common"
	"bryson.foundation/kbuildresource/dto"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/prometheus/common/log"
	"net/http"
)

type BuildJobController struct {
	beego.Controller
}

func (b *BuildJobController) CreateBuildJob() {
	var buildJobDTO dto.BuildJobDTO
	if err := json.Unmarshal(b.Ctx.Input.RequestBody, &buildJobDTO); err == nil {
		log.Infof("Create buildJob %s", buildJobDTO.Name)
		if buildJobDTO,err := async.GetRequestController().AcceptRequest(&buildJobDTO, common.BuildJobCreateRequestType); err == nil {
			b.Ctx.Output.SetStatus(http.StatusCreated)
			b.Data["json"] = common.GenerateResponse(common.ResponseSuccessResult, "create buildJob success",buildJobDTO)
		} else {
			b.Ctx.Output.SetStatus(http.StatusOK)
			b.Data["json"] = common.GenerateResponse(common.ResponseFailedResult, err.Error(), nil)
		}
	} else {
		b.Ctx.Output.SetStatus(http.StatusOK)
		b.Data["json"] = common.GenerateResponse(common.ResponseFailedResult, err.Error(), nil)
	}
	b.ServeJSON()
}