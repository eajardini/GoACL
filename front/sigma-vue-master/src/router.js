import Vue from 'vue';
import Router from 'vue-router';
import Dashboard from './components/Dashboard.vue';
import login from './components/login.vue';
import paginaInicial from './components/paginaInicial.vue';

Vue.use(Router);

export default new Router({
	mode: "history",
	routes: [
		{
			path: '/dash',
			name: 'dashboard',
			component: Dashboard
		},	
		{
			path: '/',
			name: 'paginaInicial',
			component: paginaInicial
		},
		{
			path: '/login',
			name: 'login',
			component: login
		},
	
	],
	scrollBehavior() {
		return {x: 0, y: 0};
	}
});