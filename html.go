package main

import "net/http"
import "fmt"

func css() string {
	return `
    .splash{width:300px;margin:200px auto;text-align:center;font-family:"Helvetica Neue", Helvetica, Arial}
    @media (max-width: 500px) {.splash{margin-top:100px}}
    .head{margin-bottom:40px}
    .logos{position:relative;margin-bottom:40px}
    .logo{width:100%}
    p{font-size:15px;margin:5px 0}
    select{background:none}
    input[type=submit], .form-item, div.message{
      font-size:12px;margin-top:10px;vertical-align:middle;
      display:block;text-align:center;box-sizing:
      border-box;width:100%;padding:9px}
    input[type=submit],div.message{
      font-weight:bold;border-width:0;
      text-transform:uppercase;cursor:pointer;appearence:none;
      -webkit-appearence:none;outline:0;}
      .error{color:#fff;background-color:#F4001E;text-transform:none}
    .neutral{color:#fff;background:#E01563}
    .success{color:#fff;background-color:#68C200}
    form{margin-top:20px;margin-bottom:0}
    input{color:#9B9B9B;border:1px solid #D6D6D6}
    input:focus{color:#666;border-color:#999;outline:0}
    .active{color:#E01563}

    p.signin{padding:10px 0 10px;font-size:11px}
    p.signin a{color:#E01563;text-decoration:none}
    p.signin a:hover{background-color:#E01563;color:#fff}
    footer{color:#D6D6D6;font-size:11px;margin:200px auto 0;width:300px;text-align:center}
    footer a{color:#9B9B9B;text-decoration:none;border-bottom:1px solid #9B9B9B}
    footer a:hover{color:#fff;background-color:#9B9B9B}
  `
}

func header(slack SlackInfo) string {
	users, active, _ := TotalUsers(slack.Name, slack.Token)

	return `
<html>
  <head>
    <title>Join ` + slack.DisplayName + ` on Slack!</title>
    <meta
      name="viewport"
      content="width=device-width,initial-scale=1.0,minimum-scale=1.0,user-scalable=no">
  </head>
  <div class="splash">
    <div class="logos">
      <img class="logo org" src="/banners/` + slack.Name + `.png" />
    </div>
    <p>Join <b>` + slack.DisplayName + `</b> on Slack.</p>
    <p class="status">
      <b class="active">` + fmt.Sprintf("%v", active) + `</b>
      users online now of
      <b class="total">` + fmt.Sprintf("%v", users) + `</b>
      registered.
    </p>`
}

func footer(slack SlackInfo) string {
	return `
    <p class="signin">or <a href="https://` + slack.Name + `.slack.com" target="_top">sign in</a>.</p>
    <style>` + css() + `</style>
  </div>
</html>
  `
}

func RenderBody(slack SlackInfo, w http.ResponseWriter, body string) {
	html := header(slack) + body + footer(slack)
	w.Write([]byte(html))
}

func RenderSuccess(slack SlackInfo, w http.ResponseWriter) {
	RenderBody(slack, w,
		`<div class='message success'>Sent! Check your inbox</div>`)
}

func RenderForm(slack SlackInfo, w http.ResponseWriter) {
	RenderBody(slack, w,
		`<form action='/request-invite' method='POST'>
    <input name="email" type="email" placeholder="you@yourdomain.com" autofocus=true class="form-item" />
    <input type="submit" class="neutral" value="Get my Invite" />
  </form>`)
}

func RenderError(slack SlackInfo, w http.ResponseWriter, msg string) {
	RenderBody(slack, w,
		`<div class='message error'>`+msg+`</div>`)
}
