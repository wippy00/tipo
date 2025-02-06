<script>
export default {
    props: ['show', 'message'],

    data: function () {
        return {
            reactionsTypes: ['😂', '🗿', '😐', '👍', '❤️', '🔥', '🎉', '😢', '😡'],
        }
    },

    methods: {
        closeModal() {
            this.$emit('close')
        },

        async reactMessage(reaction) {

            const conversation_id = 0
            const auth_id = sessionStorage.getItem('id')
            const message_id = this.message.id

            try {
                let response = await this.$axios.put("/conversations/" + conversation_id + "/messages/" + message_id + "/react", {
                    reaction: reaction
                }, {
                    headers: {
                        authorization: auth_id
                    }
                })

            } catch (e) {
                console.error(e.toString())
            }

        },
        handleReaction(reaction) {
            this.reactMessage(reaction)
            this.closeModal()
        }

    }
};
</script>


<template>
    <div class="modal-backdrop fade show"></div>


    <div class="modal fade show" style="display: block;" aria-modal="true" role="dialog">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">

                    <h5 class="modal-title">React</h5>
                    <button type="button" class="btn-close" @click="closeModal" aria-label="Close"></button>

                </div>

                <div class="modal-body">
                    <div class="card">
                        <div class="d-flex">
                            <img v-if="message.author.photo" :src="'data:image/jpeg;base64,' + message.author.photo" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + message.author.name" width="42" height="42" class="rounded-5 mt-2 ms-2" style="object-fit: cover;">
                            <h5 class="card-title ms-2 mt-3 text-capitalize"> {{ message.author.name }} </h5>
                        </div>
                        <div class="card-body">
                            <p class="card-text">{{ message.text }}</p>

                            <div v-if="message.reactions">
                                <span v-for="item in message.reactions" class="badge text-bg-primary m-1"><span class="text-capitalize">{{ item.user.name }}:</span> {{ item.reaction }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="d-flex mt-3">
                        <button v-for="reaction in reactionsTypes" @click="handleReaction(reaction)" class="btn btn-outline-light me-1">{{ reaction }}</button>

                    </div>

                    <button @click="handleReaction('')" class="btn btn-outline-light mt-2">Delete</button>

                </div>
            </div>
        </div>
    </div>


</template>