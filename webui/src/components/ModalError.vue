<script>
export default {
    props:['error'],
    
    methods: {
        closeModal() {
            this.$emit('close')
        },
  
  }
};
</script>


<template>
    <div v-if="error" class="modal fade show" tabindex="-1" style="display: block;" aria-modal="true" role="dialog">
        <div class="modal-dialog">
            <div class="modal-content bg-danger-subtle">
                <div class="modal-header">
                    <h5 class="modal-title">Error</h5>
                    <button type="button" class="btn-close" @click="closeModal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <p>{{ error }}</p>
                </div>
            </div>
        </div>
    </div>
</template>