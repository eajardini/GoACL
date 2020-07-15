<template>
  <div class="layout-topbar">
    <button class="p-link layout-menu-button" @click="onMenuToggle">
      <span class="pi pi-bars"></span>
    </button>
    <div class="layout-topbar-icons">
      <button class="p-link">
        <span class="layout-topbar-item-text">Eventos</span>
        <span class="layout-topbar-icon pi pi-calendar"></span>
        <span class="layout-topbar-badge">1</span>
      </button>
      <button class="p-link" @click="sairDoSistema">
        <span class="layout-topbar-item-text">Sair</span>
        <span class="layout-topbar-icon pi pi-sign-out"></span>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    onMenuToggle(event) {
      this.$emit("menu-toggle", event);
    },
    sairDoSistema() {
      this.$parent.mostraLeftBar = false;
      this.$parent.mostraTopBar = false;
      this.$parent.staticMenuInactive = true;
      this.$parent.nomeDoUsuarioApp = "";        
      this.$cookies.remove("credencial");
      this.$cookies.remove("token");
      this.$cookies.remove("bearer");
     this.$router.push("/login");
    }
  }
};
</script>