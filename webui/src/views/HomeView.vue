<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			error: null,
			username: "",
			id: 0,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async loginHandler(event) {
			event.preventDefault()

			if (this.username === "") {
				this.error = "Username cannot be empty.";
				return;
			}
			this.error = null
			try {
				let response = await this.$axios.post("/login", { name: this.username })
				console.log(response.data)
				await sessionStorage.setItem("id", response.data.id);
				await sessionStorage.setItem("username", response.data.name);
				// await sessionStorage.setItem("photo", response.data.photo);
				this.$router.push({ path: '/conversations' })
			} catch (event) {
				if (event.response && event.response.status === 400) {
					this.error = "Username should has a length between 3 - 16";
				} else if (event.response && event.response.status === 500) {
					this.error = "An internal error occurred, please try again later.";
				} else {
					this.error = event.toString();
				}
				setTimeout(() => {
					this.error = null;
				}, 5000);
			}

		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div class="container">
		<!-- <div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div> -->

		<ErrorMsg v-if="error" :msg="errormsg"></ErrorMsg>


		<div class="p-2">
			<h1 class="text-center">Login</h1>
			<form @submit.prevent="loginHandler">
				<div class="mb-3">
					<label for="username" class="form-label">Username</label>
					<input type="text" class="form-control" id="username" v-model="username"
						aria-describedby="usernameHelp">
				</div>
				<button type="submit" class="btn btn-primary">Submit</button>
			</form>
		</div>

		<!-- <h1 v-if="sessionStorage.getItem('username')">
			{{ username }}
		</h1> -->
	</div>
</template>

<style>
.container {
	height: 100vh;
}
</style>
