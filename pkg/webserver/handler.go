package webserver

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

type handler interface {
	HandleCreateKeygroup(k keygroup.Keygroup) error
	HandleDeleteKeygroup(k keygroup.Keygroup) error
	HandleRead(i data.Item) (data.Item, error)
	HandleUpdate(i data.Item) error
	HandleDelete(i data.Item) error
}