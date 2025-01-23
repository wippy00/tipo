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
                this.conversations = response.data
                return conversations
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
		this.refresh(),
        await this.fetchConversations();
   
        for (let i = 0; i < this.conversations.length; i++) {
            let userData = await this.fethcUser(this.conversations[i].last_message.author)
            this.conversations[i].last_message.author = userData
        }  
        
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
                        <h5 class="card-title">{{ item.name }}</h5>
                        <p class="card-text"><small class="text-body-secondary">{{ item.last_message.author.name +": "+item.last_message.content }}</small></p>
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
