<template>
  <div class="content">
    <div id="loginFormulario">
      <form style="text-align: center" action>
        <h1>Login</h1>
        <p>
          <label for="login">
            <span>Seu login</span>
            <span></span>
          </label>
          <InputText type="text" v-model="login" placeholder="Informe seu login" />
        </p>

        <p>
          <label for="senha_login">
            <span>Sua senha</span>
          </label>
          <InputText type="password" v-model="senha" placeholder="Informe sua senha" />
        </p>
        <p>
          <Button label="Logar" @click.prevent="realizaLogin" />
        </p>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      login: "",
      senha: "",
      autenticado: false,
      token: ""
    };
  },
  mounted() {
    this.$parent.nomeDoUsuarioApp = "";
    this.$parent.mostraLeftBar = false;
    this.$parent.mostraTopBar = false;
    this.$parent.staticMenuInactive = true;
  },
  methods: {
    async realizaLogin() {
      var md5 = require("md5");
      var senhaLocal = "";
      const formData = new FormData();
      senhaLocal = md5(this.senha);
      formData.append("username", this.login);
      formData.append("password", senhaLocal);

      
      await this.$acl
        // this.$acl
        .post("login", formData)
        .then(resp => {
          this.token = resp.data.token;          
          console.log("[login.vue|realizaLogin]Valor senha MD5:" + senhaLocal);
          sessionStorage.setItem("senha", senhaLocal);
          sessionStorage.setItem("token", this.token);         
          this.$router.push("/");
        })
        .catch(error => {
          // handle error
          console.log(error);
          if (error.message.includes("401")) {
            alert("Usuário ou senha inválidos");
            self.$refs.login.focus();
          } else {
            let erroSTR = String(error)
            if (erroSTR.includes("Network Error")) {
              
              alert("Não foi possível alcançar o servidor. Contate o suporte do sistema");
            } else {
              alert("Error geral de login:" + error);
            }
          }
        });
    },

    focusInput() {
      console.log("[login.vue|focusInput]");
      // this.$refs.login.focus();
      // self.$refs.login.focus()
    }
  }
};
</script>

<style scoped>
/* *,
*:before,
*:after {
  margin: 0;
  padding: 0;
  font-family: Arial, sans-serif;
} */
</style>