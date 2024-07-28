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
            />
        </div>
        <div class="w-full md:w-8/12">
            <MenuDetailComponent />    
        </div>
    </div>
</template>



<script setup lang="ts">
import { watchEffect } from 'vue';


const { news, error, isLoading, fetchNews } = useGetNews('news', 'GET', null)

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

const generateNewContent = () => {
    // todo generate new content
    console.log("generate new content");
}


onMounted(() => {
    console.log("mounted to the index page");
    
})

</script>