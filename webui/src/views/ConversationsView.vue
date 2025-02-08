<script>
import ConversationUserPhoto from '@/components/ConversationUserPhoto.vue'
import ConversationCard from '@/components/ConversationCard.vue'

export default {
    components: {
        ConversationUserPhoto,
        ConversationCard

    },
    data: function () {
        return {
            error: null,
            errormsg: null,
            loading: false,

            conversations: null,

            search: "",
            users: null,
            filteredUsers: null,

            auth_id: null,

            showCreateGroup: false,
            group_name: "",
            group_photo: null,
            checked_users: [],

            refreshInterval: null
        }
    },
    watch: {
        search: function () {
            // console.log(this.search)
            if (this.search.length < 3) {
                return this.filteredUsers = []

            }

            if (this.search === ":all") {
                return this.filteredUsers = this.users
            }

            this.filteredUsers = this.users.filter(user => {
                return user.name.toLowerCase().includes(this.search.toLowerCase())
            })
            // console.log(this.filteredUsers)
        }
    },
    methods: {
        async fetchConversations() {
            
            this.error = null

            this.auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/conversations", {
                    headers: {
                        authorization: this.auth_id
                    }
                })
                this.conversations = response.data

                // console.log(this.conversations)

                // expand the last message author
                for (let i = 0; i < this.conversations.length; i++) {
                    if (this.conversations[i].last_message == 0) {
                        continue
                    }

                    let authorFound = false;
                    for (let j = 0; j < this.conversations[i].participants.length; j++) {
                        if (this.conversations[i].participants[j].id == this.conversations[i].last_message.author) {
                            this.conversations[i].last_message.author = this.conversations[i].participants[j];
                            authorFound = true;
                            break;
                        }
                    }

                    if (!authorFound) {
                        for (let k = 0; k < this.users.length; k++) {
                            if (this.users[k].id == this.conversations[i].last_message.author) {
                                this.conversations[i].last_message.author = this.users[k];
                                break;
                            }
                        }
                    }
                }

                // sort by last message timestamp
                this.conversations.sort((a, b) => {
                    const dateA = a.last_message ? new Date(a.last_message.timestamp) : 0;
                    const dateB = b.last_message ? new Date(b.last_message.timestamp) : 0;
                    return dateB - dateA;
                });

            } catch (e) {
                this.error = e.toString()
            }
            this.loading = false
        },
        async fetchAllUsers() {
            
            this.error = null

            this.auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/users", {
                    headers: {
                        authorization: this.auth_id
                    }
                })
                this.users = response.data

            } catch (e) {
                this.error = e.toString()
            }
            this.loading = false

        },

        async startNewChat(user_id) {
            // console.log(user_id)
            this.error = null

            this.auth_id = sessionStorage.getItem('id')
            try {
                var formData = new FormData();
                // formData.append('name', this.group_name);
                // formData.append('photo', this.group_photo);
                formData.append('cnv_type', "chat");
                formData.append('participants', JSON.stringify([user_id]));

                let response = await this.$axios.post("/conversations", formData,
                    {
                        headers: {
                            'Content-Type': 'multipart/form-data',
                            authorization: this.auth_id
                        }
                    })
                this.$router.push('/conversations/' + response.data.id)
            } catch (e) {
                this.error = e.toString()
            }
        },
        async startNewGroup(event) {
            event.preventDefault()
            this.error = null

            this.auth_id = sessionStorage.getItem('id')

            var formData = new FormData();
            formData.append('name', this.group_name);
            formData.append('photo', this.group_photo);
            formData.append('cnv_type', "group");
            formData.append('participants', JSON.stringify(this.checked_users));

            try {
                let response = await this.$axios.post("/conversations", formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        authorization: this.auth_id
                    }
                })
                this.$router.push('/conversations/' + response.data.id)
            } catch (e) {
                this.error = e.toString()
            }
        },
        file_inputHandler(event) {
            this.group_photo = event.target.files[0];
        },
        showCreateGroupHandler() {
            this.showCreateGroup = !this.showCreateGroup
        },

        async refresh() {
            await this.fetchConversations();
            await this.fetchAllUsers();
        }

    },
    async mounted() {
        if (sessionStorage.getItem('logged_in') !== "true") {
            console.log("Not logged in")
            this.$router.push('/')
        }

        this.refresh();

        this.refreshInterval = setInterval(() => { 
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
        <!-- <ErrorMsg v-if="error" :msg="errormsg"></ErrorMsg> -->

        <h1 v-if="loading">Loading...</h1>
        <div class="row">

            <!-- search friends -->
            <div class="col-10">
                <label for="search" class="form-label fw-bold">Search Friends</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-search" viewBox="0 0 16 16">
                            <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0" />
                        </svg>
                    </span>
                    <input v-model="search" type="text" class="form-control" placeholder="For all users type ':all'" aria-label="Search" aria-describedby="Search-field">
                </div>
            </div>
            <div class="col-2">
                <label for="search" class="form-label fw-bold ">Create Group</label>
                <div>
                    <button @click="showCreateGroupHandler" type="button" class="btn btn-primary">Create</button>
                </div>
            </div>

            <ul class="col-3 list-group mt-3">
                <li v-for="(item, index) in filteredUsers" :key="index" class="list-group-item text-capitalize">
                    {{ item.name }}
                    <button @click="startNewChat(item.id)" class="btn btn-primary float-end">Chat</button>
                </li>
            </ul>
            <!-- create group -->
            <div v-if="showCreateGroup" class="col-12 border p-3 rounded-3 mt-3 mx-auto">
                <h1 class="text-center">Create new group</h1>
                <form @submit.prevent="startNewGroup">

                    <div class="mb-3 row">
                        <div class="col-6">
                            <label for="group_name" class="form-label">Name</label>
                            <input v-model="group_name" type="text" class="form-control" id="group_name" placeholder="Group name" aria-describedby="group_name">
                        </div>

                        <div class="col-6">
                            <label for="group_photo" class="form-label">Photo</label>
                            <input v-on:change="file_inputHandler" type="file" id="group_photo" accept="image/*" class="form-control" aria-describedby="group_photo">
                        </div>
                    </div>

                    <div class="mb-3">
                        <label for="photo" class="form-label">Participants</label>

                        <div v-for="(item, index) in users" :key="index" class="form-check">
                            <input v-model="checked_users" class="form-check-input" type="checkbox" :value="item.id" :id="item.id">
                            <label class="form-check-label text-capitalize" :for="item.id"> {{ item.name }} </label>
                        </div>

                    </div>

                    <button type="submit" class="btn btn-primary">Submit</button>
                </form>
            </div>
        </div>




        <!-- conversation list section -->
        <div v-for="(item, index) in conversations" :key="index" class="row card my-4">

            <ConversationCard v-if="item.cnv_type == 'group'">

                <template v-slot:conversationImage>
                    <img v-if="item.photo" :src="'data:image/jpeg;base64,' + item.photo" width="100" height="100" class="rounded-1" style="object-fit: cover;">
                    <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + item.name" width="100" height="100" class="rounded-1" style="object-fit: cover;">
                </template>

                <template v-slot:conversationName>
                    <RouterLink :to="'/conversations/' + item.id">
                        <h5 class="card-title text-capitalize">{{ item.name }}</h5>
                    </RouterLink>
                </template>

                <template v-if="item.last_message.id != 0" v-slot:conversationMessage>
                    <p v-if="item.last_message.text != null" class="card-text text-capitalize mb-0">
                        {{ item.last_message.author.name + ": " }}
                        <small class="text-body-secondary">{{ item.last_message.text }}</small>
                        <small v-if="item.last_message.photo != null"> 🖼️</small>
                    </p>
                    <p v-if="item.last_message.text == null" class="card-text text-capitalize mb-0">
                        {{ item.last_message.author.name + ": " }}
                        <small v-if="item.last_message.photo != null">🖼️</small>
                    </p>

                    <small class="text-body-secondary">{{ item.last_message.timestamp }}</small>
                </template>
            </ConversationCard>

            <ConversationCard v-if="item.cnv_type == 'chat'">

                <template v-slot:conversationImage>
                    <ConversationUserPhoto :item="item" :auth_id="auth_id" width="100" height="100" />
                </template>

                <template v-slot:conversationName>
                    <RouterLink :to="'/conversations/' + item.id">
                        <h5 v-if="item.participants[0].id != auth_id" class="card-title text-capitalize">{{
                            item.participants[0].name }}</h5>
                        <h5 v-if="item.participants[1].id != auth_id" class="card-title text-capitalize">{{
                            item.participants[1].name }}</h5>
                    </RouterLink>
                </template>

                <template v-if="item.last_message.id != 0" v-slot:conversationMessage>
                    <p v-if="item.last_message.text != null" class="card-text text-capitalize mb-0">
                        {{ item.last_message.author.name + ": " }}
                        <small class="text-body-secondary">{{ item.last_message.text }}</small>
                        <small v-if="item.last_message.photo != null"> 🖼️</small>
                    </p>
                    <p v-if="item.last_message.text == null" class="card-text text-capitalize mb-0">
                        {{ item.last_message.author.name + ": " }}
                        <small v-if="item.last_message.photo != null">🖼️</small>
                    </p>

                    <small class="text-body-secondary">{{ item.last_message.timestamp }}</small>
                </template>
            </ConversationCard>

        </div>
    </div>
</template>

<style>
/* .container {
	height: 100vh;
} */
</style>
