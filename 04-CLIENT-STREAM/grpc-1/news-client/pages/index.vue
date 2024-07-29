<template>
    <div class="flex flex-wrap">
        <div class="w-full md:w-2/12">
            <MenuCategoryComponent 
            @category-clicked="clickedCategory"
            @new-content="generateNewContent"
            />
        </div>
        <div class="w-full md:w-2/12">
            <MenuListComponent 
                :news="news"
                :error="error"
                :isLoading="isLoading"
                :fetchNews="fetchNews"
                @news-clicked="newsSelected"
            />
        </div>
        <div class="w-full md:w-8/12">
            <MenuDetailComponent :selectedNews="selectedNews"/>    
        </div>
    </div>
</template>



<script setup lang="ts">
import { watchEffect } from 'vue';
import type { NewsItem } from '~/types';


const { news, error, isLoading, fetchNews } = useGetNews('news', 'GET', null)
const selectedNews = ref<NewsItem | null>(null);

watchEffect(() => {
    if (news.value) {
        console.log('News updated:', news.value)
        // Update your component state or do something with the news
    }
    if (error.value) {
        console.error('Error occurred:', error.value)
        // Handle the error
    }
})

const clickedCategory = async (key: string) => {
    try {
        console.log(key)
        // const { news: otherNews, error: otherError, isLoading: otherLoading, fetchNews: fetchOtherNews } = useGetNews('other-endpoint', 'POST', { someData: 'value' })
        // const { news: otherNews, error: otherError, isLoading: otherLoading, fetchNews: fetchOtherNews } = useGetNews('news', 'GET', null)
        console.log(key)
        await fetchNews()
        
    } catch (e) {
        console.error('Failed to fetch news:', e)
    }
}

const newsSelected = async (title: string) => {
    console.log("Selected news: ", title);
    // loop over news and get the news title same as this title, get the news content 
    // pass the whole NewsItem as props to <MenuDetailComponent />    
    const newsItem = news.value?.find(item => item.title === title) || null;
    selectedNews.value = newsItem;
    console.log(selectedNews.value);
    
}

const generateNewContent = () => {
    // todo generate new content
    console.log("generate new content");
}

onMounted(() => {
    console.log("mounted to the index page");
    
})

</script>