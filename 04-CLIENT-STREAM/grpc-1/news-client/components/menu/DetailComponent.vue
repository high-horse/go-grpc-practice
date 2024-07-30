<template>
    <div class="news-details" v-if="selectedNews" >
        <div v-if="selectedNews.title == 'UNTITLED'">
            <MenuDetailsNewDertailComponent />
        </div>
        <div v-else>
            <MenuDetailsShowDetailComponent  :selectedNews="selectedNews" />    
        </div>   
    </div>
    <div v-else>
        
    </div>

</template>

<script setup lang="ts">
    import type { NewsItem } from '~/types';

    const props = defineProps<{
        selectedNews: NewsItem | null;
    }>()
    // const isDisabled = ref(true);

    // const enableEdit = () => {
    //     // remove disabled from the UTextarea
    //     console.log("edit clicked");
        
    //     isDisabled.value = false;
    // };

    // watch(
    //     () => props.selectedNews,
    //     () => {
    //         isDisabled.value = true;
    //     }
    // );
</script>

<style scoped>
    .news-details {
        margin: 0 10px;
        padding: 10px;
        height: 100vh;
        border: 1px solid wheat;
        border-radius: 10px;
    }
</style>