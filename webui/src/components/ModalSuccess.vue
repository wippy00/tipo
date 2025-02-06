<script>
export default {
    props: ['msg'],

    methods: {
        closeModal() {
            this.$emit('close')
        },

    }
};
</script>


<template>
    <div v-if="msg" class="modal fade show" tabindex="-1" style="display: block;" aria-modal="true" role="dialog">
        <div class="modal-dialog">
            <div class="modal-content bg-success-subtle">
                <div class="modal-header">
                    <h5 class="modal-title">msg</h5>
                    <button type="button" class="btn-close" @click="closeModal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <p>{{ msg }}</p>
                </div>
            </div>
        </div>
    </div>
</template>