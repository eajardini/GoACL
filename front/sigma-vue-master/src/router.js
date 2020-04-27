import Vue from 'vue';
import Router from 'vue-router';
import Dashboard from './components/Dashboard.vue';

Vue.use(Router);

export default new Router({
	mode: "history",
	routes: [
		{
			path: '/',
			name: 'dashboard',
			component: Dashboard
		},
	
		{
			path: '/empty',
			name: 'empty',
			component: () => import('./components/EmptyPage.vue')
		},
	
	],
	scrollBehavior() {
		return {x: 0, y: 0};
	}
});