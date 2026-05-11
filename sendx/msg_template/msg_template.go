package msg_template

const (
	VerificationCode = `
		<html>
		<body>
			<h2>Change password verification code</h2>
			<p>Your verification code is:<strong>%v</strong></p>
			<p>The validity period of the verification code is <strong>10 minute</strong>，Please use it as soon as possible.</p>
			<p>If this is not your own operation, please ignore this email.</p>
			<hr>
			<p style="color: gray;">This is a system email, please do not reply directly.</p>
		</body>
		</html>
	`
)
