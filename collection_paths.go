package filters

import (
	"strings"

	vocab "github.com/go-ap/activitypub"
)

const (
	// ActorsType is a constant that represents the URL path for the local actors collection.
	// It is used as the parent for all To IDs
	ActorsType = vocab.CollectionPath("actors")
	// ActivitiesType is a constant that represents the URL path for the local activities collection
	// It is used as the parent for all Activity IDs
	ActivitiesType = vocab.CollectionPath("activities")
	// ObjectsType is a constant that represents the URL path for the local objects collection
	// It is used as the parent for all non To, non Activity Object IDs
	ObjectsType = vocab.CollectionPath("objects")

	// BlockedType is an internally used collection, to store a list of actors the actor has blocked
	BlockedType = vocab.CollectionPath("blocked")

	// IgnoredType is an internally used collection, to store a list of actors the actor has ignored
	IgnoredType = vocab.CollectionPath("ignored")
)

// TODO(marius): here we need a better separation between the collections which are exposed in the HTTP API
//
//	(activities,actors,objects) and the ones that are internal (blocked,ignored)
var (
	HiddenCollections = vocab.CollectionPaths{
		BlockedType,
		IgnoredType,
	}

	FedBOXCollections = vocab.CollectionPaths{
		ActivitiesType,
		ActorsType,
		ObjectsType,
		BlockedType,
		IgnoredType,
	}

	validActivityCollection = vocab.CollectionPaths{
		ActivitiesType,
	}

	validObjectCollection = vocab.CollectionPaths{
		ActorsType,
		ObjectsType,
	}
)

func getValidActivityCollection(typ vocab.CollectionPath) vocab.CollectionPath {
	for _, t := range validActivityCollection {
		if strings.EqualFold(string(typ), string(t)) {
			return t
		}
	}
	return vocab.Unknown
}

func getValidObjectCollection(typ vocab.CollectionPath) vocab.CollectionPath {
	for _, t := range validObjectCollection {
		if strings.EqualFold(string(typ), string(t)) {
			return t
		}
	}
	return vocab.Unknown
}

// ValidCollection shows if the current ActivityPub end-point type is a valid collection
func ValidCollection(typ vocab.CollectionPath) bool {
	return ValidActivityCollection(typ) || ValidObjectCollection(typ)
}

// ValidActivityCollection shows if the current ActivityPub end-point type is a valid collection for handling Activities
func ValidActivityCollection(typ vocab.CollectionPath) bool {
	return getValidActivityCollection(typ) != vocab.Unknown || vocab.ValidActivityCollection(typ)
}

// ValidObjectCollection shows if the current ActivityPub end-point type is a valid collection for handling Objects
func ValidObjectCollection(typ vocab.CollectionPath) bool {
	return getValidObjectCollection(typ) != vocab.Unknown || vocab.ValidObjectCollection(typ)
}
