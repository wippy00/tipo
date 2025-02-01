<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data: function () {
		return {
			error: null,
			errormsg: null,
			loading: false,
		    
			auth_id: null,
			auth_name: null,
			auth_photo: null
        }
        
	},
	methods: {
		async logout() {
			sessionStorage.removeItem('id');
			sessionStorage.removeItem('name');
			localStorage.removeItem('photo');
			
			sessionStorage.setItem("logged_in", false);
			this.refreshData();

			this.$router.push('/');
		},
		refreshData() {
			this.auth_id = sessionStorage.getItem('id');
			this.auth_name = sessionStorage.getItem('name');
			this.auth_photo = localStorage.getItem('photo');
		}
	},
	mounted() {
		this.refreshData();
	}
}
</script>

<template>
	<header>
		<nav class="navbar navbar-expand-lg bg-dark fixed-top">
			<div class="container-fluid">
				<a class="navbar-brand">Wasa Text</a>
				<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
				</button>
				<div class="collapse navbar-collapse" id="navbarNav">
					<ul class="navbar-nav">
						<li class="nav-item">
						<RouterLink to="/" class="nav-link ">
							Home
						</RouterLink>
						</li>
						<li class="nav-item">
						<RouterLink to="/profile" class="nav-link ">
							Profile
						</RouterLink>
						</li>
					</ul>
				</div>
				<div v-if="auth_id" class="d-flex profilo">
					<RouterLink to="/profile" class="nav-link text-light">
						<img v-if="auth_photo == null" :src="'data:image/jpeg;base64,' + auth_photo" width="42" height="42" class="rounded-circle" style="object-fit: cover;">
						<img v-else :src="'https://placehold.co/100x100/orange/white?text=' + auth_name" width="42" height="42" class="rounded-circle" style="object-fit: cover;">
					</RouterLink>
					<p class="text-capitalize ms-2 pt-2 d-block">Hi {{ auth_name }} </p>
					<button type="button" class="btn btn-primary ms-2" @click="logout">Logout</button>
				</div>
				<!-- <form class="d-flex" role="search">
					<input class="form-control me-2" type="search" placeholder="Search" aria-label="Search">
					<button class="btn btn-outline-success" type="submit">Search</button>
				</form> -->
			</div>
			<hr>
		</nav>
	</header>
	
	<div class="mt-5 p-3"></div>
	
	<main>
		<RouterView @login-success="refreshData" />
	</main>

</template>