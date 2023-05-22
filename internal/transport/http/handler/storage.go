package handler

import (
	"fmt"
	"github.com/ivahaev/go-logger"
	"github.com/kataras/iris/v12"
)

func (h *Handler) GetFileByKey(ctx iris.Context) {
	key := ctx.Params().GetString("key")

	file := h.Service.Storage.GetFileByKey(key)
	if file == nil {
		ctx.StatusCode(405)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"code":    405,
				"message": "File is not exists",
			},
			"data": []iris.Map{},
		})
		return
	}

	ctx.ContentType(fmt.Sprintf(`%s/%s`, file.Type, file.Ext))
	err := ctx.ServeFile(fmt.Sprintf(`./%s/%s.%s`, file.Path, file.Key, file.Ext))
	if err != nil {
		logger.Error(err)
		ctx.StatusCode(405)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"code":    500,
				"message": "Server uploaded a file with an error",
			},
			"data": []iris.Map{},
		})
		return
	}
}

func (h *Handler) TestVast(ctx iris.Context) {
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"code":    500,
			"message": "Server uploaded a file with an error",
		},
		"data": `<VAST version="3.0">
    <Ad id="123" type="front">
        <InLine>
            <AdSystem><![CDATA[DSP]]></AdSystem>
            <AdTitle><![CDATA[adTitle]]></AdTitle>
            <Impression id="11111"><![CDATA[http://impressionv1.track.com]]></Impression>
            <Impression id="11112"><![CDATA[http://impressionv2.track.com]]></Impression>
            <Creatives>
                <Creative id="987">
                    <Linear>
                        <Duration>00:00:15</Duration>
                        <TrackingEvents>
                            <Tracking event="start"><![CDATA[http://track.xxx.com/q/start?xx]]></Tracking>
                            <Tracking event="firstQuartile"><![CDATA[http://track.xxx.com/q/firstQuartile?xx]]></Tracking>
                            <Tracking event="midpoint"><![CDATA[http://track.xxx.com/q/midpoint?xx]]></Tracking>
                            <Tracking event="thirdQuartile"><![CDATA[http://track.xxx.com/q/thirdQuartile?xx]]></Tracking>
                            <Tracking event="complete"><![CDATA[http://track.xxx.com/q/complete?xx]]></Tracking>
                        </TrackingEvents>
                        <MediaFiles>
                            <MediaFile delivery="progressive" type="video/mp4" width="1024" height="576"><![CDATA[https://samplelib.com/lib/preview/mp4/sample-5s.mp4]]></MediaFile>
                        </MediaFiles>
                    </Linear>
                </Creative>
            </Creatives>
            <Description></Description>
            <Survey></Survey>
        </InLine>
    </Ad>
</VAST>`,
	})
	return
}
