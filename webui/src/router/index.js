import {createRouter, createWebHashHistory} from 'vue-router'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component:  () => import('@/views/HomeView.vue')},
		{path: '/conversations', component: () => import('@/views/ConversationsView.vue')},
		{path: '/conversations/:id', component: () => import('@/views/ChatView.vue')},
		
		
		{path: '/profile', component: () => import('@/views/ProfileView.vue')},
		// {path: '/some/:id/link', component: HomeView},
	]
})

export default router
