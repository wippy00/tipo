<script>
import ModalError from '@/components/ModalError.vue';
import ModalSuccess from '@/components/ModalSuccess.vue';
import ChatHeader from '@/components/ChatHeader.vue';

export default {
    components: {
        ChatHeader,
        ModalError,
        ModalSuccess
    },

    data: function () {
        return {
            error: null,
            errormsg: null,
            loading: false,
            msg: null,

            auth_id: null,
            conversations: [],
            users: [],

            name_input: "",
            photo_input: null,


            search: "",
            filteredUsers: null,
            checked_users: []
        }
    },
    watch: {
        search: function() {
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
        async reload() {
            await this.fetchConversations(this.$route.params.id)
            // Filtra gli utenti che non sono presenti in conversations.participants
            let participantIds = this.conversations.participants.map(participant => participant.id);
            this.users = this.users.filter(user => !participantIds.includes(user.id));
            this.checked_users = []
            this.search = ""
        },
        async fetchConversations(conversations_id) {
            this.loading = true
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
        async fetchAllUsers() {
            this.loading = true
            this.error = null

			this.auth_id = sessionStorage.getItem('id')

            try {
                let response = await this.$axios.get("/users", {
                    headers: {
                        authorization: this.auth_id
                    }
                })
                let allUsers = response.data

                // Filtra gli utenti che non sono presenti in conversations.participants
                let participantIds = this.conversations.participants.map(participant => participant.id);
                this.users = allUsers.filter(user => !participantIds.includes(user.id));

            } catch (e) {
                this.error = e.toString()
            }
            this.loading = false

        },

        async editName(event) {
            event.preventDefault();
            this.loading = true
            this.error = null

            if (this.name_input === "") {
                this.error = "Name cannot be empty.";
                return;
            }

            this.loading = true;
            this.error = null;
            this.errormsg = null;

            try {
                let response = await this.$axios.put("/conversations/"+ this.$route.params.id +"/name", {
                    name: this.name_input
                }, {
                    headers: {
                        authorization: this.auth_id
                    }
                });

                if (response.status === 200) {
                    this.msg = "Name updated successfully.";
                    this.reload();
                } else {
                    let json = await response.json();
                    this.error = response.status;
                    this.errormsg = json.message;
                }
            } catch (error) {
                this.error = error;
            }

            this.loading = false;
        },
        async editPhoto(event) {
            event.preventDefault();
            this.loading = true
            this.error = null
            if (this.photo_input == null) {
                this.error = "Photo cannot be empty.";
                return;
            }

            this.loading = true;
            this.error = null;
            this.errormsg = null;

            try {
                let formData = new FormData();
                formData.append('photo', this.photo_input);

                let response = await this.$axios.put("/conversations/"+ this.$route.params.id +"/photo", formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'authorization': this.auth_id
                    }
                });

                if (response.status === 200) {
                    this.msg = "Photo updated successfully.";
                    this.reload();
                } else {
                    let json = await response.json();
                    this.error = response.status;
                    this.errormsg = json.message;
                }
            } catch (error) {
                this.error = error;
            }

            this.loading = false;
        },
        async addUsers() {
            this.loading = true
            this.error = null

            if (this.checked_users.length === 0) {
                this.error = "Select at least one user.";
                return;
            }

            this.loading = true;
            this.error = null;
            this.errormsg = null;

            for (let i = 0; i < this.checked_users.length; i++) {
                try {
                    let response = await this.$axios.post("/conversations/" + this.$route.params.id +"/add/" + this.checked_users[i], {}, {
                        headers: {
                            authorization: this.auth_id
                        }
                    })
                
                    if (response.status === 200) {
                        this.msg = "Users added successfully.";
                        this.reload();
                    } else {
                        let json = await response.json();
                        this.error = response.status;
                        this.errormsg = json.message;
                    }
                } catch (error) {
                    let json = await response.json();
                    this.error = response.status;
                    this.errormsg = json.message;
                    break
                }
            }

            this.loading = false;
        },
        photo_inputHandler(event) {
            this.photo_input = event.target.files[0];
        },

        
        
    },
    async mounted() {
        if (sessionStorage.getItem('logged_in') !== "true") {
            console.log("Not logged in")
            this.$router.push('/')
        }
        

        this.auth_id = sessionStorage.getItem('id')
        await this.fetchConversations(this.$route.params.id)

        if (this.conversations.cnv_type !== 'group'){
            this.$router.push('/')
        }

        await this.fetchAllUsers()
    }
}
</script>

<template>
    <div class="container">
        <ChatHeader :conversations="conversations" :auth_id="auth_id" />

        <ModalError :error="error" @close="error = null" />
        
        <ModalSuccess :msg="msg" @close="msg = null" />        

        <!-- edit name -->
        <label class="form-label fw-bold mt-5">Edit Name:</label>
        <form @submit.prevent="editName" class="input-group mb-3">
            <input v-model="name_input" id="name_input" type="text" class="form-control" placeholder="Name"
            aria-label="Name" aria-describedby="name_input">
            <button class="btn btn-primary" type="submit" id="name_input">Edit</button>
        </form>
        
        <!-- edit photo -->
        <label class="form-label fw-bold">Edit Photo:</label>
        <form @submit.prevent="editPhoto" class="input-group mb-3">
            <input v-on:change="photo_inputHandler" id="photo_input" type="file"  accept="image/*" class="form-control" placeholder="Type a message"
            aria-label="select profile photo" aria-describedby="photo_input">
            <button class="btn btn-primary" type="submit" id="photo_input">Edit</button>
        </form>

        <!-- add friends -->
        <label for="search" class="form-label fw-bold">Add Friends</label>
        <div class="input-group">
            <span class="input-group-text">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-search" viewBox="0 0 16 16">
                    <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0"/>
                </svg>
            </span>
            <input v-model="search" type="text" class="form-control" placeholder="For all users type ':all'" aria-label="Search" aria-describedby="Search-field">
            <button class="btn btn-primary" @click="addUsers" id="photo_input">Add</button>
        </div>
        <ul class="col-3 list-group mt-3 z-1">
            <li v-for="(item, index) in filteredUsers" :key="index"  class="list-group-item text-capitalize">  
                <input v-model="checked_users" class="form-check-input" type="checkbox" :value="item.id" :id="item.id">
                <label class="form-check-label text-capitalize ms-2" :for="item.id"> {{ item.name }} </label>
            </li>
        </ul>




        <div class="row mt-5">
            <div class="col-md-12">
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">Participants</h5>
                        <ul class="list-group list-group-flush">
                            <li v-for="participant in conversations.participants" class="list-group-item">
                                <img v-if="participant.photo" :src="'data:image/jpeg;base64,' + participant.photo" width="100" height="100" class="rounded-5"  style="object-fit: cover;">
                                <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + participant.name" width="100" height="100" class="rounded-5" style="object-fit: cover;">
                                <span class="fs-4 fw-bold text-capitalize ms-2 ">
                                    {{ participant.name }}
                                </span>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>


    
</template>