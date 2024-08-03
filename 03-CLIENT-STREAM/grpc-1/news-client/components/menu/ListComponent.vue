<template>
    <div>
        <div v-if="isLoading">Loading ...</div>
        <div v-else-if="error">{{ error?.message }}</div>
        <div class="bordered-div bordered-parent" v-else-if="news" >
            <ul>
                <li 
                    v-for="item in news" 
                    :key="item?.source.id" 
                    :class="{ selected: selectedItem === item }"
                    @click="clicked_news(item)"
                >
                    <div class="bordered-div">
                        <b>{{ item?.title }}</b> <i>{{ item?.author }}</i>
                    </div>
                </li>
            </ul>
        </div>
       
    </div>
</template>

<script setup lang="ts">
/* __placeholder__ */
import type { NewsItem } from '~/types';

const props = defineProps<{
    news: Ref<NewsItem[] | null>;
    error: Ref<Error | null>;
    isLoading: Ref<boolean>;
    fetchNews: () => Promise<void>;
}>();

const emit = defineEmits<{
    (event: 'news-clicked', key: string): void
}>()

const selectedItem = ref<NewsItem | null>(null);

const clicked_news = (item: NewsItem) => {
    selectedItem.value = item;
    emit('news-clicked', item.title);
}
</script>

<style scoped>
.bordered-div{
    border: 1px solid rgb(17, 85, 233);
    border-radius: 10px;
    padding: 10px 5px;
    margin: 10px 0;
}

.bordered-parent {
    height: 100vh;
    overflow-y: scroll;
}

.selected {
    background-color: #2f0cf3;
}
</style>