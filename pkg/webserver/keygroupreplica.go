package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

func getKeygroupReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		r, err := h.HandleGetKeygroupReplica(keygroup.Keygroup{
			Name: kgname,
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.JSON(http.StatusOK, r)
		return
	}
}

func postKeygroupReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		nodeid := context.Params.ByName("nodeid")

		err := h.HandleAddKeygroupReplica(keygroup.Keygroup{
			Name: kgname,
		}, replication.Node{
			ID: replication.ID(nodeid),
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusConflict, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}

func deleteKeygroupReplica(h exthandler.Handler) func(context *gin.Context) {
	return func(context *gin.Context) {
		kgname := context.Params.ByName("kgname")

		nodeid := context.Params.ByName("nodeid")

		err := h.HandleRemoveKeygroupReplica(keygroup.Keygroup{
			Name: kgname,
		}, replication.Node{
			ID: replication.ID(nodeid),
		})

		if err != nil {
			_ = context.AbortWithError(http.StatusNotFound, err)
			return
		}

		context.Status(http.StatusOK)
		return
	}
}