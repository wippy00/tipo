<script>
export default {
	data: function () {
		return {
			error: null,
			errormsg: null,
			loading: false,
		    
            conversations: null
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
        async fetchConversations() {
            this.loading = true
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations", {
                    headers: {
                        authorization: auth_id
                    }
                })
                conversations = response.data

                for (let i = 0; i < conversations.length; i++) {
                    let userData = await this.fethcUser(conversations[i].last_message.author)
                    conversations[i].last_message.author = userData
                }
                this.conversations = conversations
                // return conversations
            } catch (e) {
                this.error = e.toString()
            }
            this.loading = false
        },
        async fethcUser(id) {
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/user/"+id, {
                    headers: {
                        authorization: auth_id
                    }
                })
                return response.data
            } catch (e) {
                this.error = e.toString()
            }
        }
	},
	async mounted() {
        if (sessionStorage.getItem('logged_in') === false || sessionStorage.getItem('logged_in') === null) {
            console.log("Not logged in")
            this.$router.push('/')
        }

		this.refresh();
        await this.fetchConversations();
    }
}
</script>

<template>
	<div class="container">
		<!-- <ErrorMsg v-if="error" :msg="errormsg"></ErrorMsg> -->

        <h1 v-if="loading">Loading...</h1>

        <div v-for="item in conversations" class="row card my-4" >
            <div class="row g-0">

                <div class="col-md-1 col-2">
                    <!-- <img src="..." class="img-thumbnail" alt="..."> -->
                    <img src="https://placehold.co/100" class="img-thumbnail" alt="...">
                    
                </div>

                <div class="col-md-11 col-10">
                    <div class="card-body">
                        <RouterLink :to="'/conversations/' + item.id">
                            <h5 class="card-title text-capitalize">{{ item.name }}</h5>
						</RouterLink>
                        <p class="card-text text-capitalize">{{ item.last_message.author.name +": " }}<small class="text-body-secondary">{{ item.last_message.content }}</small></p>
                    </div>
                </div>

            </div>
        </div>
		
	</div>
</template>

<style>
.container {
	height: 100vh;
}
</style>
