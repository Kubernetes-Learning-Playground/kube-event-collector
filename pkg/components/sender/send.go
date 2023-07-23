package sender

import (
	"crypto/tls"
	"fmt"
	"github.com/practice/kube-event/pkg/config"
	"github.com/practice/kube-event/pkg/model"
	"gopkg.in/mail.v2"
	"k8s.io/klog/v2"
)

//var (
// GlobalSend *Sender
//)
//
//func init() {
// GlobalSend = NewSender() // FIXME: 目前有sender实例无法重复利用的bug
//}

// Sender 邮件发送器
type Sender struct {
	dialer *mail.Dialer
	cfg    *config.Sender
}

// Send 发送邮件
func (sender *Sender) Send(event *model.Event) error {
	// TODO: 需要抛出错误
	m := mail.NewMessage()
	m.SetHeader("From", sender.cfg.Email)
	m.SetHeader("To", sender.cfg.Targets)
	title := fmt.Sprintf("new event generated, event name: %s, event type: %s", event.Name, event.Type)
	m.SetHeader("Subject", title)
	content := fmt.Sprintf("event message: %s, event reason: %s", event.Message, event.Reason)
	m.SetBody("text/plain", content)
	if err := sender.dialer.DialAndSend(m); err != nil {
		klog.Error("send err: ", err)
		return err
	}
	klog.Info("send email.....")
	return nil
}

// NewSender 创建邮件发送器
func NewSender(config *config.Config) *Sender {
	d := mail.NewDialer(config.Sender.Remote, config.Sender.Port,
		config.Sender.Email, config.Sender.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &Sender{
		dialer: d,
		cfg:    &config.Sender,
	}
}
