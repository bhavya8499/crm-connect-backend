package aws_pkg

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"sync"
)

var sess *session.Session
var once sync.Once

func GetAWSSession() *session.Session {
	if sess == nil {
		CreateAWSSession()
	}
	return sess
}

func CreateAWSSession() {
	once.Do(func() {
		sess = session.Must(session.NewSession())
	})
}
