<script>
export default {
    props: ['conversations', 'auth_id'],
    methods: {
        openSettings() {
            this.$router.push('/conversations/'+this.$route.params.id+'/settings')
        },
        async leaveConversation() {
            let auth_id = sessionStorage.getItem('id')

            try {
                await this.$axios.delete("/conversations/"+this.$route.params.id+"/leave", {
                    headers: {
                        authorization: auth_id
                    }
                })

                this.$router.push('/')

            } catch (e) {
                this.error = e.toString()
            }
        },
    }
};
</script>

<template>
    <div class="card p-2 bg-body-tertiary col-12">
        <div class="row">
            <div v-if="conversations.cnv_type == 'group'" class="col-10 d-flex align-items-center">
                <img v-if="conversations.photo" :src="'data:image/jpeg;base64,' + conversations.photo" width="100"
                    height="100" class="rounded-1" style="object-fit: cover;">
                <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + conversations.name" width="100"
                    height="100" class="rounded-1" style="object-fit: cover;">
                <div class="d-flex flex-column ms-2">
                    <h1 class="text-capitalize">{{ conversations.name }}</h1>
                    <ul class="list-inline">
                        <li class="list-inline-item text-capitalize" v-for="user in conversations.participants">{{user.name + " -" }}</li>
                        <li class="list-inline-item">...</li>
                    </ul>
                </div>
            </div>
            <div v-if="conversations.cnv_type == 'chat'" class="col-10 d-flex align-items-center">
                <img v-if="conversations.participants[0].photo && conversations.participants[0].id != auth_id" :src="'data:image/jpeg;base64,' + conversations.participants[0].photo" width="100" height="100" class="rounded-1" style="object-fit: cover;">
                <img v-if="!conversations.participants[0].photo && conversations.participants[0].id != auth_id" :src="'https://placehold.co/100x100/orange/white?text=' + conversations.participants[0].name" width="100" height="100" class="rounded-1" style="object-fit: cover;">
                <img v-if="conversations.participants[1].photo && conversations.participants[1].id != auth_id" :src="'data:image/jpeg;base64,' + conversations.participants[1].photo" width="100" height="100" class="rounded-1" style="object-fit: cover;">
                <img v-if="!conversations.participants[1].photo && conversations.participants[1].id != auth_id" :src="'https://placehold.co/100x100/orange/white?text=' + conversations.participants[1].name" width="100" height="100" class="rounded-1" style="object-fit: cover;">
                <div class="d-flex flex-column ms-2">
                    <h1 v-if="conversations.participants[0].id != auth_id" class="text-capitalize">{{ conversations.participants[0].name }}</h1>
                    <h1 v-if="conversations.participants[1].id != auth_id" class="text-capitalize">{{ conversations.participants[1].name }}</h1>
                </div>
            </div>
            <div v-if="conversations.cnv_type == 'group'" class="col-2 d-flex align-items-center justify-content-center">
                <button @click="openSettings" class="btn btn-primary">Settings</button>
                <button @click="leaveConversation" class="btn btn-primary ms-2">Leave</button>
            </div>
            <div v-if="conversations.cnv_type == 'chat'" class="col-2 d-flex align-items-center justify-content-center">
                <button @click="leaveConversation" class="btn btn-primary ms-2">Delete</button>
            </div>
        </div>
    </div>
</template>
