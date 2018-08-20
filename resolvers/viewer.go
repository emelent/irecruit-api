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

// -----------------
// Root Resolver Methods struct
// -----------------

// View resolves "view" gql query
func (r *RootResolver) View(args struct {
	Token   *string
	Enforce *string
}) (*viewerResolver, error) {

	// Did we get a token?
	if args.Token == nil {
		// return a guest viewer
		return &viewerResolver{&guestViewerResolver{}}, nil
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
		return nil, er.NewInvalidTokenError()
	}
	account := transformAccount(rawAccount)

	// enforce an enforceable
	if args.Enforce != nil {
		switch *args.Enforce {
		case "RECRUIT": // try to enforce recruit

			// check if account has recruit profile
			if utils.IsNullID(account.RecruitID) {
				return nil, er.NewInputError("Failed to enfore 'RECRUIT'.")
			}

			// retrieve Recruit profile
			rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
			if err != nil {
				log.Println("Failed to find recruit =>", err)
				return nil, er.NewGenericError()
			}

			// return RecruitViewer
			recruit := transformRecruit(rawRecruit)
			viewer := &recruitViewerResolver{&recruit, &account, r.crud}
			return &viewerResolver{viewer}, nil

		case "HUNTER": // try to enforce hunter
			return nil, er.NewInputError("Unimplemented")

		case "SYSTEM": // try to enforce system
			// check if account is sys account
			if !utils.IsSysAccount(&account) {
				return nil, er.NewInputError("Failed to enfore 'SYSTEM'.")
			}

			// return sysViewer
			viewer := &sysViewerResolver{&account, r.crud}
			return &viewerResolver{viewer}, nil

		case "ACCOUNT":
			// return accountViewer
			return &viewerResolver{&accountViewerResolver{&account}}, nil

		}
	}

	// check if account is sys account
	if utils.IsSysAccount(&account) {
		viewer := &sysViewerResolver{&account, r.crud}
		return &viewerResolver{viewer}, nil
	}

	// check if account has recruit profile
	if !utils.IsNullID(account.RecruitID) {
		rawRecruit, err := r.crud.FindID(config.RecruitsCollection, account.RecruitID)
		if err != nil {
			log.Println("Failed to find recruit =>", err)
			return nil, er.NewGenericError()
		}
		recruit := transformRecruit(rawRecruit)
		viewer := &recruitViewerResolver{&recruit, &account, r.crud}
		return &viewerResolver{viewer}, nil
	}

	return nil, er.NewGenericError()
}

// -----------------
// viewer interface
// -----------------
type viewer interface {
	ID() graphql.ID
	Name() string
	Surname() string
	Email() string
}

// -----------------
// recruitViewerResolver struct
// -----------------
type recruitViewerResolver struct {
	r    *models.Recruit
	a    *models.Account
	crud *db.CRUD
}

func (r *recruitViewerResolver) ID() graphql.ID {
	return graphql.ID(r.r.ID.Hex())
}

func (r *recruitViewerResolver) Name() string {
	return r.a.Name
}

func (r *recruitViewerResolver) Surname() string {
	return r.a.Surname
}

func (r *recruitViewerResolver) Email() string {
	return r.a.Email
}

func (r *recruitViewerResolver) Profile() (*recruitResolver, error) {
	defer r.crud.CloseCopy()

	// retrieve account
	rawAccount, err := r.crud.FindID(config.AccountsCollection, r.a.ID)
	if err != nil {
		log.Println(err)
		return nil, er.NewGenericError()
	}
	account := transformAccount(rawAccount)

	return &recruitResolver{r.r, &account}, nil
}

// -----------------
// hunterViewerResolver struct
// -----------------
// type hunterViewerResolver struct {
// 	r *models.Hunter
// 	crud *db.CRUD
// }

// func (r *hunterViewerResolver) ID() graphql.ID {
// 	return graphql.ID(r.r.ID.Hex())
// }

// func (r *hunterViewerResolver) Name() string {
// 	return graphql.ID(r.r.Name)
// }

// func (r *hunterViewerResolver) Surname() string {
// 	return graphql.ID(r.r.Surname)
// }

// func (r *hunterViewerResolver) Email() string {
// 	return graphql.ID(r.r.Email)
// }

// func (r *hunterViewerResolver) Recruit(args struct{ID graphql.ID}) (*recruitResolver, error){
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

// -----------------
// sysViewerResolver struct
// -----------------
type sysViewerResolver struct {
	a    *models.Account
	crud *db.CRUD
}

func (r *sysViewerResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

func (r *sysViewerResolver) Name() string {
	return r.a.Name
}

func (r *sysViewerResolver) Surname() string {
	return r.a.Surname
}

func (r *sysViewerResolver) Email() string {
	return r.a.Email
}

func (r *sysViewerResolver) Accounts() ([]*accountResolver, error) {
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

// -----------------
// guestViewerResolver struct
// -----------------
type guestViewerResolver struct {
	crud *db.CRUD
}

func (r *guestViewerResolver) ID() graphql.ID {
	return graphql.ID("GUEST")
}

func (r *guestViewerResolver) Name() string {
	return "Guest"
}

func (r *guestViewerResolver) Surname() string {
	return ""
}

func (r *guestViewerResolver) Email() string {
	return ""
}

// -----------------
// accountViewerResolver struct
// -----------------
type accountViewerResolver struct {
	a *models.Account
}

func (r *accountViewerResolver) ID() graphql.ID {
	return graphql.ID(r.a.ID.Hex())
}

func (r *accountViewerResolver) Name() string {
	return r.a.Name
}

func (r *accountViewerResolver) Surname() string {
	return r.a.Surname
}

func (r *accountViewerResolver) Email() string {
	return r.a.Email
}

func (r *accountViewerResolver) IsHunter() bool {
	return utils.IsNullID(r.a.HunterID)
}

func (r *accountViewerResolver) IsRecruit() bool {
	return utils.IsNullID(r.a.RecruitID)
}

func (r *accountViewerResolver) CheckPassword(args struct{ Password string }) bool {
	return r.a.CheckPassword(args.Password)
}

// -----------------
// viewerResolver struct
// -----------------
type viewerResolver struct {
	viewer
}

func (r *viewerResolver) ToRecruitViewer() (*recruitViewerResolver, bool) {
	v, ok := r.viewer.(*recruitViewerResolver)
	return v, ok
}

func (r *viewerResolver) ToSysViewer() (*sysViewerResolver, bool) {
	v, ok := r.viewer.(*sysViewerResolver)
	return v, ok
}

func (r *viewerResolver) ToGuestViewer() (*guestViewerResolver, bool) {
	v, ok := r.viewer.(*guestViewerResolver)
	return v, ok
}

func (r *viewerResolver) ToAccountViewer() (*accountViewerResolver, bool) {
	v, ok := r.viewer.(*accountViewerResolver)
	return v, ok
}
