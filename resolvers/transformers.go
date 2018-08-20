package resolvers

import (
	models "../models"
	"gopkg.in/mgo.v2/bson"
)

// TransformAccount transforms interface to Account model
func TransformAccount(in interface{}) models.Account {
	var account models.Account
	switch v := in.(type) {
	case bson.M:
		account.ID = v["_id"].(bson.ObjectId)
		account.Email = v["email"].(string)
		account.Name = v["name"].(string)
		account.Surname = v["surname"].(string)
		account.Password = v["password"].(string)
		account.AccessLevel = v["access_level"].(int)
		account.HunterID = v["hunter_id"].(bson.ObjectId)
		account.RecruitID = v["recruit_id"].(bson.ObjectId)
	case models.Account:
		account = v
	}

	return account
}

// TransformTokenManager transforms interface into TokenManager model
func TransformTokenManager(in interface{}) models.TokenManager {
	var tokenMgr models.TokenManager
	switch v := in.(type) {
	case bson.M:
		tokenMgr.ID = v["_id"].(bson.ObjectId)
		tokenMgr.RefreshToken = v["refresh_token"].(string)
		tokenMgr.MaxTokens = v["max_tokens"].(int)
		tokenMgr.AccountID = v["account_id"].(bson.ObjectId)
		tokenMgr.Tokens = v["tokens"].([]string)

	case models.TokenManager:
		tokenMgr = v
	}

	return tokenMgr
}

// TransformQA transforms interface into QA model
func TransformQA(in interface{}) models.QA {
	var qa models.QA
	switch v := in.(type) {
	case bson.M:
		qa.Question = v["question"].(string)
		qa.Answer = v["answer"].(string)
	case models.QA:
		qa = v
	}
	return qa
}

// TransformRecruit transforms interface into Recruit model
func TransformRecruit(in interface{}) models.Recruit {
	var recruit models.Recruit
	switch v := in.(type) {
	case bson.M:
		recruit.ID = v["_id"].(bson.ObjectId)
		recruit.BirthYear = v["birth_year"].(int32)
		recruit.Province = v["province"].(string)
		recruit.City = v["city"].(string)
		recruit.Gender = v["gender"].(string)
		recruit.Disability = v["disability"].(string)
		recruit.Vid1Url = v["vid1_url"].(string)
		recruit.Vid2Url = v["vid2_url"].(string)
		recruit.Phone = v["phone"].(string)
		recruit.Email = v["email"].(string)
		recruit.Qa1 = TransformQA(v["qa1"])
		recruit.Qa2 = TransformQA(v["qa2"])

	case models.Recruit:
		recruit = v
	}
	return recruit
}

// TransformIndustry transforms interface into Industry model
func TransformIndustry(in interface{}) models.Industry {
	var industry models.Industry
	switch v := in.(type) {
	case bson.M:
		industry.ID = v["_id"].(bson.ObjectId)
		industry.Name = v["name"].(string)

	case models.Industry:
		industry = v
	}

	return industry
}

// TransformQuestion transforms interface into Question model
func TransformQuestion(in interface{}) models.Question {
	var question models.Question
	switch v := in.(type) {
	case bson.M:
		question.ID = v["_id"].(bson.ObjectId)
		question.IndustryID = v["industry_id"].(bson.ObjectId)
		question.Question = v["question"].(string)

	case models.Question:
		question = v
	}

	return question
}

// TransformDocument transforms interface into Document model
func TransformDocument(in interface{}) models.Document {
	var document models.Document
	switch v := in.(type) {
	case bson.M:
		document.ID = v["_id"].(bson.ObjectId)
		document.OwnerID = v["owner_id"].(bson.ObjectId)
		document.URL = v["url"].(string)
		document.DocType = v["doc_type"].(string)
		document.OwnerType = v["owner_type"].(string)

	case models.Document:
		document = v
	}

	return document
}
