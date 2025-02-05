<script>
// import ErrorMsg from '@/components/ErrorMsg.vue'
import ChatHeader from '@/components/ChatHeader.vue';

import ConversationCard from '@/components/ConversationCard.vue';
import ConversationUserPhoto from '@/components/ConversationUserPhoto.vue';

import MessageReactionPopup from '@/components/MessageReactionPopup.vue';

export default {
    components: {
        ChatHeader,
        ConversationCard,
        ConversationUserPhoto,

        MessageReactionPopup,
    },
	data: function () {
		return {
			error: null,
			errormsg: null,
			loading: false,
            
            auth_id: null,
            auth_photo: null,

            messages: [],
            conversation: {},
            allConversations: [],
            participants: {},

            message_text: "",
            message_photo: null,
            message_forward: 0,
            message_reply: 0,

            showForwardMessage: false,
            showMessageReactionPopup: false,
            messageReactionPopupData: null,

            forward_conversation_id: null,

            replyMessage_data: null,

            refreshInterval: null
            
        }
        
	},
	methods: {
        async refresh() {
			// this.errormsg = null;
			
            this.auth_id = sessionStorage.getItem('id');
            // await this.fetchconversation(this.$route.params.id);
            await this.fetchMessages(this.$route.params.id);

            // this.$nextTick(() => {
            //     this.scrollToBottom();
            // });
		},
        async fetchAll(conversation_id) {
            await this.fetchconversation(conversation_id);
            await this.fetchMessages(conversation_id);
        },
        async fetchMessages(conversation_id) {
            // this.error = null

			const auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations/"+conversation_id+"/messages", {
                    headers: {
                        authorization: auth_id
                    }
                })
                let messages = response.data
                
                for (let i = 0; i < messages.length; i++) {
                    let userData = this.getUser(messages[i].author);
                    messages[i].author = userData
                    

                    // expand forwoard
                    if (messages[i].forward != 0) {
                        if (messages[i].forward in this.participants) {
                            messages[i].forward = this.getUser(messages[i].forward);
                        } 
                        else {
                            messages[i].forward = await this.fetchUser(messages[i].forward);
                        }
                    }

                    // expand reactions
                    if (messages[i].reactions)  {
                        messages[i].reactions.forEach(item => {
                            item.user = this.getUser(item.user)
                        })
                    }
                    
                }

                var messages_dict = {}
                messages.forEach(message => messages_dict[message.id] = message);
                
                // expand repy message
                messages.forEach(message => {
                    if (message.reply != 0) {
                        message.reply = messages_dict[message.reply]
                    }
                });

                this.messages = messages

                // return messages
            } catch (e) {
                this.error = e.toString()
            }
        },
        async fetchconversation(conversation_id) {
            // this.error = null

			const auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations/" + conversation_id, {
                    headers: {
                        authorization: auth_id
                    }
                })
                this.conversation = response.data

                // get all participants
                for (let j = 0; j < this.conversation.participants.length; j++) {
                    this.participants[this.conversation.participants[j].id] = this.conversation.participants[j]
                }
                // console.log(this.participants)

            } catch (e) {
                this.error = e.toString()
            }
        },
        async fetchAllConversation() {
            this.loading = true
            this.error = null

			this.auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations", {
                    headers: {
                        authorization: this.auth_id
                    }
                })
                this.allConversations = response.data

                // console.log(this.conversations)
                
                // expand the last message author
                for (let i = 0; i < this.allConversations.length; i++) {
                    if (this.allConversations[i].last_message == 0) {
                        continue
                    }

                    for (let j = 0; j < this.allConversations[i].participants.length; j++) {
                        if (this.allConversations[i].participants[j].id == this.allConversations[i].last_message.author) { 
                            this.allConversations[i].last_message.author = this.allConversations[i].participants[j] 
                        }
                
                    }
                }
            
            } catch (e) {
                this.error = e.toString()
            }
            this.loading = false
        },
        getUser(user_id) {

            try {
                return this.participants[user_id]
            } catch (e) {
                this.error = e.toString()
            }
        },
        async fetchUser(user_id) {
            this.auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/users/"+user_id, {
                    headers: {
                        authorization: this.auth_id
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

            const auth_id = sessionStorage.getItem('id')

            const formData = new FormData();
            formData.append('text', this.message_input);
            if (this.message_photo) {
                formData.append('photo', this.message_photo);
            }

            if (this.message_reply) {
                formData.append('reply', this.message_reply);
            }

            try {
                let response = await this.$axios.post("/conversations/"+this.$route.params.id+"/messages", formData, {
                    headers: {
                        authorization: auth_id
                    }
                })

                this.message_input = ""
                this.message_photo = null
                this.message_reply = 0
                this.replyMessage_data = null

                this.refresh()

            } catch (e) {
                this.error = e.toString()
            }

            this.$nextTick(() => {
                this.scrollToBottom();
            });
        },
        async deleteMessage(event) {
            event.preventDefault()
            
            console.log(event.target.message_id.value)
            
            const message_id = event.target.message_id.value;
            const auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.delete("/conversations/"+ this.$route.params.id +"/messages/" + message_id, {
                    headers: {
                        authorization: auth_id
                    }
                })
                this.refresh()
                
            } catch (e) {
                this.error = e.toString()
            }
        },
        replyMessage(event) {
            event.preventDefault()
            
            // console.log(event.target.message_id.value)
            let message_id = event.target.message_id.value
            
            this.replyMessage_data = this.messages.find(message => message.id == message_id)

            this.message_reply = this.replyMessage_data.id
            // console.log(this.replyMessage_data)
        },

        showForwardMessageHandler(event){
            this.showForwardMessage = true
            this.fetchAllConversation()
            console.log(this.allConversations)
            event.preventDefault()
            let message_id = event.target.message_id.value
            this.message_forward = message_id
        },

        showMessageReactionPopupHandler(event) {
            event.preventDefault()
            console.log("showMessageReactionPopup")
            let message_id = event.target.message_id.value
            this.messageReactionPopupData = this.messages.find(message => message.id == message_id)
            this.showMessageReactionPopup = true
        },
        closeMessageReactionPopup() {
            this.showMessageReactionPopup = false;
        },

        forwardMessage(event) {
            event.preventDefault()
            let conversation_id = event.target.forward_conversation_id.value
            let user_id = this.message_forward
            const auth_id = sessionStorage.getItem('id')

            try {
                let response = this.$axios.post("/conversations/"+conversation_id+"/messages/"+user_id+"/forward", {}, {
                    headers: {
                        authorization: auth_id
                    }
                })

                this.allConversations = {}
                this.message_forward = null
                this.refresh()
                
            } catch (e) {
                this.error = e.toString()
            }

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

        await this.fetchAll(this.$route.params.id)

        // await this.fetchMessages(this.$route.params.id)

        // await this.fetchconversation(this.$route.params.id);
        
        this.$nextTick(() => {
            this.scrollToBottom();
        });

        this.refreshInterval = setInterval(() => { // Salva l'ID dell'intervallo
            this.refresh();
        }, 2000);
    },
    unmounted() {
        clearInterval(this.refreshInterval)
    }
}
</script>

<template>

    <!-- <ErrorMsg v-if="error" :msg="error"></ErrorMsg> -->

    
    <div class="container">

        <MessageReactionPopup v-if="showMessageReactionPopup" :message="messageReactionPopupData" @close="closeMessageReactionPopup"></MessageReactionPopup>

        <ConversationHeader :conversations="conversation" :auth_id="auth_id" />
        <ChatHeader :conversations="conversation" :auth_id="auth_id" />

        <div v-for="message in messages">






            <!-- Se sono io -->
            <div v-if="message.author.id == auth_id" class="card my-4 bg-body-tertiary offset-md-7 col-5">
                
                <!-- sezione risposta messaggio -->
                <div v-if="message.reply != 0">
                    <h6 class="m-2">Reply To:</h6>
                    <div class="card m-2 bg-body-tertiary p-2">
                        <div class="d-flex">
                            <img v-if="message.reply.author.photo" :src="'data:image/jpeg;base64,' + message.reply.author.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.reply.author.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <h6 class="card-title ms-2 mt-3 text-capitalize"> {{ message.reply.author.name }} </h6>
                        </div>
                        <div class="card-body">
                            <img v-if="message.reply.photo" :src="'data:image/jpeg;base64,' + message.reply.photo" class="card-img-top rounded-3" alt="...">
                            <p class="card-text mt-2">{{ message.reply.text }}</p>
                        </div>
                    </div>
                </div>


                <!-- sezione header messaggio -->
                <div class="d-flex">
                    <img v-if="message.author.photo" :src="'data:image/jpeg;base64,' + message.author.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.author.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <h5 class="card-title ms-2 mt-3 text-capitalize"> {{ message.author.name }} </h5>
                    <div class="ms-auto m-1 row">
                        <form @submit.prevent="deleteMessage" class="col">
                            <input type="hidden" name="message_id" :value="message.id">
                            <button class="btn btn-primary" type="submit">Delete</button>
                        </form>

                        <form @submit.prevent="replyMessage" class="col">
                            <input type="hidden" name="message_id" :value="message.id">
                            <button class="btn btn-primary" type="submit">Reply</button>
                        </form>

                        <form @submit.prevent="showForwardMessageHandler" class="col">
                            <input type="hidden" name="message_id" :value="message.id">
                            <button class="btn btn-primary" type="submit">Forward</button>
                        </form>

                        <form @submit.prevent="showMessageReactionPopupHandler" class="col">
                            <input type="hidden" name="message_id" :value="message.id">
                            <button class="btn btn-primary" type="submit">React</button>
                        </form>

                    </div>
                </div>

                <!-- sezione forward messaggio -->
                <div v-if="message.forward != 0">
                    <div class="card m-2 bg-body-tertiary p-2">
                        <h6 class="m-2">Forward from:</h6>
                        <div class="d-flex">
                            <img v-if="message.forward.photo" :src="'data:image/jpeg;base64,' + message.forward.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.forward.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <h6 class="card-title ms-2 mt-3 text-capitalize"> {{ message.forward.name }} </h6>
                        </div>
                    </div>
                </div>

                <!-- sezione testo foto messaggio -->
                <div class="card-body">
                    <img v-if="message.photo" :src="'data:image/jpeg;base64,' + message.photo" class="card-img-top rounded-3" alt="...">
                    <p class="card-text mt-2">{{ message.text }}</p>
                </div>
                <small class="text-end p-2">{{ message.timestamp }}</small>

                <!-- sezione reactions messaggio -->
                <div v-if="message.reactions">
                    <span v-for="item in message.reactions" class="badge text-bg-primary m-1">
                        <span class="text-capitalize">{{ item.user.name }}:</span>  
                        {{ item.reaction }}
                    </span>
                </div>


            </div>
            
            <!-- Se sono gli altri -->
            <div v-else class="card my-4 bg-body-tertiary col-5">

                <!-- sezione risposta messaggio -->
                <div v-if="message.reply != 0">
                    <h6 class="m-2">Reply To:</h6>
                    <div class="card m-2 bg-body-tertiary p-2">
                        <div class="d-flex">
                            <img v-if="message.reply.author.photo" :src="'data:image/jpeg;base64,' + message.reply.author.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.reply.author.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <h6 class="card-title ms-2 mt-3 text-capitalize"> {{ message.reply.author.name }} </h6>
                        </div>
                        <div class="card-body">
                            <img v-if="message.reply.photo" :src="'data:image/jpeg;base64,' + message.reply.photo" class="card-img-top rounded-3" alt="...">
                            <p class="card-text mt-2">{{ message.reply.text }}</p>
                        </div>
                    </div>
                </div>


                <!-- sezione header messaggio -->
                <div class="d-flex">
                    <img v-if="message.author.photo" :src="'data:image/jpeg;base64,' + message.author.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.author.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                    <h5 class="card-title ms-2 mt-3 text-capitalize"> {{ message.author.name }} </h5>
                    <div class="ms-auto m-1 row">
                        <form @submit.prevent="replyMessage" class="col">
                            <input type="hidden" name="message_id" :value="message.id">
                            <button class="btn btn-primary" type="submit">Reply</button>
                        </form>

                        <form @submit.prevent="showForwardMessageHandler" class="col">
                            <input type="hidden" name="message_id" :value="message.id">
                            <button class="btn btn-primary" type="submit">Forward</button>
                        </form>

                        <form @submit.prevent="showMessageReactionPopupHandler" class="col">
                            <input type="hidden" name="message_id" :value="message.id">
                            <button class="btn btn-primary" type="submit">React</button>
                        </form>
                    </div>
                </div>

                <!-- sezione forward messaggio -->
                <div v-if="message.forward != 0">
                    <div class="card m-2 bg-body-tertiary p-2">
                        <h6 class="m-2">Forward from:</h6>
                        <div class="d-flex">
                            <img v-if="message.forward.photo" :src="'data:image/jpeg;base64,' + message.forward.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.forward.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <h6 class="card-title ms-2 mt-3 text-capitalize"> {{ message.forward.name }} </h6>
                        </div>
                    </div>
                </div>

                <!-- sezione testo foto messaggio -->
                <div class="card-body">
                    <img v-if="message.photo" :src="'data:image/jpeg;base64,' + message.photo" class="card-img-top rounded-3" alt="...">
                    <p class="card-text mt-2">{{ message.text }}</p>
                </div>
                <small class="text-end p-2">{{ message.timestamp }}</small>

                <!-- sezione reactions messaggio -->
                <div v-if="message.reactions">
                    <span v-for="item in message.reactions" class="badge text-bg-primary m-1"><span class="text-capitalize">{{ item.user.name }}:</span>  {{ item.reaction }}</span>
                </div>
            
            </div>

        </div>

        <!-- preview risposta -->
        <div v-if="replyMessage_data" class="card my-4 bg-body-tertiary col-12 p-2">
            <h5>Replying to: <span class="text-capitalize">{{ replyMessage_data.author.name }}</span></h5>
            <p>{{ replyMessage_data.text }}</p>
            <button @click="replyMessage_data = null, message_reply = 0" class="btn btn-primary">Cancel</button>
        </div>
        
        <!-- Invio messaggio -->
        <form @submit.prevent="sendMessage" class="input-group mb-3">
            <input v-model="message_input" id="message_input" type="text" class="form-control" placeholder="Type a message" aria-label="Type a message" aria-describedby="message_input">
            <button class="btn btn-primary" type="submit" id="send">Send</button>
        </form>

        <!-- <button  @click="refresh" class="btn btn-primary" type="refresh" id="send">Refresh</button> -->
        <div>
            <div v-for="(item, index) in allConversations" :key="index" class="row card my-4">
            
                <ConversationCard v-if="item.cnv_type == 'group'">

                    <template v-slot:conversationImage>
                        <img v-if="item.photo" :src="'data:image/jpeg;base64,' + item.photo" width="100" height="100" class="rounded-1"  style="object-fit: cover;">
                        <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + item.name" width="100" height="100" class="rounded-1" style="object-fit: cover;">
                    </template>

                    <template v-slot:conversationName>
                        <h5 class="card-title text-capitalize">{{ item.name }}</h5>
                    </template>

                    <template v-slot:conversationMessage>
                        <form @submit.prevent="forwardMessage" class="col">
                            <input type="hidden" name="forward_conversation_id" :value="item.id">
                            <button class="btn btn-primary" type="submit">Forward message</button>
                        </form>
                    </template>

                </ConversationCard>

                <ConversationCard v-if="item.cnv_type == 'chat'">

                    <template v-slot:conversationImage>
                        <ConversationUserPhoto :item="item" :auth_id="auth_id" width="100" height="100"/>
                    </template>

                    <template v-slot:conversationName>
                        <h5 v-if="item.participants[0].id != auth_id" class="card-title text-capitalize">{{ item.participants[0].name }}</h5>
                        <h5 v-if="item. participants[1].id != auth_id" class="card-title text-capitalize">{{ item.participants[1].name }}</h5>
                    </template>

                    <template v-slot:conversationMessage>
                        <form @submit.prevent="forwardMessage" class="col">
                            <input type="hidden" name="forward_conversation_id" :value="item.id">
                            <button class="btn btn-primary" type="submit">Forward message</button>
                        </form>
                    </template>
                    
                </ConversationCard>
            
            </div>
        </div>
            
    </div>
</template>
