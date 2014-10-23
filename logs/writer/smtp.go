// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"bytes"
	"io"
	"net/smtp"
	"strings"
)

// 实现io.Writer接口的邮件发送。
type Smtp struct {
	username string   // smtp账号
	password string   // smtp密码
	host     string   // smtp主机，需要带上端口
	sendTo   []string // 接收者。
	subject  string   // 邮件主题。

	// 邮件内容的缓存
	cache *bytes.Buffer
	// 邮件头部分的长度
	headerLen int

	// BUG(caixw): 缓存smtp.Auth，一般情况下是没有问题，若
	// smtp.Auth的实现者在实例里保存状态值之类的东西，则不
	// 能缓存只能在每次SendMail的时候实时申请，会造成大量的
	// 内存碎片，可考虑sync.Pool或是直接重写smtp.SendMail()
	// 函数来提升性能。
	auth smtp.Auth
}

var _ io.Writer = &Smtp{}

// 新建Smtp对象。
// username为smtp的账号；
// password为smtp对应的密码；
// subject为发送邮件的主题；
// host为smtp的主机地址，需要带上端口号；
// sendTo为接收者的地址。
func NewSmtp(username, password, subject, host string, sendTo []string) *Smtp {
	ret := &Smtp{
		username: username,
		password: password,
		subject:  subject,
		host:     host,
		sendTo:   sendTo,
	}
	ret.init()

	return ret
}

// 初始化一些基本内容。
//
// 像To,From这些内容都是固定的，可以先写入到缓存中，这样
// 这后就不需要再次构造这些内容。
func (s *Smtp) init() {
	s.cache = bytes.NewBufferString("")
	s.cache.Grow(1024)

	// to
	s.cache.WriteString("To: ")
	s.cache.WriteString(strings.Join(s.sendTo, ";"))
	s.cache.WriteString("\r\n")

	// from
	s.cache.WriteString("From: ")
	s.cache.WriteString(s.username) // <...>有需要吗？
	s.cache.WriteString("\r\n")

	// subject
	s.cache.WriteString("Subject: ")
	s.cache.WriteString(s.subject)
	s.cache.WriteString("\r\n")

	// mime-version
	s.cache.WriteString("MIME-Version: ")
	s.cache.WriteString("1.0\r\n")

	// contentType
	s.cache.WriteString(`Content-Type: text/plain; charset="utf-8"`)
	s.cache.WriteString("\r\n\r\n")

	s.headerLen = s.cache.Len()

	// 去掉端口部分
	h := strings.Split(s.host, ":")[0]
	s.auth = smtp.PlainAuth("", s.username, s.password, h)
}

// io.Writer
func (s *Smtp) Write(msg []byte) (int, error) {
	s.cache.Write(msg)

	err := smtp.SendMail(
		s.host,
		s.auth,
		s.username,
		s.sendTo,
		s.cache.Bytes(),
	)
	l := s.cache.Len()

	s.cache.Truncate(s.headerLen)

	return l, err
}
