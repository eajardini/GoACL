import Vue from 'vue';
import App from './App.vue';
import router from './router';
// import store from './store'
import Button from 'primevue/button';
import Breadcrumb from 'primevue/breadcrumb';
import Dialog from 'primevue/dialog';
import Dropdown from 'primevue/dropdown';
import InputText from 'primevue/inputtext';
import MegaMenu from 'primevue/megamenu';
import Menu from 'primevue/menu';
import Menubar from 'primevue/menubar';
import Message from 'primevue/message';
import Panel from 'primevue/panel';
import PanelMenu from 'primevue/panelmenu';
import Password from 'primevue/password';
import ToastService from 'primevue/toastservice';
import Tooltip from 'primevue/tooltip';
import Tree from 'primevue/tree';
import TreeTable from 'primevue/treetable';

import 'primevue/resources/themes/nova-light/theme.css';
import 'primevue/resources/primevue.min.css';
import 'primeflex/primeflex.css';
import 'primeicons/primeicons.css';
import 'prismjs/themes/prism-coy.css';
import '@fullcalendar/core/main.min.css';
import '@fullcalendar/daygrid/main.min.css';
import '@fullcalendar/timegrid/main.min.css';
import './assets/layout/layout.scss';

import "./plugins/axios"


Vue.use(ToastService);
Vue.directive('tooltip', Tooltip);

Vue.config.productionTip = false;

Vue.component('Breadcrumb', Breadcrumb);
Vue.component('Button', Button);
Vue.component('Dialog', Dialog);
Vue.component('Dropdown', Dropdown);
Vue.component('InputText', InputText);
Vue.component('MegaMenu', MegaMenu);
Vue.component('Menu', Menu);
Vue.component('Menubar', Menubar);
Vue.component('Message', Message);
Vue.component('Panel', Panel);
Vue.component('PanelMenu', PanelMenu);
Vue.component('Password', Password);
Vue.component('Tree', Tree);
Vue.component('TreeTable', TreeTable);

new Vue({
	router,
	// store,
	render: h => h(App)
}).$mount('#app');
