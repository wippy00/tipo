<script>
export default {
	data: function () {
		return {
			error: null,
			errormsg: null,
			loading: false,
            
            auth_id: null,
            auth_photo: null,

            messages: [],
            conversations: null,

            message_input: "",
            photo_input: null,

            refreshInterval: null
            
        }
        
	},
	methods: {
		async refresh() {
			this.errormsg = null;
			
            this.auth_id = sessionStorage.getItem('id');
            // await this.fetchConversations(this.$route.params.id);
            await this.fetchMessages(this.$route.params.id);

            // this.$nextTick(() => {
            //     this.scrollToBottom();
            // });
		},
        async fetchMessages(conversations_id) {
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations/"+conversations_id+"/messages", {
                    headers: {
                        authorization: auth_id
                    }
                })
                let messages = response.data
                
                for (let i = 0; i < messages.length; i++) {
                    let userData = await this.fetchUser(messages[i].author);
                    messages[i].author = userData
                }
                
                this.messages = messages

                // return messages
            } catch (e) {
                this.error = e.toString()
            }
        },
        async fetchConversations(conversations_id) {
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations/" + conversations_id, {
                    headers: {
                        authorization: auth_id
                    }
                })
                this.conversations = response.data

            } catch (e) {
                this.error = e.toString()
            }
        },
        async fetchUser(user_id) {
            this.error = null

			let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/users/"+user_id, {
                    headers: {
                        authorization: auth_id
                    }
                })
                return response.data
            } catch (e) {
                this.error = e.toString()
            }
        },
        async sendMessage(event) {
            event.preventDefault()
            if (this.message_input === "") {
				this.error = "Message cannot be empty.";
				return;
			}
            this.error = null

            let auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.post("/conversations/"+this.$route.params.id+"/messages", {
                    text: this.message_input
                }, {
                    headers: {
                        authorization: auth_id
                    }
                })

                this.messages.push(response.data)
                this.refresh()

            } catch (e) {
                this.error = e.toString()
            }

            this.$nextTick(() => {
                this.scrollToBottom();
            });
        },
        scrollToBottom() {
            window.scrollTo(0, document.body.scrollHeight);
        }
	},
	async mounted() {
        if (sessionStorage.getItem('logged_in') !== "true") {
            console.log("Not logged in")
            this.$router.push('/')
        }
        
        this.auth_id = sessionStorage.getItem('id');

        await this.fetchMessages(this.$route.params.id)

        await this.fetchConversations(this.$route.params.id);
        
        this.$nextTick(() => {
            this.scrollToBottom();
        });

        this.refreshInterval = setInterval(() => { // Salva l'ID dell'intervallo
            this.refresh();
        }, 1000);
    },
    unmounted() {
        clearInterval(this.refreshInterval)
    }
}
</script>

<template>
    <div class="container">
        <!-- <div class="card p-2 bg-body-tertiary col-12">
            <div class="d-flex">
                <img v-if="conversations.photo" :src="'data:image/jpeg;base64,' + conversations.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + conversations.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                <h1 class="ms-2 text-capitalize"> {{ conversations.name }}</h1>
            </div>
            <ul class="list-inline mt-2 ms-2">
                <li class="list-inline-item text-capitalize" v-for="user in conversations.participants">{{ user.name  + " -" }}</li>
                <li class="list-inline-item">...</li>
            </ul>
        </div> -->
        <div v-for="message in messages">
            <!-- Se sono io -->
            <div v-if="message.author.id == auth_id" class="card my-4 bg-body-tertiary offset-md-7 col-5">

                <div class="d-flex">
                    <img v-if="message.author.photo" :src="'data:image/jpeg;base64,' + message.author.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.author.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <h5 class="card-title ms-2 mt-3 text-capitalize"> {{ message.author.name }} </h5>
                </div>

                <div class="card-body">
                    <img v-if="message.photo" :src="'data:image/jpeg;base64,' + message.photo" class="card-img-top rounded-3" alt="...">
                    <p class="card-text mt-2">{{ message.text }}</p>
                </div>
                <small class="text-end p-2">{{ message.timestamp }}</small>

            </div>
            
            <!-- Se sono gli altri -->
            <div v-else class="card my-4 bg-body-tertiary col-5">

                <div class="d-flex">
                    <img v-if="message.author.photo" :src="'data:image/jpeg;base64,' + message.author.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.author.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <h5 class="card-title ms-2 mt-3 text-capitalize"> {{ message.author.name }} </h5>
                </div>

                <div class="card-body">
                    <img v-if="message.photo" :src="'data:image/jpeg;base64,' + message.photo" class="card-img-top rounded-3" alt="...">
                    <p class="card-text mt-2">{{ message.text }}</p>
                </div>
                <small class="text-end p-2">{{ message.timestamp }}</small>
            
            </div>

        </div>

        <form @submit.prevent="sendMessage" class="input-group mb-3">
            <input v-model="message_input" id="message_input" type="text" class="form-control" placeholder="Type a message" aria-label="Type a message" aria-describedby="message_input">
            <button class="btn btn-primary" type="submit" id="send">Send</button>
        </form>

        <button  @click="refresh" class="btn btn-primary" type="refresh" id="send">Refresh</button>

            
    </div>
</template>

<style>
/* .container{
    height: 90vh;
} */

</style>