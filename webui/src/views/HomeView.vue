<script>
export default {
	data: function () {
		return {
			error: null,
			errormsg: null,
			loading: false,

			logged_in: false,
			name: ""
		}
	},
	methods: {
		async loginHandler(event) {
			event.preventDefault()

			if (this.name === "") {
				this.error = "name cannot be empty.";
				return;
			}
			this.error = null

			try {
				let response = await this.$axios.post("/login", { name: this.name })
				
				sessionStorage.setItem("logged_in", true);
				sessionStorage.setItem("id", response.data.id);
				sessionStorage.setItem("name", response.data.name);
				localStorage.setItem("photo", response.data.photo);

				this.$emit("login-success")
				this.$router.push({ path: '/conversations' })

			} catch (event) {
				if (event.response && event.response.status === 400) {
					this.error = "name should has a length between 3 - 16";
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
		if (sessionStorage.getItem('logged_in') === "true") {
            this.$router.push('/conversations')
        }

		this.logged_in = sessionStorage.getItem('logged_in');
		console.log(this.logged_in)
	}
}
</script>

<template>
	<div class="container">

		<ErrorMsg v-if="error" :msg="errormsg"></ErrorMsg>

		<div v-if="logged_in !== 'true'">
			<div class="p-2">
				<h1 class="text-center">Login</h1>
				<form @submit.prevent="loginHandler">
					<div class="mb-3">
						<label for="name" class="form-label">name</label>
						<input type="text" class="form-control" id="name" v-model="name"aria-describedby="nameHelp">
					</div>
					<button type="submit" class="btn btn-primary">Submit</button>
				</form>
			</div>
		</div>
		<div v-else>
			<h1 class="text-center">Welcome {{ name }}</h1>
		</div>
		

		

		

	</div>
</template>

<style>

</style>
