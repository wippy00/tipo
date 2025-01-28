<script>
import ModalError from '@/components/ModalError.vue';
import ModalSuccess from '@/components/ModalSuccess.vue';

export default {
    components: {
        ModalError,
        ModalSuccess
    },
    data: function () {
        return {
            error: null,
            errormsg: null,
            msg: null,
            loading: false,

            auth_id: null,
            auth_name: null,
            auth_photo: null,

            name_input: "",
            photo_input: null

        }
    },

    methods: {
        async editName(event) {
            event.preventDefault();

            if (this.name_input === "") {
                this.error = "Name cannot be empty.";
                return;
            }

            this.loading = true;
            this.error = null;
            this.errormsg = null;

            try {
                let response = await this.$axios.put("/setUserName", {
                    name: this.name_input
                }, {
                    headers: {
                        authorization: this.auth_id
                    }
                });

                if (response.status === 200) {
                    sessionStorage.setItem('name', this.name_input);
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

        photo_inputHandler(event) {
            this.photo_input = event.target.files[0];
        },

        async editPhoto(event) {
            event.preventDefault();
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

                let response = await this.$axios.put("/setUserPhoto", formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'authorization': this.auth_id
                    }
                });

                if (response.status === 200) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                       // rimuove data:image/jpeg;base64,
                        const base64String = e.target.result.split(',')[1];
                        localStorage.setItem('photo', base64String);
                    }
                    reader.readAsDataURL(this.photo_input);

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

        async reload() {
            this.auth_id = sessionStorage.getItem('id');
            this.auth_name = sessionStorage.getItem('name');
            this.auth_photo = localStorage.getItem('photo');
            // console.log(this.auth_photo)
        }
    },

    mounted() {
        if (sessionStorage.getItem('logged_in') !== "true") {
            console.log("Not logged in")
            this.$router.push('/')
        }
        this.reload();
    },
}
</script>

<template>
    <div class="container">
        
        <ModalError :error="error" @close="error = null" />
        
        <ModalSuccess :msg="msg" @close="msg = null" />

        <div class="d-flex mb-5">
            <img v-if="auth_photo !== 'null'" :src="'data:image/jpeg;base64,' + auth_photo" width="164" height="164" class="rounded-circle" style="object-fit: cover;">
            <img v-else :src="'https://placehold.co/100x100/orange/white?text=' + auth_name" width="164" height="164" class="rounded-circle" style="object-fit: cover;">

            <h1 class="text-capitalize ms-3 mt-5 pt-2"> {{ auth_name }} </h1>
        </div>
        

        <label class="form-label">Edit name:</label>
        <form @submit.prevent="editName" class="input-group mb-3">
            <input v-model="name_input" id="name_input" type="text" class="form-control" placeholder="Name"
                aria-label="Name" aria-describedby="name_input">
            <button class="btn btn-primary" type="submit" id="name_input">Edit</button>
        </form>

        <label class="form-label">Edit photo:</label>
        <form @submit.prevent="editPhoto" class="input-group mb-3">
            <input v-on:change="photo_inputHandler" id="photo_input" type="file"  accept="image/*" class="form-control" placeholder="Type a message"
                aria-label="Type a message" aria-describedby="photo_input">
            <button class="btn btn-primary" type="submit" id="photo_input">Edit</button>
        </form>

    </div>

</template>