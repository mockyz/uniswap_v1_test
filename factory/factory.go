package factory
import (
	sdk "github.com/ontio/ontology-go-sdk"
)

type UniswapFactory struct {
	ontSdk *sdk.OntologySdk
	templateCode string
	tokenHash string
}

