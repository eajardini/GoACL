export	default	{
  fazLogin:	(credencial,	loginInformado)	=>	{
    console.log("[logins.js|fazLogin] login:" + loginInformado.login + " senha: " + loginInformado.senha)
   // if ((loginInformado.login == "admin") && (loginInformado.senha == "123")) {
    credencial.usuario.login = loginInformado.login
    credencial.usuario.senha = loginInformado.senha
    credencial.usuario.logado = true
    //}    
  },
  
}