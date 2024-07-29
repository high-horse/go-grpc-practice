<template>
    <div class="news-details" v-if="selectedNews" >
        <b>
            <h1>{{ selectedNews.title }}</h1>
        </b>
        <p>
            <i>Author: {{ selectedNews.author || "unknown" }},&nbsp;&nbsp;&nbsp;&nbsp;</i>
            <i>source: {{ selectedNews.source.name || "unknown" }}</i>
        </p>
        
        <p>Description:</p>
        <UTextarea 
            ref="textareaRef" 
            autoresize
            :disabled="isDisabled"
            placeholder="News Description..." 
            v-model="selectedNews.description" 
            :rows="20"
            :maxrows="20"
        />
            
        
        <p class="flex justify-end"><i>{{ selectedNews.publishedAt }}</i></p>
        <div>
            <div style="width: 80%;"></div>
            <div class="flex justify-end" style="width: 20%;">
                <UButton
                    icon="i-heroicons-pencil-square"
                    size="sm"
                    color="primary"
                    variant="solid"
                    label="Edit"
                    :trailing="false"
                    class="mr-2"
                />
                <UButton
                    icon="i-heroicons-pencil-square"
                    size="sm"
                    color="primary"
                    variant="solid"
                    label="Save"
                    :trailing="true"
                    @click="enableEdit"
                />
            </div>
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
const isDisabled = ref(true);

const enableEdit = () => {
    // remove disabled from the UTextarea
    console.log("edit clicked");
    
    isDisabled.value = false;
}
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