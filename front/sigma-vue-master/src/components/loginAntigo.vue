<template>
<div>
  <head>
    <meta charset="UTF-8" />
    <title>Formulário de Login e Registro com HTML5 e CSS3</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  </head>
  <body>
    <div class="container">
      <a class="links" id="paracadastro"></a>
      <a class="links" id="paralogin"></a>

      <div class="content">
        <!--FORMULÁRIO DE LOGIN-->
        <div id="login">
          <form method="post" action>
            <h1>Login</h1>
            <p>
              <label for="nome_login">Seu login</label>
              <input
                ref="login"
                id="nome_login"
                v-model="login"
                required="required"
                type="text"
                placeholder="ex. contato@htmlecsspro.com"
              />
            </p>

            <p>
              <label for="senha_login">Sua senha</label>
              <input
                id="senha_login"
                v-model="senha"
                required="required"
                type="password"
                placeholder="ex. senha"
              />
            </p>
            <p>
              <input v-on:click.prevent="realizaLogin" type="submit" value="Logar" />
            </p>
          </form>
        </div>
      </div>
    </div>
  </body>
</div>
</template>

<script>
export default {
  data() {
    return {
      login: "",
      senha: ""
    };
  },
  mounted() {
    this.focusInput();
  },
  methods: {
    realizaLogin() {
      console.log(
       
      );
      this.$store.commit("fazLogin", { login: this.login, senha: this.senha });
      let usuario = this.$store.getters.getCredencial;
      if (usuario.login == "") {
        alert("Usuário não logado");
        this.$refs.login.focus();
      } else {
        this.$router.replace("/")
      }

    },
    focusInput() {
      console.log("[login.vue|focusInput]");
      this.$refs.login.focus();
      // self.$refs.login.focus()
    }
  }
};
</script>



<style>
*,
*:before,
*:after {
  margin: 0;
  padding: 0;
  font-family: Arial, sans-serif;
}

/* remove a linha dos links */
a {
  text-decoration: none;
}

/* esconde as ancoras da tela */
a.links {
  display: none;
}

.content {
  width: 500px;
  min-height: 560px;
  margin: 0px auto;
  position: relative;
}

h1 {
  font-size: 48px;
  color: #066a75;
  padding: 2px 0 10px 0;
  font-family: Arial, sans-serif;
  font-weight: bold;
  text-align: center;
  padding-bottom: 30px;
}

h1:after {
  content: " ";
  display: block;
  width: 100%;
  height: 2px;
  margin-top: 10px;
  background: -webkit-linear-gradient(
    left,
    rgba(147, 184, 189, 0) 0%,
    rgba(147, 184, 189, 0.8) 20%,
    rgba(147, 184, 189, 1) 53%,
    rgba(147, 184, 189, 0.8) 79%,
    rgba(147, 184, 189, 0) 100%
  );
  background: linear-gradient(
    left,
    rgba(147, 184, 189, 0) 0%,
    rgba(147, 184, 189, 0.8) 20%,
    rgba(147, 184, 189, 1) 53%,
    rgba(147, 184, 189, 0.8) 79%,
    rgba(147, 184, 189, 0) 100%
  );
}

p {
  margin-bottom: 15px;
}

.content p:first-child {
  margin: 0px;
}

label {
  color: #405c60;
  position: relative;
}

/* placeholder */
::-webkit-input-placeholder {
  color: #bebcbc;
  font-style: italic;
}

input:-moz-placeholder,
textarea:-moz-placeholder {
  color: #bebcbc;
  font-style: italic;
}
input {
  outline: none;
}

/*estilo dos input,  menos o checkbox */
input:not([type="checkbox"]) {
  width: 95%;
  margin-top: 4px;
  padding: 10px;
  border: 1px solid #b2b2b2;

  -webkit-border-radius: 3px;
  border-radius: 3px;

  -webkit-box-shadow: 0px 1px 4px 0px rgba(168, 168, 168, 0.6) inset;
  box-shadow: 0px 1px 4px 0px rgba(168, 168, 168, 0.6) inset;

  -webkit-transition: all 0.2s linear;
  transition: all 0.2s linear;
}

/*estilo do botão submit */
input[type="submit"] {
  width: 100% !important;
  cursor: pointer;
  background: rgb(61, 157, 179);
  padding: 8px 5px;
  color: #fff;
  font-size: 20px;
  border: 1px solid #fff;
  margin-bottom: 10px;
  text-shadow: 0 1px 1px #333;

  -webkit-border-radius: 5px;
  border-radius: 5px;

  transition: all 0.2s linear;
}

/*estilo do botão submit no hover */
input[type="submit"]:hover {
  background: #4ab3c6;
}

.link {
  position: absolute;
  background: #e1eaeb;
  color: #7f7c7c;
  left: 0px;
  height: 20px;
  width: 440px;
  padding: 17px 30px 20px 30px;
  font-size: 16px;
  text-align: right;
  border-top: 1px solid #dbe5e8;

  -webkit-border-radius: 0 0 5px 5px;
  border-radius: 0 0 5px 5px;
}

.link a {
  font-weight: bold;
  background: #f7f8f1;
  padding: 6px;
  color: rgb(29, 162, 193);
  margin-left: 10px;
  border: 1px solid #cbd518;

  -webkit-border-radius: 4px;
  border-radius: 4px;

  -webkit-transition: all 0.4s linear;
  transition: all 0.4s linear;
}

.link a:hover {
  color: #39bfd7;
  background: #f7f7f7;
  border: 1px solid #4ab3c6;
}

#login {
  position: absolute;
  top: 0px;
  width: 88%;
  padding: 18px 6% 60px 6%;
  margin: 0 0 35px 0;
  background: #f7f7f7;
  border: 1px solid rgba(147, 184, 189, 0.8);

  -webkit-box-shadow: 5px;
  border-radius: 5px;

  -webkit-animation-duration: 0.5s;
  -webkit-animation-timing-function: ease;
  -webkit-animation-fill-mode: both;

  animation-duration: 0.5s;
  animation-timing-function: ease;
  animation-fill-mode: both;
}
</style>