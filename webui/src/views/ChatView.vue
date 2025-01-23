<script>
export default {
	data: function () {
		return {
			error: null,
			errormsg: null,
			loading: false,
		    
            messages: null,
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
        async fetchMessages(conversations_id) {
            this.loading = true
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations/"+conversations_id+"/messages", {
                    headers: {
                        authorization: auth_id
                    }
                })
                this.messages = response.data
                // return messages
            } catch (e) {
                this.error = e.toString()
            }
            this.loading = false
        },
        async fetchConversations(conversations_id) {
            this.loading = true
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations/"+conversations_id, {
                    headers: {
                        authorization: auth_id
                    }
                })
                this.conversations = response.data
                // return conversations
            } catch (e) {
                this.error = e.toString()
            }
            this.loading = false
        },
        async fetchUser(user_id) {
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/user/"+user_id, {
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
        this.auth_id = sessionStorage.getItem('id');
		this.refresh();
        await this.fetchConversations(this.$route.params.id);
        await this.fetchMessages(this.$route.params.id);
   
        for (let i = 0; i < this.messages.length; i++) {
            let userData = await this.fetchUser(this.messages[i].author);
            console.log(userData)
            this.messages[i].author = userData
        }  
        
    }
}
</script>

<template>
    <div class="container">
        <!-- <h1> {{ conversations.name }}</h1> -->

        <div v-for="message in messages">
            <div v-if="message.author.id == auth_id" class="card my-4 bg-light offset-md-7 col-5">
                <h5 class="card-title ms-3 mt-1">{{ message.author.name }}</h5>
                <div class="card-body">
                    <img v-if="message.photo" :src="'data:image/jpeg;base64,' + message.photo" class="card-img-top rounded-3" alt="...">
                    <p class="card-text mt-2">{{ message.content }}</p>
                </div>
                <small class="text-end p-2">{{ message.timestamp }}</small>
            </div>
            <div v-else class="card my-4 bg-light col-5">
                <h5 class="card-title ms-3 mt-1">{{ message.author.name }}</h5>
                <div class="card-body">
                    <img v-if="message.photo" :src="'data:image/jpeg;base64,' + message.photo" class="card-img-top rounded-3" alt="...">
                    <p class="card-text mt-2">{{ message.content }}</p>
                </div>
                <small class="text-end p-2">{{ message.timestamp }}</small>
            </div>
            
        </div>

    </div>
</template>