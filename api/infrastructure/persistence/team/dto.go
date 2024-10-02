package team

type Dto struct {
	EntityId string `dynamodbav:"EntityId"`
	Id       string `dynamodbav:"Id"`
	Category int    `dynamodbav:"Category"`
	Sport    string `dynamodbav:"Sport"`
}
