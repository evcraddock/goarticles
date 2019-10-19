package repos

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"github.com/evcraddock/goarticles/internal/services"
	"github.com/evcraddock/goarticles/pkg/links"
)

type LinkRepository struct {
	Server       string
	DatabaseName string
}

func CreateLinkRepository(server, databaseName string) *LinkRepository {
	return &LinkRepository{
		Server:       server,
		DatabaseName: databaseName,
	}
}

func (r *LinkRepository) GetLinks(query map[string]interface{}) (*links.Links, error) {
	log.Debugf("Connecting to database %v", r.Server)
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("links")
	results := links.Links{}
	if err := services.NewError(
		c.Find(query).Sort("-createddate").All(&results),
		"error retrieving data",
		"DatabaseError",
		false); err != nil {
		return nil, err
	}

	return &results, nil
}

func (r *LinkRepository) AddLink(link links.Link) (*links.Link, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()

	link.ID = bson.NewObjectId()
	if err := services.NewError(
		session.DB(r.DatabaseName).C("links").Insert(link),
		"failed to create link",
		"DatabaseError",
		false); err != nil {
		return nil, err
	}

	log.Debug("Added Link ID: ", link.ID)

	return &link, nil
}

func (r *LinkRepository) DeleteLink(id string) error {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return err
	}

	defer session.Close()
	c := session.DB(r.DatabaseName).C("links")
	oid, err := r.linkExists(c, id)
	if err != nil {
		return err
	}

	if err = c.RemoveId(oid); err != nil {
		return services.NewError(err, "failed to delete link", "DatabaseError", false)
	}

	log.Debug("Delete Link ID: ", oid)

	return nil
}

func (r *LinkRepository) LinkExists(id string) (bool, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return false, err
	}

	defer session.Close()
	c := session.DB(r.DatabaseName).C("links")
	if _, err := r.linkExists(c, id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *LinkRepository) linkExists(collection *mgo.Collection, id string) (*bson.ObjectId, error) {
	if !bson.IsObjectIdHex(id) {
		err := services.NewError(fmt.Errorf("invalid id"), "can not find record: invalid id", "NotFound", false)
		return nil, err
	}

	oid := bson.ObjectIdHex(id)
	count, err := collection.FindId(oid).Count()
	if err != nil {
		return nil, services.NewError(err, "could not find link", "DatabaseError", false)
	}

	if count < 1 {
		return nil, services.NewError(fmt.Errorf("link does not exist"), "link does not exist", "NotFound", false)
	}

	return &oid, nil
}
