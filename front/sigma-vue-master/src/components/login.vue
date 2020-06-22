<template>
  <div class="p-grid">
    <div class="p-col-12">
      <div class="card">
        <div class="container">
          <a class="links" id="paracadastro"></a>
          <a class="links" id="paralogin"></a>

          <div class="content">
            <!--FORMULÁRIO DE LOGIN-->
            <div id="login">
              <form style="text-align: center" action>
                <h1>Login</h1>
                <p>
                  <label for="nome_login"><span> Seu login </span><span> </span></label>
                  <InputText type="text" v-model="login" placeholder="Informe seu login" />
                </p>

                <p>
                  <label for="senha_login"> <span>Sua senha </span>  </label>
                  <InputText type="password" v-model="senha" placeholder="Informe sua senha" />
                </p>
                <p>                  
									<Button label="Logar" @click="realizaLogin"/>
                </p>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
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
    // this.focusInput();
  },
  methods: {
    realizaLogin() {
      console.log();
      this.$store.commit("fazLogin", { login: this.login, senha: this.senha, logado: true });
      let usuario = this.$store.getters.getCredencial;
      if (usuario.logado == true) {
        this.$router.replace("/");
       
      } else {
         alert("Usuário não logado");
        this.$refs.login.focus();
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

<style scoped>
*,
*:before,
*:after {
  margin: 0;
  padding: 0;
  font-family: Arial, sans-serif;
}
</style>