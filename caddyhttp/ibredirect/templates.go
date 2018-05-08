package ibredirect

const authPageTemplate = `
<html xmlns="http://www.w3.org/1999/xhtml">    
<head><title>Authentication required</title></head>
<body>
<p>{{.Message}}</p>
<p>Authentication required. Please enter your credentials:</p>
<form method="post" action="{{.ActionURL}}">
<p>URL to browse on: <input type="text" name="ib-origin-url" value="{{.OriginURL}}"/></p>
<p>Security token: <input type="text" name="ib-auth-token" value=""/></p>
<input type="submit" value="Submit"/>
</form>
</body>
</html>
`

const stubPage = `
<html xmlns="http://www.w3.org/1999/xhtml">    
  <head>      
    <title>Test</title>    
  </head>    
  <body>Debug</body>  
</html>     
`

const redirectPage = `
<html xmlns="http://www.w3.org/1999/xhtml">    
  <head>      
    <meta http-equiv="refresh" content="0;URL='%s'" />    
  </head>    
  <body>Redirecting, please wait...</body>  
</html>     
`
