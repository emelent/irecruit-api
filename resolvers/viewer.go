package resolvers

import (
	"log"

	config "../config"
	db "../database"
	er "../errors"
	models "../models"
	utils "../utils"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

type viewer interface {
	ID() graphql.ID
	Name() string
	Surname() string
	Email() string
}

type rViewerResolver struct {
	r    *models.Recruit
	a    *models.Account
	crud *db.CRUD
}

func (r *rViewerResolver) ID() graphql.ID {
	return graphql.ID(r.r.ID.Hex())
}

func (r *rViewerResolver) Name() string {
	return r.a.Name
}

func (r *rViewerResolver) Surname() string {
	return r.a.Surname
}

func (r *rViewerResolver) Email() string {
	return r.a.Email
}

func (r *rViewerResolver) Profile() (*recruitResolver, error) {
	defer r.crud.CloseCopy()

	// retrieve account
	rawAccount, err := r.crud.FindID(config.AccountsCollection, r.r.ID)
	if err != nil {
		log.Println(err)
		return nil, er.NewGenericError()
	}
	account := transformAccount(rawAccount)

	return &recruitResolver{r.r, &account}, nil
}

// type hViewerResolver struct {
// 	r *models.Hunter
// 	crud *db.CRUD
// }

// func (r *hViewerResolver) ID() graphql.ID {
// 	return graphql.ID(r.r.ID.Hex())
// }

// func (r *hViewerResolver) Name() string {
// 	return graphql.ID(r.r.Name)
// }

// func (r *hViewerResolver) Surname() string {
// 	return graphql.ID(r.r.Surname)
// }

// func (r *hViewerResolver) Email() string {
// 	return graphql.ID(r.r.Email)
// }

// func (r *hViewerResolver) Recruit(args struct{ID graphql.ID}) (*recruitResolver, error){
// 	id := string(args.ID)
// 	if !bson.IsObjectIdHex(id){
// 		log.Println("Invalid id!")
// 		return nil, er.GenericError()
// 	}

// 	rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
// 	if err !=  nil{
// 		log.Println(err)
// 		return nil, er.GenericError()
// 	}
// 	recruit  := transformRecruit(rawRecruit)
// 	rawAccount, err := r.crud.FindOne(
// 		config.AccountsCollection,
// 		&bson.M{"recruit_id": recruit.ID},
// 	)
// 	if err !=  nil{
// 		log.Println(err)
// 		return nil, er.GenericError()
// 	}
// 	account :=  transformAccount(rawAccount)
// 	return &recruitResolver{&recruit, &account}, nil
// }

type sViewerResolver struct {
	a    *models.Account
	crud *db.CRUD
}

func (r *sViewerResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

func (r *sViewerResolver) Name() string {
	return r.a.Name
}

func (r *sViewerResolver) Surname() string {
	return r.a.Surname
}

func (r *sViewerResolver) Email() string {
	return r.a.Email
}

func (r *sViewerResolver) Accounts() ([]*accountResolver, error) {
	defer r.crud.CloseCopy()

	// fetch all accounts
	rawAccounts, err := r.crud.FindAll(accountsCollection, nil)
	results := make([]*accountResolver, 0)
	for _, r := range rawAccounts {
		account := transformAccount(r)
		results = append(results, &accountResolver{&account})
	}
	return results, err
}

type viewerResolver struct {
	viewer
}

func (r *viewerResolver) ToRecruitViewer() (*rViewerResolver, bool) {
	v, ok := r.viewer.(*rViewerResolver)
	return v, ok
}

func (r *viewerResolver) ToSysViewer() (*sViewerResolver, bool) {
	v, ok := r.viewer.(*sViewerResolver)
	return v, ok
}

// View resolves "view" gql query
func (r *RootResolver) View(args struct{ Token *string }) (*viewerResolver, error) {

	// Did we get a token?
	if args.Token == nil {
		// TODO return a Guest ViewerResolver
		return nil, er.NewInvalidTokenError()
	}
	// get token claims
	claims, err := utils.GetTokenClaims(*args.Token)
	if err != nil {
		return nil, err
	}

	// no refresh tokens allowed
	if claims.Refresh {
		return nil, er.NewInvalidTokenError()
	}

	// get account
	if !bson.IsObjectIdHex(claims.AccountID) {
		return nil, er.NewInvalidTokenError()
	}

	rawAccount, err := r.crud.FindID(config.AccountsCollection, bson.ObjectIdHex(claims.AccountID))
	if err != nil {
		log.Println("Failed to find account by ID from token =>", err)
		return nil, er.NewGenericError()
	}
	account := transformAccount(rawAccount)

	// check if account is sys account
	if utils.IsSysAccount(&account) {
		viewer := &sViewerResolver{&account, r.crud}
		return &viewerResolver{viewer}, nil
	}

	// check if account has recruit profile
	if utils.IsNullID(account.RecruitID) {
		log.Println("This account has no recruit profile.")
		return nil, er.NewInputError("Sorry, but you're not 'viewer' material.")
	}

	rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
	if err != nil {
		log.Println("Failed to find recruit =>", err)
		return nil, er.NewGenericError()
	}
	recruit := transformRecruit(rawRecruit)
	viewer := &rViewerResolver{&recruit, &account, r.crud}
	return &viewerResolver{viewer}, nil
}
