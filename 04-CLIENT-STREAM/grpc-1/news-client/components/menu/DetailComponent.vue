<template>
    <div v-if="selectedNews">
        <h1>{{ selectedNews.title }}</h1>
        <p>{{ selectedNews.author }}</p>
        
        <p>Description:
            {{ selectedNews.description }}
        </p>
        <p>{{ selectedNews.publishedAt }}</p>
    </div>
    <div v-else>
        
    </div>

</template>

<script setup lang="ts">
import type { NewsItem } from '~/types';

const props = defineProps<{
    selectedNews: NewsItem | null
}>()

</script>