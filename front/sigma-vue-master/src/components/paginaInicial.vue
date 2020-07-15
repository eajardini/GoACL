<template>
  <div class="p-grid">
    <div class="p-col-12">
      <div class="card">
        <h1>Bem vindo</h1>
        <p>Sistema Modelo</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      token: "",
      bearer: ""
    };
  },
  beforeCreate: function() {
    if (this.$cookies.get("token") == null) {
      this.$router.push("/login");
    }
    this.token = this.$cookies.get("token");
    // console.log("[paginainicial.vue|beforeCreate] Valor token:" + this.token);
  },
  // async mounted() {
  mounted() {
    this.setaAmbiente();
    this.$acl.defaults.headers.common["Authorization"] = this.$cookies.get(
      "bearer"
    );
    this.$acl
      .get("/acl/MontaMenu")
      .then(resp => {
        this.$parent.menu = resp.data;

        console.log(
          "[paginainicial.vue|monted] Valor resp:" + resp.data.ItemsNivel1Locais
        );
      })
      .catch(error => {
        console.log("Erro de retorno:" + error);
        // console.log("Erro de dados(data):" + error.response.data);
        // console.log("Erro do status:" + error.response.status);
        // console.log("Erro headers:" + error.response.headers);

        // this.verificaSeEstaLogado();
      });
  },
  methods: {
    setaAmbiente() {
      // this.token = sessionStorage.getItem("token");
      // this.token = this.$session.get("token")
      // console.log("[paginainicial.vue|setaToken] Valor do token:" + this.token )
      // if (this.token == null) {
      //   this.$router.push("/login");
      // } else {

      // this.$parent.menu = [
      //   { label: "Dashboard", icon: "pi pi-fw pi-home", to: "/" },
      //   {
      //     label: "Menu Modes",
      //     icon: "pi pi-fw pi-cog",
      //     items: [
      //       {
      //         label: "Static Menu",
      //         icon: "pi pi-fw pi-bars",
      //         command: () => (this.layoutMode = "static")
      //       }
      //     ]
      //   }
      // ];

      this.$parent.nomeDoUsuarioApp = "Pipoca";
      this.$parent.mostraLeftBar = true;
      this.$parent.mostraTopBar = true;
      this.$parent.staticMenuInactive = false;

      // }
    },

    setCabecalho() {
      this.bearer = "Bearer ";
      this.token = this.$cookies.get("token");
      // Vue.prototype.$http.defaults.headers.common['Authorization'] = this.$cookie.get('tokencliente')
      this.bearer = this.bearer + this.token;
      return this.bearer;
    },

    verificaSeEstaLogado() {
      let usuario = true;
      // = this.$store.getters.getCredencial;
      // console.log("[paginaInicial.vue| mounted] valor token:" + usuario.token)
      // console.log("[paginaInicial.vue| mounted] valor token:" + usuario.token)

      if (usuario == true) {
        this.$parent.menu = [
          { label: "Dashboard", icon: "pi pi-fw pi-home", to: "/" },
          {
            label: "Menu Modes",
            icon: "pi pi-fw pi-cog",
            items: [
              {
                label: "Static Menu",
                icon: "pi pi-fw pi-bars",
                command: () => (this.layoutMode = "static")
              }
            ]
          }
        ];
        this.$parent.nomeDoUsuarioApp = "Pipoca";
        this.$parent.mostraLeftBar = true;
        this.$parent.mostraTopBar = true;
        this.$parent.staticMenuInactive = false;
      } else {
        this.$router.push("/login");
      }
    }
  }
};
</script>

<style scoped>
</style>