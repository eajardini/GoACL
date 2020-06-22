<template>
  <div class="p-grid">
    <div class="p-col-12">
      <div class="card">
        <h1>Bem vindo</h1>
        <p>Sistema Modelo.</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
     
    };
  },
  async mounted() {
    await this.$acl
      .get("/acl/MontaMenu")
      .then(resp => {
        // this.menu = resp.data.resposta
        this.menu = resp.data;
        //   this.menu = [{"label":"Financeiro","items":[{"label":"Contas a Pagar","items":[{"label":"Cadastro"}]}]},{"label":"CRM","items":null}]
      })
      .catch(error => {
        // handle error
        console.log("Erro de retorno:" + error);
        // console.log("Erro de dados(data):" + error.response.data);
        // console.log("Erro do status:" + error.response.status);
        // console.log("Erro headers:" + error.response.headers);
        this.verificaSeEstaLogado();
      });
  },   
  methods: {
    verificaSeEstaLogado() {
      let usuario = this.$store.getters.getCredencial;       
      if (usuario.logado == true) {        
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