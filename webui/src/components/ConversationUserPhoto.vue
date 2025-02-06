<script>
export default {
    props: ['item', 'auth_id', 'width', 'height'],
}
</script>

<template>
    <img v-if="item.participants[0].photo && item.participants[0].id != auth_id" :src="'data:image/jpeg;base64,' + item.participants[0].photo" :width="width" :height="height" class="rounded-1" style="object-fit: cover;">
    <img v-if="!item.participants[0].photo && item.participants[0].id != auth_id" :src="'https://placehold.co/100x100/orange/white?text=' + item.participants[0].name" :width="width" :height="height" class="rounded-1" style="object-fit: cover;">
    <img v-if="item.participants[1].photo && item.participants[1].id != auth_id" :src="'data:image/jpeg;base64,' + item.participants[1].photo" :width="width" :height="height" class="rounded-1" style="object-fit: cover;">
    <img v-if="!item.participants[1].photo && item.participants[1].id != auth_id" :src="'https://placehold.co/100x100/orange/white?text=' + item.participants[1].name" :width="width" :height="height" class="rounded-1" style="object-fit: cover;">
</template>