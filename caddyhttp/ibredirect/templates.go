package ibredirect

const authPageTemplate = `
<html xmlns="http://www.w3.org/1999/xhtml">    
<head><title>Authentication required</title></head>
<body>
<p>{{.Message}}</p>
<p>Authentication required. Please enter your credentials:</p>
<form method="post" action="{{.ActionURL}}">
<p>Security token: <input type="text" name="ib-auth-token" value=""/></p>
<input type="submit" value="Submit"/>
</form>
</body>
</html>
`
